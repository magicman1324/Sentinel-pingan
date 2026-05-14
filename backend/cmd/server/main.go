package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pingan/monitor-backend/internal/api/rest"
	"github.com/pingan/monitor-backend/internal/repository"
	"github.com/pingan/monitor-backend/internal/service"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	dsn := getEnv("TIDB_DSN", "root:@tcp(127.0.0.1:4000)/monitor?parseTime=true")
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("tidb: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	rdb := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})
	defer rdb.Close()

	ruleRepo := repository.NewRuleRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	svc := service.NewService(ruleRepo, alertRepo, rdb)

	// Preload rules cache
	if err := svc.SyncEnabledRulesToRedis(context.Background()); err != nil {
		log.Printf("redis preload: %v", err)
	}

	// REST
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	rest.RegisterRoutes(router, svc, db.DB, rdb)

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// gRPC
	grpcSrv := grpc.NewServer()
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() { log.Printf("rest: %v", httpSrv.ListenAndServe()) }()
	go func() { log.Printf("grpc: %v", grpcSrv.Serve(lis)) }()

	log.Println("[backend] :8080 (REST) :9090 (gRPC)")

	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpSrv.Shutdown(shutdownCtx)
	grpcSrv.GracefulStop()
	log.Println("[backend] stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
