package main

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Metadata struct {
	Type         string
	ResponseType string
}

type Tasks map[string]Metadata

func extractAutoGenTask(filename string) Tasks {
	tasks := Tasks{}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open %s, %v", filename, err)
	}
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if !strings.HasPrefix(line, "//autogen ") {
			break
		}
		fields := strings.Fields(line[10:])
		if len(fields) > 0 {
			structName := fields[0]
			md := Metadata{}
			if len(fields) > 1 {
				md.Type = fields[1]
				if len(fields) > 2 {
					md.ResponseType = fields[2]
				}
			}
			tasks[structName] = md
		}
	}
	f.Close()
	return tasks
}

type Context struct {
	pkgTypePrefix string
	depth         int
	counter       uint64
	r             bytes.Buffer
	w             bytes.Buffer
}

func (ctx *Context) unique() string {
	ctx.counter++
	return "v" + strconv.FormatUint(ctx.counter, 10)
}

func (ctx *Context) pad(w *bytes.Buffer) {
	write(w, strings.Repeat("\t", ctx.depth))
}

func (ctx *Context) both(s string) {
	ctx.read(s)
	ctx.write(s)
}

func (ctx *Context) read(s string) {
	for _, line := range strings.Split(s, "\n") {
		ctx.pad(&ctx.r)
		write(&ctx.r, line)
		write(&ctx.r, "\n")
	}
}

func (ctx *Context) write(s string) {
	for _, line := range strings.Split(s, "\n") {
		ctx.pad(&ctx.w)
		write(&ctx.w, line)
		write(&ctx.w, "\n")
	}
}

func (ctx *Context) readVar(t, v string) {
	ctx.read("if " + v + ", err = read" + t + "(r); err != nil {\n\treturn\n}")
}

func (ctx *Context) writeVar(t, v string) {
	ctx.write("if err = write" + t + "(w, " + v + "); err != nil {\n\treturn\n}")
}

func (ctx *Context) readTypedVar(t, v, r string) {
	tmp := ctx.unique()
	ctx.read("var " + tmp + " " + strings.ToLower(t))
	ctx.read("if " + tmp + ", err = read" + t + "(r); err != nil {\n\treturn\n}")
	ctx.read(v + " = " + r + "(" + tmp + ")")
}

func (ctx *Context) writeTypedVar(t, v string) {
	ctx.write("if err = write" + t + "(w, " + strings.ToLower(t) + "(" + v + ")); err != nil {\n\treturn\n}")
}

func (ctx *Context) readCall(s string) {
	ctx.read("if err = " + s + ".readFrom(r); err != nil {\n\treturn\n}")
}

func (ctx *Context) writeCall(s string) {
	ctx.write("if err = " + s + ".writeTo(w); err != nil {\n\treturn\n}")
}

func (ctx *Context) resolveType(t types.Type) (typeName string, realType types.Type) {
	typeName = ""
	realType = t
	for {
		named, ok := realType.(*types.Named)
		if !ok {
			break
		}
		if typeName == "" {
			typeName = named.String()
		}
		realType = named.Underlying()
	}
	if typeName == "" {
		switch realType.(type) {
		case *types.Basic:
			typeName = realType.(*types.Basic).Name()
		case *types.Slice:
			// typeName is resolved later
		default:
			log.Fatalf("unable to resolve type %#v", t)
		}
	} else {
		if strings.HasPrefix(typeName, ctx.pkgTypePrefix) {
			typeName = typeName[len(ctx.pkgTypePrefix):]
		}
		if !token.IsIdentifier(typeName) {
			log.Fatal("unexpected typeName: ", typeName)
		}
	}
	return
}

