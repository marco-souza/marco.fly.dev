package main

import (
	"log"
	"os"
	"path"
)

// cli to create packages into project folders
func main() {
	// get first argument
	folder := os.Args[1]
	pkg := os.Args[2]

	// check if folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		log.Fatalln("folder does not exist:", err)
	}

	// create pkg folder
	if err := os.Mkdir(path.Join(folder, pkg), 0755); err != nil {
		log.Fatalln("error creating package folder:", err)
	}

	// create pkg go file
	f, err := os.Create(path.Join(folder, pkg, pkg+".go"))
	defer f.Close()
	if err != nil {
		log.Fatalln("error creating package file:", err)
	}

	// write package template
	_, err = f.WriteString("package " + pkg + "\n\nfunc main() {\n\t// code here\n}\n")
	if err != nil {
		log.Fatalln("error writing package file:", err)
	}

	// create test file
	f, err = os.Create(path.Join(folder, pkg, pkg+"_test.go"))
	defer f.Close()
	if err != nil {
		log.Fatalln("error creating package test file:", err)
	}

	// write test template
	_, err = f.WriteString("package " + pkg + "\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {\n\t// test here\n}\n")
	if err != nil {
		log.Fatalln("error writing package test file:", err)
	}

	// close file
	if err := f.Close(); err != nil {
		log.Fatalln("error closing file:", err)
	}
}
