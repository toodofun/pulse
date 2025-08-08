// Copyright 2025 The Toodofun Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/toodofun/pulse/internal/config"
	"github.com/toodofun/pulse/internal/controller"
	"github.com/toodofun/pulse/internal/infra"
	"github.com/toodofun/pulse/internal/service"
)

type Server struct {
	ctx    context.Context
	server *http.Server
}

//go:embed static
var fs embed.FS

func New(cfg *config.Config) (*Server, error) {
	logrus.Infof("running with config: %+v", config.Current())
	ctx := context.WithValue(
		context.Background(),
		config.ContextKeyInstanceId,
		fmt.Sprintf("main-%s", uuid.NewString()[:5]),
	)

	ctx = context.WithValue(ctx, config.ContextKeyConfig, cfg)

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), Static("/", NewStaticFileSystem(fs, "static")))
	engine.MaxMultipartMemory = 8 << 20
	engine.NoRoute(func(ctx *gin.Context) {
		controller.Reply(ctx, controller.CodeNotFound, nil)
	})

	api := engine.Group(cfg.Server.Prefix)
	api.Use(AuthCheck())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	db, err := infra.NewDatabase(cfg.Database)
	if err != nil {
		return nil, err
	}

	monitorService := service.NewMonitorService()
	userService := service.NewUserService()
	services := []Service{
		monitorService,
		userService,
	}

	for _, svc := range services {
		if err := svc.Initialize(db); err != nil {
			return nil, fmt.Errorf("failed to initialize service: %w", err)
		}
	}

	controllers := []Controller{
		controller.NewMonitorController(monitorService),
		controller.NewUserController(userService),
	}

	for _, ctrl := range controllers {
		ctrl.RegisterRoute(api)
	}

	return &Server{
		ctx:    ctx,
		server: server,
	}, nil
}

func (s *Server) Run() error {
	errCh := make(chan error, 1)

	go func() {
		logrus.Infof("Server listening on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case ch := <-sig:
		logrus.Infof("signal received: %s, shutting down server...", ch.String())
		return s.Shutdown()
	case err := <-errCh:
		return err
	}
}

func (s *Server) Shutdown() error {
	if s.server == nil {
		return fmt.Errorf("server is not initialized")
	}

	cfg := s.ctx.Value(config.ContextKeyConfig).(*config.Config)
	ctx, cancel := context.WithTimeout(s.ctx, cfg.Server.GracePeriod)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