func (ctx *Context) Process(varName string, varType types.Type) {
	ctx.depth++
	ctx.both("// Field " + varName)
	varTypeName, varRealType := ctx.resolveType(varType)
	switch varRealType.(type) {
	case *types.Basic:
		basicType := varRealType.(*types.Basic)
		ctx.both("// Basic " + basicType.Name())
		switch basicType.Kind() {
		case types.Uint64:
			if varTypeName == "uint64" {
				ctx.readVar("Uint64", varName)
				ctx.writeVar("Uint64", varName)
			} else {
				ctx.both("// Typed " + varTypeName)
				ctx.readTypedVar("Uint64", varName, varTypeName)
				ctx.writeTypedVar("Uint64", varName)
			}
		case types.Uint32:
			if varTypeName == "uint32" {
				ctx.readVar("Uint32", varName)
				ctx.writeVar("Uint32", varName)
			} else {
				ctx.both("// Typed " + varTypeName)
				ctx.readTypedVar("Uint32", varName, varTypeName)
				ctx.writeTypedVar("Uint32", varName)
			}
		case types.Int32:
			if varTypeName == "int32" {
				ctx.readVar("Int32", varName)
				ctx.writeVar("Int32", varName)
			} else {
				ctx.both("// Typed " + varTypeName)
				ctx.readTypedVar("Int32", varName, varTypeName)
				ctx.writeTypedVar("Int32", varName)
			}
		case types.Uint8:
			if varTypeName == "uint8" || varTypeName == "byte" {
				ctx.read("if " + varName + ", err = r.ReadByte(); err != nil {\n\treturn\n}")
				ctx.write("if err = w.WriteByte(" + varName + "); err != nil {\n\treturn\n}")
			} else {
				ctx.both("// Typed " + varTypeName)
				tmp := ctx.unique()
				ctx.read("var " + tmp + " uint8")
				ctx.read("if " + tmp + ", err = " + "r.ReadByte(); err != nil {\n\treturn\n}")
				ctx.read(varName + " = " + varTypeName + "(" + tmp + ")")
				ctx.write("if err = w.WriteByte(uint8(" + varName + ")); err != nil {\n\treturn\n}")
			}
		case types.Bool:
			if varTypeName != "bool" {
				log.Fatal("typed bool is not supported")
			}
			ctx.readVar("Bool", varName)
			ctx.writeVar("Bool", varName)
		case types.String:
			if varTypeName != "string" {
				log.Fatal("typed string is not supported")
			}
			ctx.readVar("String", varName)
			ctx.writeVar("String", varName)
		default:
			log.Fatalf("basic type %#v is unsupported", basicType)
		}
	case *types.Struct:
		ctx.both("// Struct " + varTypeName)
		ctx.readCall(varName)
		ctx.writeCall(varName)
	case *types.Array:
		arrayType := varRealType.(*types.Array)
		if arrayType.Len() > 0x7FFFFFFF {
			log.Fatalf("array too long")
		}
		elemType := arrayType.Elem()
		elemTypeName, _ := ctx.resolveType(elemType)
		if elemTypeName == "" {
			log.Fatalf("array of slice is not supported directly")
		}
		length := strconv.FormatInt(arrayType.Len(), 10)
		ctx.both("// Array [" + length + "]" + elemTypeName)
		if elemBasicType, ok := elemType.(*types.Basic); ok && elemBasicType.Kind() == types.Byte {
			ctx.read("if err = readBytes(r, " + varName + "[:]); err != nil {\n\treturn\n}")
			ctx.write("if err = writeBytes(w, " + varName + "[:]); err != nil {\n\treturn\n}")
		} else {
			loopvar := ctx.unique()
			ctx.both("for " + loopvar + " := uint32(0); " + loopvar + " < " + length + "; " + loopvar + "++ {")
			ctx.Process(varName+"["+loopvar+"]", elemType)
			ctx.both("}")
		}
	case *types.Slice:
		sliceType := varRealType.(*types.Slice)
		elemType := sliceType.Elem()
		elemTypeName, _ := ctx.resolveType(elemType)
		if elemTypeName == "" {
			log.Fatalf("slice of slice is not supported directly")
		}
		tmp := ctx.unique()
		ctx.both("// Slice []" + elemTypeName)
		ctx.read("var " + tmp + " uint32")
		ctx.readVar("Uint32", tmp)
		ctx.read(varName + " = make([]" + elemTypeName + ", " + tmp + ")")
		ctx.write("if len(" + varName + ") > 0x7FFFFFFF {\n\terr = errTooLong\n\treturn\n}")
		ctx.write("var " + tmp + " = uint32(len(" + varName + "))")
		ctx.writeVar("Uint32", tmp)
		if elemBasicType, ok := elemType.(*types.Basic); ok && elemBasicType.Kind() == types.Byte {
			ctx.read("if err = readBytes(r, " + varName + "); err != nil {\n\treturn\n}")
			ctx.write("if err = writeBytes(w, " + varName + "); err != nil {\n\treturn\n}")
		} else {
			loopvar := ctx.unique()
			ctx.both("for " + loopvar + " := uint32(0); " + loopvar + " < " + tmp + "; " + loopvar + "++ {")
			ctx.Process(varName+"["+loopvar+"]", elemType)
			ctx.both("}")
		}
	default:
		log.Fatalf("unknown type %#v", varRealType)
	}
	ctx.depth--
}

