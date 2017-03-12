package main

import (
	"fmt"
	"log"
	"strings"
	"path/filepath"
	"io/ioutil"
	"os"
)

var (
	// Configuration from CLI args
	DirToParse string
	OutFileName string
	GeneratorFileName string
)

func PathExists(path string) bool {
	if _, e := os.Stat(path); e != nil {
		// either non existent or smth wrong with the fs
		return false
	}

	return true
}

func ParseCLI() {
	var e error

	for idx, arg := range os.Args {
		//Skip first arg (process name)
		if idx>0 && arg[0] != '-' {
			DirToParse, e = filepath.Abs(arg)
			if e != nil {
				log.Fatal("filepath.Abs err:", e)
			}

			//TODO handle case where path to process is passed as last arg
			continue
		}

		arg = arg[1:]
		
		argComponents := strings.Split(arg, "=")
		argName := argComponents[0]
		argVal := ""
		if len(argComponents)>1 { argVal = argComponents[1] }

		switch argName {

		case "outfile":
			OutFileName = argVal

		case "generator":
			GeneratorFileName = "generate_godescribe.go"
			if len(argVal)>0 { GeneratorFileName = argVal } 
		}
	}
}

func MustValidateConfig() {
	if !PathExists(DirToParse) {
		log.Fatal("Could not access path:", DirToParse)
	}
}

func Usage() {
	fmt.Println(`godescribe [-flag1 ... -flagN]
Without arguments : godescribe will create a generator file to call itself to parse the local package and list functions.

-outfile=filename.go : will output json encoded symbols as a global var to the specified Go source file
-generator : will create a go generate generator to call itself at build time and create a symbols file named gogenerate_godescribe.go
-generator=filename.go same as previously with the ability to name the go generate source file used for storing symbols information.
`)
}

func PrintConfig() {
	fmt.Println("Dir to parse  :", DirToParse)
	fmt.Println("Generator file:", GeneratorFileName)
}

func main() {
	ParseCLI()
	MustValidateConfig()

	if len(DirToParse)>0 {
		symbolsJson, e := ParseSymbols(DirToParse)
		if e != nil {
			log.Fatal("Failed to parse symbols with err:", e)
		}

		if len(OutFileName) > 0 {
			gocode, e := GenGenerator(symbolsJson)

			if e = ioutil.WriteFile(OutFileName, []byte(gocode), 0644); e != nil {
				log.Fatal("Failed to write to output file: %s with err: %v", OutFileName, e)
			}

		} else {
			fmt.Println(symbolsJson)
		}

	// generator mode : simply create a source file containing a go generate directive
	// that will call godescribe to generate a symbols source file 
	} else if len(GeneratorFileName) > 0 {	
		generatorCode := `package main

//go:generate godescribe . -outfile=godescribe_symbols.go`

		if e := ioutil.WriteFile(GeneratorFileName, []byte(generatorCode), 0644); e != nil {
			log.Fatalf("Failed to create generator file: %s with err: %v", GeneratorFileName , e)
		}
	}
}
