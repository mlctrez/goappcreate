package goapp

import (
	"crypto/sha1"
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"time"
)

var (
	Version        string
	Commit         string
	runtimeVersion string
)

func RuntimeVersion() string {
	if runtimeVersion == "" {
		if IsDevelopment() {
			t := time.Now().UTC().String()
			runtimeVersion = fmt.Sprintf(`%x`, sha1.Sum([]byte(t)))
		} else {
			runtimeVersion = fmt.Sprintf("%s@%s", Version, Commit)
		}
	}
	return runtimeVersion
}

func IsDevelopment() bool {
	return app.Getenv("DEV") != ""
}

func UpdateInterval() time.Duration {
	if IsDevelopment() {
		return 3 * time.Second
	}
	return 24 * time.Hour
}
