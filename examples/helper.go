package examples

import (
	"os"
	"path/filepath"
)

const catalogPath = "./examples/catalog.html"

func getCatalogFile() *os.File {
	path := filepath.Join(getCurrentDir(), catalogPath)
	file, err := os.Open(path)
	raisePanic(err)
	return file
}

func getCurrentDir() string {
	ex, err := os.Executable()
	raisePanic(err)
	return filepath.Dir(ex)
}

func raisePanic(err error) {
	if err != nil {
		panic(err)
	}
}
