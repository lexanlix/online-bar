package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	bar_api "restapi/internal/adapters/api/bar"
	event_api "restapi/internal/adapters/api/event"
	ingredients_api "restapi/internal/adapters/api/ingredients"
	user_api "restapi/internal/adapters/api/user"
	bar_db "restapi/internal/adapters/db/bar"
	event_db "restapi/internal/adapters/db/event"
	ingredients_db "restapi/internal/adapters/db/ingredients"
	menu_db "restapi/internal/adapters/db/menu"
	session_db "restapi/internal/adapters/db/session"
	user_db "restapi/internal/adapters/db/user"
	"restapi/internal/config"
	"restapi/internal/domain/bar"
	"restapi/internal/domain/event"
	"restapi/internal/domain/ingredients"
	"restapi/internal/domain/menu"
	"restapi/internal/domain/user"
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

	// неправильно это, но пока для теста:
	hub := bar_api.NewHub()

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
	sessionRepository := session_db.NewSessionRepository(postgreSQLClient, logger)

	logger.Info("creating event repository")
	eventRepository := event_db.NewRepository(postgreSQLClient, logger)

	logger.Info("creating ingredients repository")
	ingredientsRepository := ingredients_db.NewRepository(postgreSQLClient, logger)

	logger.Info("creating bar repository")
	barRepository := bar_db.NewRepository(postgreSQLClient, logger)

	logger.Info("creating menu repository")
	menuRepository := menu_db.NewRepository(postgreSQLClient, logger)

	logger.Info("register user service")
	userService := user.NewService(userRepository, sessionRepository, logger, hasher, tokenManager,
		cfg.Tokens.AccessTokenTTL, cfg.Tokens.RefreshTokenTTL)

	logger.Info("register event service")
	eventService := event.NewService(eventRepository, logger)

	logger.Info("register ingredients service")
	ingredientsService := ingredients.NewService(ingredientsRepository, logger)

	logger.Info("register bar service")
	barService := bar.NewService(barRepository, logger)

	logger.Info("register menu service")
	menuService := menu.NewService(menuRepository, logger)

	logger.Info("register user handler")
	userHandler := user_api.NewHandler(logger, userService, eventService, barService, menuService)

	logger.Info("register event handler")
	eventHandler := event_api.NewHandler(logger, eventService, userService, barService)

	logger.Info("register ingredients handler")
	ingredientsHandler := ingredients_api.NewHandler(logger, ingredientsService)

	logger.Info("register bar handler")
	barHandler := bar_api.NewHandler(logger, barService, userService, hub)

	userHandler.Register(router)
	eventHandler.Register(router)
	barHandler.Register(router)
	ingredientsHandler.Register(router)

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
