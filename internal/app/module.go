package app

import (
	"encoding/json"
	"log"
	"os/exec"
)

type Module struct {
	Path string
	Dir  string
}

func getModule() Module {
	res, err := exec.Command("go", []string{"list", "-json", "-m"}...).Output()
	if err != nil {
		log.Fatalln(err)
	}

	mod := Module{}

	err = json.Unmarshal(res, &mod)
	if err != nil {
		log.Fatalln(err)
	}

	return mod
}
