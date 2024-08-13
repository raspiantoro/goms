package app

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/raspiantoro/goms/internal/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CreateSeed(config Config) {
	name := strings.ReplaceAll(config.Name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	config.Filename = "seed_" + time.Now().Format("20060102150405") + "_" + name + ".go"

	var baseFuncName string

	for _, n := range strings.Split(name, "_") {
		baseFuncName += cases.Title(language.English, cases.Compact).String(n)
	}

	searchPath := "."

	if config.Path != "" {
		searchPath = config.Path
	}

	path, err := scanGoms(searchPath, seederSubDir)
	if err != nil {
		log.Fatalln(err)
	}

	baseDir := filepath.Dir(path)

	mod := getModule()

	t, err := template.New("seed").Parse(string(templates.SeederTemplate()))
	if err != nil {
		log.Fatalln(err)
	}

	templateProps := map[string]string{
		"UpFuncName":      "Up" + baseFuncName,
		"DownFuncName":    "Down" + baseFuncName,
		"SeederSubDir":    seederSubDir,
		"SeedsStructName": cases.Title(language.English, cases.Compact).String(seederSubDir),
	}

	var b bytes.Buffer
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(mod.Dir+"/"+baseDir+"/"+config.Filename, b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
}

func CreateMigration(config Config) {
	name := strings.ReplaceAll(config.Name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	config.Filename = "migration_" + time.Now().Format("20060102150405") + "_" + name + ".go"

	var baseFuncName string

	for _, n := range strings.Split(name, "_") {
		baseFuncName += cases.Title(language.English, cases.Compact).String(n)
	}

	searchPath := "."

	if config.Path != "" {
		searchPath = config.Path
	}

	path, err := scanGoms(searchPath, migrationSubDir)
	if err != nil {
		log.Fatalln(err)
	}

	baseDir := filepath.Dir(path)

	mod := getModule()

	t, err := template.New("seed").Parse(string(templates.MigrationTemplate()))
	if err != nil {
		log.Fatalln(err)
	}

	templateProps := map[string]string{
		"UpFuncName":           "Up" + baseFuncName,
		"DownFuncName":         "Down" + baseFuncName,
		"MigrationsSubDir":     migrationSubDir,
		"MigrationsStructName": cases.Title(language.English, cases.Compact).String(migrationSubDir),
	}

	var b bytes.Buffer
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(mod.Dir+"/"+baseDir+"/"+config.Filename, b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
}

func scanGoms(path string, name string) (string, error) {
	var dir string
	var filenames []string

	structName := cases.Title(language.English, cases.Compact).String(name)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Base(path) != name+".go" {
			return nil
		}

		var filename string
		filename, errs := findStruct(path, structName)
		if errs != nil {
			return errs
		}

		if filename != "" {
			filenames = append(filenames, filename)
		}

		return nil
	})
	if err != nil {
		return dir, err
	}

	if len(filenames) == 0 {
		return dir, errors.New("goms files not found. you need to initialize first")
	}

	if len(filenames) > 1 {
		return dir, fmt.Errorf("multiple %s.go files with a %s struct found. provide the -d flag to specify your goms folder", name, structName)
	}

	dir = filenames[0]

	return dir, nil
}

func findStruct(path string, structName string) (string, error) {
	var filename string
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return "", err
	}

	ast.Inspect(f, func(n ast.Node) bool {

		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		if typeSpec.Name.Name == structName {
			filename = path
		}

		return true
	})

	return filename, nil
}
