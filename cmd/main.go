package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/Invan2/invan_kanban_service/config"
	"github.com/Invan2/invan_kanban_service/events"
	"github.com/Invan2/invan_kanban_service/pkg/kafka"
	"github.com/Invan2/invan_kanban_service/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.Load()
	log := logger.New(cfg.LogLevel, cfg.ServiceName)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	psqlConn, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Error("poostgres", logger.Error(err))
	}

	kafka, err := kafka.NewKafka(ctx, cfg, log)
	if err != nil {
		log.Error("kafka", logger.Error(err))
		return
	}

	pubsubServer, err := events.New(cfg, log, psqlConn, kafka)
	if err != nil {
		log.Fatal("error creating pubSubServer", logger.Error(err))
		return
	}

	server := grpc.NewServer()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.HttpPort))
	if err != nil {
		log.Error("http", logger.Error(err))
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		server.GracefulStop()
		pubsubServer.Shutdown()
	}()

	go pubsubServer.Run(ctx)

	if err := server.Serve(lis); err != nil {
		log.Fatal("serve", logger.Error(err))
		return
	}

}
