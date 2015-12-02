package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var goFilePattern = regexp.MustCompile(`^.*\.go$`)

var context = &CompatContext{Packages: map[string]*Package{}}

func process(path string, f os.FileInfo, err error) error {
	if goFilePattern.Match([]byte(path)) {
		fileSet := token.NewFileSet()
		fileContent, _ := ioutil.ReadFile(path)
		file, _ := parser.ParseFile(fileSet, path, fileContent, parser.ParseComments)

		ProcessFile(fileSet, file, context)
	}

	return nil
}

func main() {
	// Scan project files.
	filepath.Walk(".", process)

	// If .gocompat is present - compare.
	if content, err := ioutil.ReadFile("./.gocompat"); err == nil {
		older := map[string]*Package{}
		d := gob.NewDecoder(bytes.NewReader([]byte(content)))
		err = d.Decode(&older)
		if err != nil {
			fmt.Println("Error when decoding.", err)
		}

		if err := ComparePackages(older, context.Packages); err == nil {
			fmt.Println("OK")
			os.Exit(0)
		} else {
			fmt.Println("Not OK")
			os.Exit(1)
		}
	}

	// Store context objects in .gocompat
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(context.Packages)
	if err != nil {
		fmt.Println("Error when encoding.", err)
	}
	ioutil.WriteFile("./.gocompat", buffer.Bytes(), 0644)
}
