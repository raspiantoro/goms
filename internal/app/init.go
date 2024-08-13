package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/raspiantoro/goms/internal/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Init(config Config) {
	if err := os.Setenv("GOWORK", "off"); err != nil {
		log.Fatalln(err)
	}
	defer os.Setenv("GOWORK", "on")

	mod := getModule()

	if mod.Path == "command-line-arguments" {
		fmt.Println("ERROR: `go.mod` file not found in the current directory")
		return
	}

	config.Module = mod

	if config.Path == "" {
		config.Path = defaultPath
	}

	err := initSeeder(config)
	if err != nil {
		log.Fatalln(err)
	}

	err = initMigration(config)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("go", "mod", "tidy")

	err = cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func initSeeder(cfg Config) error {
	seederDir := fmt.Sprintf("%s/%s/%s", cfg.Path, seederBaseDir, seederSubDir)

	if _, err := os.Stat(seederDir); err != nil {
		err = os.MkdirAll(seederDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	t, err := template.New("seeds").Parse(string(templates.SeedTemplate()))
	if err != nil {
		return err
	}

	templateProps := map[string]string{
		"SeederSubDir":    seederSubDir,
		"SeedsStructName": cases.Title(language.English, cases.Compact).String(seederSubDir),
	}

	var b bytes.Buffer
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(seederDir+"/seeds.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	t, err = template.New("seeder-cli").Parse(string(templates.CliSeederTemplate()))
	if err != nil {
		log.Fatalln(err)
	}

	templateProps = map[string]string{
		"SeedModuleName": cfg.Module.Path + "/" + seederDir,
		"SeederSubDir":   seederSubDir,
		"Host":           "{{database-host}}",
		"Username":       "{{database-username}}",
		"Password":       "{{database-password}}",
		"DbName":         "{{database-name}}",
		"Port":           "{{database-port}}",
	}

	b.Reset()
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(cfg.Module.Dir+"/"+filepath.Clean(filepath.Join(seederDir, ".."))+"/seeder.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	return nil
}

func initMigration(cfg Config) error {
	migrationDir := fmt.Sprintf("%s/%s/%s", cfg.Path, migrationBaseDir, migrationSubDir)

	if _, err := os.Stat(migrationDir); err != nil {
		err = os.MkdirAll(migrationDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	t, err := template.New("migrations").Parse(string(templates.MigrationsTemplate()))
	if err != nil {
		return err
	}

	templateProps := map[string]string{
		"MigrationsSubDir":     migrationSubDir,
		"MigrationsStructName": cases.Title(language.English, cases.Compact).String(migrationSubDir),
	}

	var b bytes.Buffer
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(migrationDir+"/migrations.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	t, err = template.New("seeder-cli").Parse(string(templates.CliMigrationTemplate()))
	if err != nil {
		log.Fatalln(err)
	}

	templateProps = map[string]string{
		"MigrationModuleName": cfg.Module.Path + "/" + migrationDir,
		"MigrationSubDir":     migrationSubDir,
		"Host":                "{{database-host}}",
		"Username":            "{{database-username}}",
		"Password":            "{{database-password}}",
		"DbName":              "{{database-name}}",
		"Port":                "{{database-port}}",
	}

	b.Reset()
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(cfg.Module.Dir+"/"+filepath.Clean(filepath.Join(migrationDir, ".."))+"/migrator.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	return nil
}
