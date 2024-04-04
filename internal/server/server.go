package server

import (
	"context"
	"database/sql"
	"wb_level_0/internal/nats"
	"wb_level_0/internal/order"
	"wb_level_0/internal/server/orderservice"
	"wb_level_0/internal/server/repository"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	PgHost = "localhost"
	PgUser = "postgres"
	PgPort = "5432"
	PgPass = "1234"
	PgBase = "wb0"
)

type server struct {
	db 			*sql.DB
	httpServer	*http.Server
	ctx 		*context.Context
	service		*orderservice.OrderService
}

func NewServer() *server {
	
	db, err := NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := orderservice.NewOrderService(
		repository.NewCacheRepository(),
		repository.NewSqlRepository(db),
	)

	if err := service.CacheRecovery(); err != nil {
		log.Fatal(err)
	}
	
	httpServer := &http.Server{
		Addr:           "localhost:8080",
		Handler:        newHttpHandler(service),   //.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return &server{
		db: db,
		httpServer: httpServer,
		service: service,
		ctx: &ctx,
	}
}

func (s *server) Start() {

	dataCh := make(chan []byte)

	go s.httpServer.ListenAndServe()

	go nats.Subscriber(dataCh)

	go func() {
		for {
			ord, err := order.UnmarshalOrder(<-dataCh)
			if err != nil {
				log.Println(err)
			}

			if err := s.service.SaveOrder(&ord); err != nil {
				log.Println(err)
			}
		}
	}()
}

func (s *server) Shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.db.Close()
	s.httpServer.Shutdown(*s.ctx)
}

func NewConnection() (*sql.DB, error) {
	db, err := sql.Open(PgUser, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", PgHost, PgPort, PgUser, PgPass, PgBase))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}