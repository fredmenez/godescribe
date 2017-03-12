## godescribe

godescribe is a command line tool for Go source code that parses a directory and generates a list of exported functions signatures in the JSON format.

Go functions names and signature are encoded in a simple JSON model : functionName : functionSignature, eg.:

```
{
  "GetSensorValue" : {
  	"name"  : "GetSensorValue",
  	"params": [
  			{ "name" : "idx", "type":"uint32" }
  		],
  	"results": [
  			{ "type":"float32" }
  		]
  	}
}
```

## Usage

godescribe dirPath [-outfile=filename.go] [-generator[=filename.go]]

-outfile=filename.go : will create filename.go with a variable named `SymbolsJson` assigned with the JSON representation string.

-generator : generates a source file named `generate_godescribe.go` containing a go generate directive to call godescribe for the current package.

-generator=filename.go : same as -generator but will output to the specified filename.

## Examples

parse current dir and print to standard output :
```
godescribe .
```

parse directory and output results to a Go source file :
```
godescribe /path/to/package -outfile=symbols.go
```

## Installation

```
go get github.com/fredmenez/godescribe
go install github.com/fredmenez/godescribe
```

## License

This tool is licensed under the [3-Clause BSD License](http://opensource.org/licenses/BSD-3-Clause). See [license](LICENSE) file.

## Status

"Early days" : some changes can be expected based on feedback and new features.

