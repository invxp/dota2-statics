package io

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CurrentExecutablePath() (string, string) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Panic(err)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		log.Panic(err)
	}

	return path[0:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator), path[strings.LastIndex(path, string(os.PathSeparator))+len(string(os.PathSeparator)):]
}
