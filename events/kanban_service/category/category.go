package category

import (
	"github.com/Invan2/invan_kanban_service/config"
	"github.com/Invan2/invan_kanban_service/pkg/kafka"
	"github.com/Invan2/invan_kanban_service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type CategoryService struct {
	cfg   config.Config
	log   logger.Logger
	kafka kafka.KafkaI
}

func New(cfg config.Config, log logger.Logger, db *sqlx.DB, kafka kafka.KafkaI) *CategoryService {
	return &CategoryService{
		cfg:   cfg,
		log:   log,
		kafka: kafka,
	}
}

func (c *CategoryService) RegisterConsumers() {
	catalogRoute := "v1.catalog_service.category"

	c.kafka.AddConsumer(
		catalogRoute+".create", // topic
		c.Create,               // handlerFunction
	)
}
