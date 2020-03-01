package main

import (
	"encoding/json"
	"fmt"
	"github.com/PaluMacil/decisive-oak/parse"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	files, err := filepath.Glob("*.data.txt")
	if err != nil {
		fmt.Printf("finding data files: %v", err)
		os.Exit(1)
	}
	for _, filename := range files {
		sample, err := parse.FromFile(filename)
		if err != nil {
			fmt.Printf("parsing %s: %v", filename, err)
			os.Exit(1)
		}
		jsonData, err := json.MarshalIndent(sample, "", "  ")
		if err != nil {
			fmt.Printf("marshalling sample to JSON: %v", err)
			os.Exit(1)
		}
		outFilename := strings.TrimSuffix(filename, path.Ext(filename)) + ".json"
		err = ioutil.WriteFile(outFilename, jsonData, 0644)
		if err != nil {
			fmt.Printf("writing output filename: %v", err)
			os.Exit(1)
		}
	}
}


