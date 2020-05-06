package product_service

import (
	"github.com/aristat/golang-example-app/app/logger"
	"github.com/nats-io/stan.go"
)

type natsService struct {
	logger logger.Logger
}

func (s *natsService) workerHanlder(m *stan.Msg) {
	s.logger.Info("[NATS] Received a message: %s\n", logger.Args(string(m.Data)))
	m.Ack()
}
