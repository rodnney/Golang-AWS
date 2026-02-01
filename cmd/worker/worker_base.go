package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

type Worker struct {
	sqs      *aws.SQSClient
	queueURL string
	process  func(ctx context.Context, tx *domain.Transaction) error
	logger   logger.Logger
	wg       *sync.WaitGroup
	quit     chan struct{}
}

func NewWorker(sqs *aws.SQSClient, queueURL string, process func(ctx context.Context, tx *domain.Transaction) error, logger logger.Logger, wg *sync.WaitGroup, quit chan struct{}) *Worker {
	return &Worker{
		sqs:      sqs,
		queueURL: queueURL,
		process:  process,
		logger:   logger,
		wg:       wg,
		quit:     quit,
	}
}

func (w *Worker) Start(concurrency int) {
	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go w.run()
	}
}

func (w *Worker) run() {
	defer w.wg.Done()
	w.logger.Info("Worker started", "queue", w.queueURL)

	for {
		select {
		case <-w.quit:
			w.logger.Info("Worker shutting down", "queue", w.queueURL)
			return
		default:
			msgs, err := w.sqs.ReceiveMessages(context.Background(), w.queueURL, 10)
			if err != nil {
				w.logger.Error("Failed to receive messages", "queue", w.queueURL, "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if len(msgs) == 0 {
				continue
			}

			for _, msg := range msgs {
				var tx domain.Transaction
				if err := json.Unmarshal([]byte(msg.Body), &tx); err != nil {
					w.logger.Error("Failed to unmarshal transaction", "body", msg.Body, "error", err)
					// Delete invalid message to avoid loop
					w.sqs.DeleteMessage(context.Background(), w.queueURL, msg.ReceiptHandle)
					continue
				}

				if err := w.process(context.Background(), &tx); err != nil {
					w.logger.Error("Failed to process transaction", "id", tx.ID, "error", err)
					// Exponential backoff or let it return to queue?
					// SQS visibility timeout will eventually make it available again.
					continue
				}

				if err := w.sqs.DeleteMessage(context.Background(), w.queueURL, msg.ReceiptHandle); err != nil {
					w.logger.Error("Failed to delete message", "id", tx.ID, "error", err)
				}
			}
		}
	}
}
