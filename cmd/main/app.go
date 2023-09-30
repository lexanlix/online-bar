package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	user_api "restapi/internal/adapters/api/user"
	user_db "restapi/internal/adapters/db/user"
	"restapi/internal/config"
	"restapi/internal/domain/user"
	"restapi/internal/session"
	"restapi/pkg/auth"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/hash"
	"restapi/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("create router")

	router := httprouter.New()

	cfg := config.GetConfig()

	hasher := hash.NewSHA1Hasher("passHasherSalt123")

	tokenManager, err := auth.NewManager(cfg.Tokens.SigningKey)
	if err != nil {
		logger.Fatalf("failed to create token manager: %v", err)
	}

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	logger.Info("creating user repository")
	userRepository := user_db.NewRepository(postgreSQLClient, logger, hasher)

	logger.Info("creating session repository")
	sessionRepository := session.NewSessionRepository(postgreSQLClient, logger)

	logger.Info("register user service")

	userService := user.NewService(userRepository, sessionRepository, logger, hasher, tokenManager)

	logger.Info("register user handler")
	userHandler := user_api.NewHandler(logger, userService)
	userHandler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		// Нужно определить, где находится бинарник запуска
		// path/to/binary
		// Dir() -- /path/to
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)

		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		panic(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
