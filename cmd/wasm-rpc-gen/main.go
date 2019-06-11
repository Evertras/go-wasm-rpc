package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type RequiredFunction struct {
	Name    string
	ReqType string
	ResType string
}

type TemplateInput struct {
	PackageExport      string
	PackageImportShort string
	PackageImportFull  string
	Methods            []RequiredFunction
}

const outputTemplateText = `// This file is generated by go-wasm-rpc.  DO NOT EDIT.

package <<.PackageExport>>

import (
	"context"
	"syscall/js"

	proto "github.com/golang/protobuf/proto"

	"<<.PackageImportFull>>"
)

// Note: https://github.com/golang/go/commit/c468ad04177c422534ad1ed4547295935f84743d
// will make this much nicer... this is stolen from https://github.com/golang/go/issues/31335
// and is absolutely awful for performance, but okay for proof of concept for now
func typedArrayToByteSlice(arg js.Value) []byte {
	length := arg.Length()
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(arg.Index(i).Int())
	}
	return bytes
}

var server <<.PackageImportShort>>.WasmServiceServer = &wasmServer{}

func RegisterWasmCallbacks(base js.Value) {
	base.Set("ready", js.ValueOf(true))
<<range .Methods>>	base.Set("<<.Name>>", js.FuncOf(<<.Name>>))
<<end>>}
<<$packageImportShort := .PackageImportShort>><<range .Methods>>
func <<.Name>>(this js.Value, args []js.Value) interface{} {
	go func() {
		msg := &<<$packageImportShort>>.<<.ReqType>>{}

		// This is awful and should feel awful, but proof of concept... (see above)
		byteSlice := typedArrayToByteSlice(args[0])

		err := proto.Unmarshal(byteSlice, msg)

		if err != nil {
			args[1].Invoke(js.ValueOf(err.Error()))
		}

		res, err := server.<<.Name>>(context.Background(), msg)

		if err != nil {
			args[1].Invoke(js.ValueOf(err.Error()))
		}

		marshaled, err := proto.Marshal(res)

		if err != nil {
			args[1].Invoke(js.ValueOf(err.Error()))
		}

		returnedArray := js.TypedArrayOf(marshaled)
		args[1].Invoke(js.Undefined(), returnedArray)
		returnedArray.Release()
	}()

	return nil
}
<<end>>
`

func main() {
	flag.Parse()

	if flag.NArg() != 4 {
		fmt.Println("Usage: wasm-rpc-gen <input-file> <output-file> <import-package-name> <export-package-name>")
		fmt.Println("       <input-file> should be the generated protobuf")
		fmt.Println("                    .go file that includes your service")
		fmt.Println()
		fmt.Println("       <output-file> should be in the same directory as")
		fmt.Println("                     your service implementation")
		fmt.Println()
		fmt.Println("       <import-package-name> should be the *full* package name of")
		fmt.Println("                             the input generated proto package")
		fmt.Println()
		fmt.Println("       <export-package-name> should be the *short* package name of")
		fmt.Println("                             the output directory")
		os.Exit(1)
	}

	inputFilename := flag.Arg(0)
	outputFilename := flag.Arg(1)
	packageImportName := flag.Arg(2)
	packageExportName := flag.Arg(3)

	packageImportSplit := strings.Split(packageImportName, "/")

	packageImportShort := packageImportSplit[len(packageImportSplit)-1]
	fmt.Println(packageImportShort)

	outputTemplate := template.Must(template.New("interface").Delims("<<", ">>").Parse(outputTemplateText))
	src, err := ioutil.ReadFile(inputFilename)

	if err != nil {
		fmt.Println("Error reading file")
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputFilename, src, 0)

	if err != nil {
		fmt.Println("Error parsing file")
		panic(err)
	}

	input := TemplateInput{
		PackageExport:      packageExportName,
		PackageImportFull:  packageImportName,
		PackageImportShort: packageImportShort,
		Methods:            make([]RequiredFunction, 0),
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch outer := n.(type) {
		case *ast.TypeSpec:
			if outer.Name.Name == "WasmServiceServer" {
				i := outer.Type.(*ast.InterfaceType)

				for _, method := range i.Methods.List {
					fn := method.Type.(*ast.FuncType)
					reqType := fn.Params.List[1].Type.(*ast.StarExpr)
					resType := fn.Results.List[0].Type.(*ast.StarExpr)

					input.Methods = append(input.Methods, RequiredFunction{
						Name:    method.Names[0].Name,
						ReqType: fmt.Sprintf("%s", reqType.X),
						ResType: fmt.Sprintf("%s", resType.X),
					})
				}

				return false
			}
		}

		return true
	})

	fout, err := os.Create(outputFilename)

	if err != nil {
		fmt.Println("Error creating file", outputFilename)
		panic(err)
	}

	defer fout.Close()

	err = outputTemplate.Execute(fout, input)

	if err != nil {
		fmt.Println("Error running template")
		panic(err)
	}

	fmt.Println("Generated successfully to", outputFilename)
}
