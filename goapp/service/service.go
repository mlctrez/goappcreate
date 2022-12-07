//go:build !wasm

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	brotli "github.com/anargu/gin-brotli"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goappcreate/goapp"
	"github.com/mlctrez/goappcreate/goapp/compo"
	"github.com/mlctrez/servicego"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func Entry() {
	compo.Routes()
	servicego.Run(&Service{})
}

var _ servicego.Service = (*Service)(nil)

var DevEnv = os.Getenv("DEV")
var IsDev = DevEnv != ""

type Service struct {
	servicego.Defaults
	serverShutdown func(ctx context.Context) error
}

func (s *Service) Start(_ service.Service) (err error) {

	var listener net.Listener
	address := listenAddress()
	if listener, err = net.Listen("tcp4", address); err != nil {
		return
	}
	fmt.Printf("listening on http://%s\n", address)

	engine := buildGinEngine()
	if err = setupStaticHandlers(engine); err != nil {
		return
	}
	if err = setupApiEndpoints(engine); err != nil {
		return
	}
	if err = setupGoAppHandler(engine); err != nil {
		return
	}

	server := &http.Server{Handler: engine}
	s.serverShutdown = server.Shutdown

	go func() {
		var serveErr error
		if strings.HasSuffix(listener.Addr().String(), ":443") {
			serveErr = server.ServeTLS(listener, "cert.pem", "cert.key")
		} else {
			serveErr = server.Serve(listener)
		}
		if serveErr != nil && serveErr != http.ErrServerClosed {
			_ = s.Log().Error(err)
		}
	}()

	return nil
}

func (s *Service) Stop(_ service.Service) (err error) {
	if s.serverShutdown != nil {

		stopContext, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		err = s.serverShutdown(stopContext)
		if errors.Is(err, context.Canceled) {
			os.Exit(-1)
		}
	}
	_ = s.Log().Info("http.Server.Shutdown success")
	return
}

func listenAddress() string {
	if address := os.Getenv("ADDRESS"); address != "" {
		return address
	}
	if port := os.Getenv("PORT"); port == "" {
		return "localhost:8080"
	} else {
		return "localhost:" + port
	}

}

func buildGinEngine() (engine *gin.Engine) {

	if !IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	engine = gin.New()

	// required for go-app to work correctly
	engine.RedirectTrailingSlash = false

	if IsDev {
		// omit some common paths to reduce startup logging noise
		skipPaths := []string{
			"/app.css", "/app.js", "/app-worker.js", "/manifest.webmanifest", "/wasm_exec.js",
			"/web/logo-192.png", "/web/logo-512.png", "/web/app.wasm"}
		engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPaths}))
	}
	engine.Use(gin.Recovery(), brotli.Brotli(brotli.DefaultCompression))

	return
}

func setupStaticHandlers(engine *gin.Engine) (err error) {

	var wasmFile fs.File
	if wasmFile, err = goapp.WebFs.Open("web/app.wasm"); err != nil {
		return
	}
	defer func() { _ = wasmFile.Close() }()

	var stat fs.FileInfo
	if stat, err = wasmFile.Stat(); err != nil {
		return
	}
	wasmSize := stat.Size()

	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Wasm-Content-Length", fmt.Sprintf("%d", wasmSize))
		c.Next()
	})

	staticHandler := http.FileServer(http.FS(goapp.WebFs))
	engine.GET("/web/:path", gin.WrapH(staticHandler))

	if _, err = fs.Stat(goapp.WebFs, "web/app.css"); err == nil {
		//  use provided web/app.css instead of app.css provided by go-app
		engine.GET("/app.css", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/web/app.css")
		})
	} else {
		err = nil
	}

	return
}

func setupApiEndpoints(engine *gin.Engine) error {
	// setup other api endpoints here
	return nil
}

func setupGoAppHandler(engine *gin.Engine) (err error) {

	var handler *app.Handler

	// if dynamic customization of other app.Handler fields is required,
	// just build programmatically and skip the goAppHandlerFromJson() call
	if handler, err = goAppHandlerFromJson(); err != nil {
		return
	}

	handler.WasmContentLengthHeader = "Wasm-Content-Length"
	handler.Env["DEV"] = os.Getenv("DEV")

	if IsDev {
		handler.AutoUpdateInterval = time.Second * 3
		handler.Version = ""
	} else {
		handler.AutoUpdateInterval = time.Hour
		handler.Version = fmt.Sprintf("%s@%s", goapp.Version, goapp.Commit)
	}

	engine.NoRoute(gin.WrapH(handler))
	return nil
}

func goAppHandlerFromJson() (handler *app.Handler, err error) {

	var file fs.File
	if file, err = goapp.WebFs.Open("web/handler.json"); err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	handler = &app.Handler{}
	if err = json.NewDecoder(file).Decode(handler); err != nil {
		return
	}

	return
}
