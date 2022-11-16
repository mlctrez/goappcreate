package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/mlctrez/goappcreate/internal"
	"io/fs"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"
)

//go:embed goapp/* Makefile
var Files embed.FS

type Step func() error

var steps = []Step{validateGoMod, createGoAppFiles, createLogos, createMakeFile}

var copyAppCss bool

func main() {

	flag.BoolVar(&copyAppCss, "dark", false, "copies in the dark theme app.css")
	flag.Parse()

	for i, step := range steps {
		fmt.Println(i, "execute", runtime.FuncForPC(reflect.ValueOf(step).Pointer()).Name())
		err := step()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func validateGoMod() (err error) {
	if _, err = os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("missing go.mod - it is needed to determine module path")
	}
	var modulePath string
	if modulePath, err = internal.GetModulePath(); err != nil {
		return
	}
	if modulePath == "" {
		return fmt.Errorf("unable to get module path from go.mod")
	}
	return
}

func createGoAppFiles() (err error) {

	modulePath, _ := internal.GetModulePath()
	oldModule := "github.com/mlctrez/goappcreate"

	err = fs.WalkDir(Files, "goapp", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(path, 0755)
		}
		var fileBytes []byte
		fileBytes, err = Files.ReadFile(path)
		if err != nil {
			return err
		}

		// skip logo copy, they are generated
		if strings.HasPrefix(d.Name(), "logo-") {
			return nil
		}

		if !copyAppCss && d.Name() == "app.css" {
			return nil
		}

		if strings.HasSuffix(path, ".go") {
			replaced := strings.ReplaceAll(string(fileBytes), oldModule, modulePath)
			fileBytes = []byte(replaced)
		}

		return writeFileWithCheck(path, fileBytes)
	})
	return
}

func createMakeFile() (err error) {
	var makefile []byte
	if makefile, err = Files.ReadFile("Makefile"); err != nil {
		return
	}
	return writeFileWithCheck("Makefile", makefile)
}

func createLogos() (err error) {
	var pngBytes []byte
	for _, size := range []int{192, 512} {
		if pngBytes, err = internal.CreatePng(size, size); err != nil {
			return err
		}
		path := fmt.Sprintf("goapp/web/logo-%d.png", size)
		if err = writeFileWithCheck(path, pngBytes); err != nil {
			return
		}
	}
	return
}

func writeFileWithCheck(path string, fileBytes []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.WriteFile(path, fileBytes, 0644)
	} else {
		fmt.Printf("skipping file %s - it already exists\n", path)
		return nil
	}
}
