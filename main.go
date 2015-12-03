package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const compatIndexFileName = ".gocompat"

var goFilePattern = regexp.MustCompile(`^.*\.go$`)

var context = &InterfaceContext{Packages: map[string]*Package{}}

// Flags.
var (
	forceStore = flag.Bool("f", false, "Store compatibility index even if the current API is not compatible with the previous version.")
)

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
	exitMessage := ""
	exitCode := 0
	shouldStoreIndex := true

	// Parse command-line options.
	flag.Parse()

	// Scan project files.
	filepath.Walk(".", process)

	// If index is present compare current API to the previous version.
	if content, err := ioutil.ReadFile(compatIndexFileName); err == nil {
		older := map[string]*Package{}
		decoder := gob.NewDecoder(bytes.NewReader([]byte(content)))

		if err = decoder.Decode(&older); err == nil {
			if err := ComparePackages(older, context.Packages); err == nil {
				exitMessage = "OK"
			} else {
				exitMessage = "Not OK"
				exitCode = 1
				shouldStoreIndex = false
			}
		} else {
			fmt.Println("Error when decoding compatibility index.", err)
		}
	}

	// Store context objects in index.
	if shouldStoreIndex || *forceStore {
		buffer := bytes.Buffer{}
		encoder := gob.NewEncoder(&buffer)
		err := encoder.Encode(context.Packages)
		if err != nil {
			fmt.Println("Error when encoding compatibility index.", err)
		}
		ioutil.WriteFile(compatIndexFileName, buffer.Bytes(), 0644)
	}

	fmt.Println(exitMessage)
	os.Exit(exitCode)
}
