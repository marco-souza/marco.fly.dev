package main

import (
	"log/slog"
	"os"
	"path"
)

var logger = slog.With("service", "cli")

// cli to create packages into project folders
func main() {
	// get first argument
	folder := os.Args[1]
	pkg := os.Args[2]

	// check if folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		logger.Error("folder does not exist", "err", err)
	}

	// create pkg folder
	if err := os.Mkdir(path.Join(folder, pkg), 0755); err != nil {
		logger.Error("error creating package folder", "err", err)
	}

	// create pkg go file
	f, err := os.Create(path.Join(folder, pkg, pkg+".go"))
	defer f.Close()
	if err != nil {
		logger.Error("error creating package file", "err", err)
	}

	// write package template
	_, err = f.WriteString("package " + pkg + "\n\nfunc main() {\n\t// code here\n}\n")
	if err != nil {
		logger.Error("error writing package file", "err", err)
	}

	// create test file
	f, err = os.Create(path.Join(folder, pkg, pkg+"_test.go"))
	defer f.Close()
	if err != nil {
		logger.Error("error creating package test file", "err", err)
	}

	// write test template
	_, err = f.WriteString("package " + pkg + "\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {\n\t// test here\n}\n")
	if err != nil {
		logger.Error("error writing package test file", "err", err)
	}

	// close file
	if err := f.Close(); err != nil {
		logger.Error("error closing file", "err", err)
	}
}
