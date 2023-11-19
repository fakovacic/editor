package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	ioType "github.com/fakovacic/editor/internal/app/editor/io"
	ioFile "github.com/fakovacic/editor/internal/app/editor/io/file"
	ioHttp "github.com/fakovacic/editor/internal/app/editor/io/http"
	ioMiddleware "github.com/fakovacic/editor/internal/app/editor/io/middleware"
	editorMiddleware "github.com/fakovacic/editor/internal/app/editor/middleware"
	"github.com/fakovacic/editor/internal/app/hub"
	"github.com/fakovacic/editor/internal/app/hub/colors"
	hubMiddleware "github.com/fakovacic/editor/internal/app/hub/middleware"
	versioningType "github.com/fakovacic/editor/internal/app/versioning"
	versioningFile "github.com/fakovacic/editor/internal/app/versioning/file"
	versioningHTTP "github.com/fakovacic/editor/internal/app/versioning/http"
	versioningMiddleware "github.com/fakovacic/editor/internal/app/versioning/middleware"
	"github.com/fakovacic/editor/internal/app/web"
	"github.com/fakovacic/editor/internal/app/web/handler"
	handlerMiddleware "github.com/fakovacic/editor/internal/app/web/handler/middleware"
	webMiddleware "github.com/fakovacic/editor/internal/app/web/middleware"
	"github.com/fakovacic/editor/internal/app/write/validator"
	writeValidatorMiddleware "github.com/fakovacic/editor/internal/app/write/validator/middleware"
	"github.com/fakovacic/editor/internal/health"
	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	jsoniter "github.com/json-iterator/go"
)

const errorChan int = 10

//go:embed templates/*
var content embed.FS

//go:embed static/*
var static embed.FS

func main() {
	ctx := context.Background()

	// io
	var editorIO editor.IO

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		log.Fatal(ctx, "FILE_PATH environment variable not set")
	}

	var ioFileType ioType.Type

	fileIO := os.Getenv("FILE_IO")
	if fileIO == "" {
		log.Fatal(ctx, "FILE_IO environment variable not set")
	}

	err := ioFileType.Parse(fileIO)
	if err != nil {
		log.Fatal(ctx, "FILE_IO environment variable not valid")
	}

	switch ioFileType {
	case ioType.File:
		editorIO = ioFile.New(filePath)
		editorIO = ioMiddleware.NewLogMiddleware(editorIO, ioType.File)
	case ioType.HTTP:
		editorIO = ioHttp.New(filePath, http.DefaultClient)
		editorIO = ioMiddleware.NewLogMiddleware(editorIO, ioType.HTTP)
	default:
		log.Fatal(ctx, "FILE_IO environment variable not valid")
	}

	// versioning
	versioningIO := os.Getenv("VERSIONS_IO")
	if versioningIO != "" {
		versionPath := os.Getenv("VERSIONS_PATH")
		if versionPath == "" {
			log.Fatal(ctx, "VERSIONS_PATH environment variable not set")
		}

		var versioning app.Versioning

		var versioningFileType versioningType.Type

		err = versioningFileType.Parse(versioningIO)
		if err != nil {
			log.Fatal(ctx, "VERSIONS_IO environment variable not valid")
		}

		switch versioningFileType {
		case versioningType.File:
			versioning = versioningFile.New(versionPath, time.Now)
			versioning = versioningMiddleware.NewLogMiddleware(versioning, versioningType.File)
		case versioningType.HTTP:
			versioning = versioningHTTP.New(versionPath, http.DefaultClient, time.Now)
			versioning = versioningMiddleware.NewLogMiddleware(versioning, versioningType.HTTP)
		default:
			log.Fatal(ctx, "VERSIONS_IO environment variable not set")
		}

		editorIO = ioMiddleware.NewVersioningMiddleware(editorIO, versioning)
	}

	// ttl
	var connTTL *time.Duration

	connTTLStr := os.Getenv("CONN_TTL")
	if connTTLStr != "" {
		ttl, err := time.ParseDuration(connTTLStr)
		if err != nil {
			log.Fatal(ctx, "CONN_TTL environment variable not valid")
		}

		connTTL = &ttl
	}

	// editor
	editor := editor.New(editorIO)
	editor = editorMiddleware.NewLogMiddleware(editor)

	// write validator
	writeValidator := validator.New()
	writeValidator = writeValidatorMiddleware.NewLogMiddleware(writeValidator)

	// hub
	colors := colors.New()
	hb := hub.New(colors)
	hb = hubMiddleware.NewLogMiddleware(hb)

	// web service
	service := web.New(editor, hb, writeValidator, connTTL)
	service = webMiddleware.NewLogMiddleware(service)

	// handler
	h := handler.New(service)

	engine := html.NewFileSystem(http.FS(content), ".html")
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	app := fiber.New(fiber.Config{
		Views:       engine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(handlerMiddleware.Logger())
	app.Use(handlerMiddleware.ReqID())

	app.Get("/login", h.LoginForm())
	app.Post("/login", h.Login())
	app.Get("/", h.Index())

	// app.Static("/", "./static")

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(static),
	}))

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", h.WS())

	var (
		httpAddr   = "0.0.0.0:8080"
		healthAddr = "0.0.0.0:8081"
	)

	// health router
	healthServer := health.StartServer()

	errChan := make(chan error, errorChan)

	go func() {
		log.Info(ctx, fmt.Sprintf("Health service listening on %s", healthAddr))
		errChan <- healthServer.Listen(healthAddr)
	}()
	go func() {
		log.Info(ctx, fmt.Sprintf("HTTP service listening on %s", httpAddr))
		errChan <- app.Listen(httpAddr)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case e := <-errChan:
			if e != nil {
				log.Fatal(ctx, e.Error())
			}

			return
		case s := <-signalChan:
			log.Info(ctx, "Captured %v. Exiting...", s)
			health.SetHealthStatus(http.StatusServiceUnavailable)

			err := app.Shutdown()
			if err != nil {
				log.Fatal(ctx, err.Error())
			}

			return
		}
	}
}
