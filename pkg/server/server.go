package server

import (
	"context"
	"net"
	"net/http"
	"nprn/internal/config"
	"nprn/pkg/logging"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(handler http.Handler, logger *logging.Logger, cfg *config.Config) error {

	var listener net.Listener
	var listenerErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("creating socket...")
		listener, listenerErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)

	} else {
		logger.Info("listen tcp")
		listener, listenerErr = net.Listen("tcp", cfg.Listen.BindIP+":"+cfg.Listen.Port)
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenerErr != nil {
		logger.Fatal(listenerErr)
	}

	s.httpServer = &http.Server{
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.httpServer.Serve(listener)
}

func (s *Server) Close(ctx context.Context) error {
	ctx.Done()

	return s.httpServer.Close()
}
