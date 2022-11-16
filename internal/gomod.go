package internal

import (
	"encoding/json"
	"os/exec"
)

type GoModEditJson struct {
	Module struct {
		Path string `json:"Path"`
	} `json:"Module"`
}

func GetModulePath() (string, error) {
	gme := &GoModEditJson{}

	output, err := exec.Command("go", "mod", "edit", "-json").CombinedOutput()
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(output, gme)
	if err != nil {
		return "", err
	}
	return gme.Module.Path, nil

}
