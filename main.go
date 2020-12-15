package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func create(path string) {
	_, err := os.Stat(path)

	// TODO(#1): is this error handling enough?
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func currentdir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	return path
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// // TODO(#2): Extract and configure structure.yaml to hold defaults for how each individual project will be set up
// func initStructure() {

// }

func copy(templatePath string, toPath string) {
	srcFile, err := os.Open(templatePath)
	check(err)
	defer srcFile.Close()

	destFile, err := os.Create(toPath) // creates if file doesn't exist
	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	check(err)

	err = destFile.Sync()
	check(err)
}

func getTemplate() string {
	templates := os.Getenv("PROJECT_TEMPLATES")

	if templates == " " {
		fmt.Printf("Enviromental variable not properly set")
		os.Exit(1)
	}

	var template string
	switch os.Args[2] {
	case "c":
		template = "clang.c"
	case "go":
		template = "golang.go"
	case "py":
		template = "empty.py"
	}

	return fmt.Sprintf("%v/%v", templates, template)
}

// TODO(#3): add chdir functionality, with a flag to toggle
func chdir() {

}

func mkproj() (string, error) {
	path := filepath.Join(currentdir(), os.Args[1])

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		return path, nil
	}

	return " ", fmt.Errorf("%v - Path already exists", path)
}

func getProjName() string {
	mainFile := "main"
	if len(os.Args) == 4 {
		mainFile = os.Args[3]
	}

	return mainFile
}

func getProject() string {
	project, err := mkproj()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return filepath.Join(project, fmt.Sprintf("%v.%v", getProjName(), os.Args[2]))
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Invalid number of arguments")
		os.Exit(1)
	}

	copy(getTemplate(), getProject())
}
