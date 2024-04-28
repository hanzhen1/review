// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"review-job/internal/biz"
	"review-job/internal/conf"
	"review-job/internal/data"
	"review-job/internal/job"
	"review-job/internal/server"
	"review-job/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, kafka *conf.Kafka, elasticsearch *conf.Elasticsearch, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	reader := job.NewKafkaReader(kafka)
	esClient, err := job.NewESClient(elasticsearch)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	jobWorker := job.NewJobWorker(reader, esClient, logger)
	app := newApp(logger, grpcServer, httpServer, jobWorker)
	return app, func() {
		cleanup()
	}, nil
}