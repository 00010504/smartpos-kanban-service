package events

import (
	"context"

	"github.com/Invan2/invan_kanban_service/config"
	"github.com/Invan2/invan_kanban_service/events/kanban_service/category"
	"github.com/Invan2/invan_kanban_service/pkg/kafka"
	"github.com/Invan2/invan_kanban_service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type PubsubServer struct {
	cfg   config.Config
	log   logger.Logger
	db    *sqlx.DB
	kafka kafka.KafkaI
}

func New(cfg config.Config, log logger.Logger, db *sqlx.DB, kafka kafka.KafkaI) (*PubsubServer, error) {

	// Product topics)
	kafka.AddPublisher("v2.catalog_service.product.created")

	return &PubsubServer{
		cfg:   cfg,
		log:   log,
		db:    db,
		kafka: kafka,
	}, nil
}

func (s *PubsubServer) Run(ctx context.Context) {
	categoryService := category.New(s.cfg, s.log, s.db, s.kafka)
	categoryService.RegisterConsumers()

	s.kafka.RunConsumers()
}

func (s *PubsubServer) Shutdown() error {
	return s.kafka.Shutdown()
}
