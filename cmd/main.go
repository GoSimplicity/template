package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoSimplicity/template/internal/pkg/di"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	if err := di.InitViper(); err != nil {
		panic(err)
	}

	app := di.InitializeApp()
	logger := zap.L()
	defer func() {
		_ = logger.Sync()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	for _, consumer := range app.Consumer {
		if err := consumer.Start(ctx); err != nil {
			logger.Error("启动消费者失败", zap.Error(err))
			errCh <- err
			break
		}
	}

	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.Server,
	}

	go func() {
		logger.Info("HTTP server 启动", zap.Int("port", port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case errCh <- err:
			default:
			}
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			logger.Error("服务异常退出", zap.Error(err))
		}
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server 关闭失败", zap.Error(err))
	}

	logger.Info("服务已退出")
}
