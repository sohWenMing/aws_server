package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetProjectRoot() (project_root string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		_, err = os.Stat(fmt.Sprintf("%s/go.mod", cwd))
		// if we can find the go.mod file, this is the project root, exit
		if err == nil {
			return cwd, nil
		}
		if cwd == "/" {
			log.Fatal("project root could not be found")
		}
		if errors.Is(err, os.ErrNotExist) {
			cwd = filepath.Dir(cwd)
			continue
		} else {
			log.Fatal(fmt.Sprintf("project root could not be found. expected error: %v", err))
		}
	}
}
