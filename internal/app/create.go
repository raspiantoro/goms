package app

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/raspiantoro/goms/internal/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CreateSeed(config Config) {
	name := strings.ReplaceAll(config.Name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	key := time.Now().Format("20060102150405")
	config.Filename = "seed_" + key + "_" + name + ".go"

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

	err = addSeedItem(
		mod.Dir+"/"+baseDir+"/seeds.go",
		seederNode,
		key,
		config.Name,
		templateProps["UpFuncName"],
		templateProps["DownFuncName"])
	if err != nil {
		log.Fatalln(err)
	}
}

func CreateMigration(config Config) {
	name := strings.ReplaceAll(config.Name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	key := time.Now().Format("20060102150405")
	config.Filename = "migration_" + key + "_" + name + ".go"

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

	err = addSeedItem(
		mod.Dir+"/"+baseDir+"/migrations.go",
		migrationNode,
		key,
		config.Name,
		templateProps["UpFuncName"],
		templateProps["DownFuncName"])
	if err != nil {
		log.Fatalln(err)
	}
}

func scanGoms(path string, name string) (string, error) {
	var dir string
	var filenames []string

	sliceName := cases.Title(language.English, cases.Compact).String(name)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Base(path) != name+".go" {
			return nil
		}

		var filename string
		filename, errs := findValueSpec(path, sliceName)
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
		return dir, fmt.Errorf("multiple %s.go files with a %s struct found. provide the -d flag to specify your goms folder", name, sliceName)
	}

	dir = filenames[0]

	return dir, nil
}

func findValueSpec(path string, identName string) (string, error) {
	var filename string
	fset := token.NewFileSet()
	f, err := decorator.ParseFile(fset, path, nil, 0)
	if err != nil {
		return "", err
	}

	dst.Inspect(f, func(n dst.Node) bool {

		valueSpec, ok := n.(*dst.ValueSpec)
		if !ok {
			return true
		}

		if valueSpec.Names[0].Name == identName {
			filename = path
		}

		return true
	})

	return filename, nil
}

func addSeedItem(filepath, searchNode, key, name, upFuncName, downFuncName string) error {
	fset := token.NewFileSet()
	file, err := decorator.ParseFile(fset, filepath, nil, 0)
	if err != nil {
		return err
	}

	dst.Inspect(file, func(n dst.Node) bool {
		switch x := n.(type) {
		case *dst.ValueSpec:
			if x.Names[0].Name != searchNode {
				return true
			}

			cl, ok := x.Values[0].(*dst.CompositeLit)
			if !ok {
				return true
			}

			var newChild *dst.CompositeLit

			if searchNode == seederNode {
				newChild = addSeedChildNode(key, name, upFuncName, downFuncName)
			} else {
				newChild = addMigrationChildNode(key, upFuncName, downFuncName)
			}

			cl.Elts = append(cl.Elts, newChild)

		}

		return true
	})

	stringBuf := strings.Builder{}

	decorator.Fprint(&stringBuf, file)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, []byte(stringBuf.String()), os.ModePerm)
}

func addSeedChildNode(key, name, seedFuncName, rollbackFuncName string) *dst.CompositeLit {
	newChild := new(dst.CompositeLit)

	newChild.Elts = []dst.Expr{
		createBasicLitKeyValue("Key", key),
		createBasicLitKeyValue("Name", name),
		createCallExprKeyValue("Seed", seedFuncName, "gormseeder", "SeederFunc"),
		createCallExprKeyValue("Rollback", rollbackFuncName, "gormseeder", "SeederFunc"),
	}

	newChild.Decs = dst.CompositeLitDecorations{
		NodeDecs: dst.NodeDecs{
			Before: dst.NewLine,
			After:  dst.NewLine,
		},
	}

	return newChild
}

func addMigrationChildNode(key, seedFuncName, rollbackFuncName string) *dst.CompositeLit {
	newChild := new(dst.CompositeLit)

	newChild.Elts = []dst.Expr{
		createBasicLitKeyValue("ID", key),
		createCallExprKeyValue("Migrate", seedFuncName, "gormigrate", "MigrateFunc"),
		createCallExprKeyValue("Rollback", rollbackFuncName, "gormigrate", "RollbackFunc"),
	}

	newChild.Decs = dst.CompositeLitDecorations{
		NodeDecs: dst.NodeDecs{
			Before: dst.NewLine,
			After:  dst.NewLine,
		},
	}

	return newChild
}

func createBasicLitKeyValue(key, value string) *dst.KeyValueExpr {
	return &dst.KeyValueExpr{
		Key: &dst.Ident{
			Name: key,
			Decs: dst.IdentDecorations{
				NodeDecs: dst.NodeDecs{
					Before: dst.NewLine,
					After:  dst.NewLine,
				},
			},
		},
		Value: &dst.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", value),
			Decs: dst.BasicLitDecorations{
				NodeDecs: dst.NodeDecs{
					Before: dst.NewLine,
					After:  dst.NewLine,
				},
			},
		},
	}
}

func createCallExprKeyValue(key, value, pkg, funcName string) *dst.KeyValueExpr {
	return &dst.KeyValueExpr{
		Key: &dst.Ident{
			Name: key,
		},
		Value: &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X: &dst.Ident{
					Name: pkg,
				},
				Sel: &dst.Ident{
					Name: funcName,
				},
			},
			Args: []dst.Expr{
				&dst.Ident{
					Name: value,
				},
			},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.NewLine,
				After:  dst.NewLine,
			},
		},
	}
}
