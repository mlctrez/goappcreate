package goapp

import (
	"embed"
)

// embed web and all subdirectories

//go:embed web/*
var WebFs embed.FS
