# goappcreate

[![Go Report Card](https://goreportcard.com/badge/github.com/mlctrez/goappcreate)](https://goreportcard.com/report/github.com/mlctrez/goappcreate)

## Purpose

Creates a minimal [go-app](https://go-app.dev/) project structure with the following features:

* A Makefile supporting all the build steps required to run a go-app.
* Web resources, including the app.wasm are embedded in the server binary.
* The server binary uses [servicego](github.com/mlctrez/servicego) for installation of the binary as a service
  on certain platforms.
* The gin router with brotli compression is included to facilitate building web apis. 
* Supports the wasm file size header to correctly display progress on the loading screen.
* The app.Handler configuration is loaded from a json file in the web directory.
* When running in dev mode, go-app updates are automatic.
* Creates a single folder `goapp` so it can fit well with existing codebases.
* Checks are performed before before overwriting any files.






created by [tigwen](https://github.com/mlctrez/tigwen)
