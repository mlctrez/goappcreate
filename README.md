# goappcreate

<p align="center">
  <img width="192" height="192" src="https://github.com/mlctrez/goappcreate/blob/master/goapp/web/logo-192.png?raw=true">
</p>

## Purpose

Creates a minimal [go-app](https://go-app.dev/) project structure to bootstrap new projects or experiments.

The following features are included:

* A Makefile supporting all the build steps required to run a [go-app](https://go-app.dev/) project.
* An application version is backed by reading the current tag and commit hash from git.
* A proper split between the wasm and server dependencies to reduce the final wasm size.
* Web resources, including the `app.wasm` are embedded in the server binary during the build.
* The server binary uses [servicego](github.com/mlctrez/servicego) to allow installation as a service on
  supported platforms.
* The [gin web framework](https://github.com/gin-gonic/gin) is included, with brotli compression middleware
  pre-configured.
* It supports the wasm file size header to correctly display progress on the loading screen.
* The `app.Handler` configuration is loaded from a json file in the web directory.
* When running in dev mode, go-app updates are automatic.
* Creates a single folder `goapp` so it can fit well with existing codebases.
* Checks are performed before overwriting any files.
* Less than 300 lines of code - easy to understand and modify.

## Usage

* Create a go project with modules `go mod init <modulePath>`
* Initialize a git repository
* Install using `go install github.com/mlctrez/goappcreate@v0.9.0`
* Execute `goappcreate` in the project folder.
  * There is only one option, `-dark`, to use a modified version of `app.css` which overrides the color scheme on the
    loading screen.
* Install dependencies using `go mod tidy`
* Run `make` with no arguments to start the server in dev mode.

[![Go Report Card](https://goreportcard.com/badge/github.com/mlctrez/goappcreate)](https://goreportcard.com/report/github.com/mlctrez/goappcreate)

created by [tigwen](https://github.com/mlctrez/tigwen)
