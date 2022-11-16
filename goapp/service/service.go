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
	dev := goapp.IsDevelopment()

	if dev {
		fmt.Printf("running on http://%s\n", address)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	// required for go-app to work correctly
	engine.RedirectTrailingSlash = false

	if dev {
		engine.Use(gin.Logger(), gin.Recovery(), brotli.Brotli(brotli.DefaultCompression))
	} else {
		engine.Use(gin.Recovery(), brotli.Brotli(brotli.DefaultCompression))
	}

	if err = setupGinStaticHandlers(engine); err != nil {
		return
	}

	// other api endpoints can go here

	var handler *app.Handler
	if handler, err = BuildHandler(); err != nil {
		return
	}
	engine.NoRoute(gin.WrapH(handler))

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

func setupGinStaticHandlers(engine *gin.Engine) (err error) {

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

func BuildHandler() (handler *app.Handler, err error) {

	var file fs.File
	if file, err = goapp.WebFs.Open("web/handler.json"); err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	handler = &app.Handler{}
	if err = json.NewDecoder(file).Decode(handler); err != nil {
		return
	}
	handler.Version = goapp.RuntimeVersion()
	handler.AutoUpdateInterval = goapp.UpdateInterval()
	if goapp.IsDevelopment() {
		handler.Env["DEV"] = "1"
	}
	handler.WasmContentLengthHeader = "Wasm-Content-Length"

	return
}
