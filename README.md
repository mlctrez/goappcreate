# goappcreate

<p align="center">
  <img width="192" height="192" src="https://github.com/mlctrez/goappcreate/blob/master/goapp/web/logo-192.png?raw=true">
</p>

## Purpose

Creates a minimal [go-app](https://go-app.dev/) project structure to bootstrap new projects or experiments.

The following features are included:

* A Makefile supporting all the build steps required to run a [go-app](https://go-app.dev/) project.
* An application version backed by the current tag and commit hash from git.
* A proper split between the wasm and server dependencies to reduce the final wasm size.
* Web resources, including the `app.wasm` are embedded in the server binary during the build.
* The server binary uses [servicego](github.com/mlctrez/servicego) to allow installation as a service on
  supported platforms. 
* The [gin web framework](https://github.com/gin-gonic/gin) is included and brotli compression middleware
  is pre-configured.
* It supports the wasm file size header to correctly display progress on the loading screen when using compression.
* The `app.Handler` configuration is loaded from a json file in the web directory.
* When running in dev mode, [app updates](https://go-app.dev/lifecycle#listen-for-app-updates) are automatic.
* Uses just a single folder `goapp` so it can fit into an existing codebase.
* Checks are performed to prevent overwriting any existing files.
* Less than 300 lines of code - easy to understand and modify.

## Usage

* Create a go project with modules - `go mod init <modulePath>`
* Initialize a git repository to pul the tag and hash information from.
  * This is not required for development, but is highly recommended for other environments. 
* Install using `go install github.com/mlctrez/goappcreate@v1.2.1`
* Execute `goappcreate` in the project folder.
    * There is only one option, `-dark`, to substitute in a dark theme for the loading screen.
* Install dependencies using `go mod tidy`
* Run `make` with no arguments to start the server in dev mode.

[![Go Report Card](https://goreportcard.com/badge/github.com/mlctrez/goappcreate)](https://goreportcard.com/report/github.com/mlctrez/goappcreate)

created by [tigwen](https://github.com/mlctrez/tigwen)