func write(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		log.Fatal("unable to write: ", err)
	}
}

func writeFile(f *os.File, s string) {
	_, err := f.WriteString(s)
	if err != nil {
		log.Fatal("unable to write to ", f.Name(), ": ", err)
	}
}

func clearAutogenFiles() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range files {
		name := fi.Name()
		if strings.HasSuffix(name, "_autogen.go") {
			err := os.Remove(name)
			if err != nil {
				log.Fatal("unable to remove ", name)
			}
		}
	}
}

func main() {
	// clearAutogenFiles()

	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		".",
	)
	if err != nil {
		log.Fatal("packages.Load: ", err)
	}

	if len(pkgs) != 1 {
		log.Fatal("found multiple package")
	}

	var pkg *packages.Package
	for _, pkg = range pkgs {
		log.Print("found package: ", pkg)
		log.Print("package name: ", pkg.Name)
	}

	// if len(pkg.Errors) > 0 {
	//     for _, err := range pkg.Errors {
	//         log.Print(err)
	//     }
	//     log.Print("package has unresolved error(s)")
	// }

	allTasks := map[int]Tasks{}

	for idx, filename := range pkg.GoFiles {
		// _test.go are excluded
		if strings.HasSuffix(filename, "_autogen.go") {
			continue
		}
		tasks := extractAutoGenTask(filename)
		if len(tasks) > 0 {
			allTasks[idx] = tasks
		}
	}

	for idx, filename := range pkg.GoFiles {
		tasks, ok := allTasks[idx]
		if !ok {
			continue
		}
		outputFilename := filename[:len(filename)-3] + "_autogen.go"
		output, err := os.Create(outputFilename)
		if err != nil {
			log.Fatal("unable to create ", outputFilename)
		}
		writeFile(output, `// This file is automatically generated. DO NOT MODIFY!

package `+pkg.Name+`
`)
		for _, decl := range pkg.Syntax[idx].Decls {
			genDecl, _ := decl.(*ast.GenDecl)
			if genDecl == nil {
				continue
			}
			if genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				astStructType, _ := typeSpec.Type.(*ast.StructType)
				if astStructType == nil {
					continue
				}
				structName := typeSpec.Name.Name
				task, ok := tasks[structName]
				if !ok {
					continue
				}
				delete(tasks, structName)
				if task.Type != "" {
					writeFile(output, `
func (*`+structName+`) getType() PacketType {
	return `+task.Type+`
}
`)
				}
				if task.ResponseType != "" {
					writeFile(output, `
func (*`+structName+`) getResponseType() PacketType {
	return `+task.ResponseType+`
}
`)
				}

				log.Print("struct ", structName)
				ctx := Context{pkgTypePrefix: pkg.PkgPath + "."}
				structType := pkg.TypesInfo.TypeOf(astStructType).(*types.Struct)
				for i := 0; i < structType.NumFields(); i++ {
					field := structType.Field(i)
					if !field.Exported() {
						continue
					}
					ctx.Process("this."+field.Name(), field.Type())
				}

				writeFile(output, `
func (this *`+structName+`) readFrom(r Reader) (err error) {
`+ctx.r.String()+`	return
}
`)
				writeFile(output, `
func (this *`+structName+`) writeTo(w Writer) (err error) {
`+ctx.w.String()+`	return
}
`)
			}
		}
		if len(tasks) > 0 {
			log.Fatal("autogen struct is missing:\n%#v", tasks)
		}
		output.Close()
	}
}
