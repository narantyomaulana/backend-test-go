package services

import (
	"log"

	"e-wallet-api/pkg/rabbitmq"
)

type QueueService struct {
	rabbitMQ      *rabbitmq.RabbitMQ
	walletService *WalletService
}

func NewQueueService(rabbitMQ *rabbitmq.RabbitMQ, walletService *WalletService) *QueueService {
	return &QueueService{
		rabbitMQ:      rabbitMQ,
		walletService: walletService,
	}
}

func (s *QueueService) StartTransferWorker() error {
	return s.rabbitMQ.ConsumeMessages("transfer_queue", s.handleTransferMessage)
}

func (s *QueueService) handleTransferMessage(message rabbitmq.TransferMessage) {
	log.Printf("Processing transfer: %+v", message)
	
	if err := s.walletService.ProcessTransfer(message); err != nil {
		log.Printf("Failed to process transfer %s: %v", message.TransferID, err)
		// In production, you might want to implement retry logic or dead letter queue
	}
}
