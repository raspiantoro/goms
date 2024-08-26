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
	seederDir := fmt.Sprintf("%s/%s", cfg.Path, seedersDir)

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
		"SeederSubDir":    seedersDir,
		"SeedsStructName": cases.Title(language.English, cases.Compact).String(seedersDir),
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
		"SeederSubDir":   seedersDir,
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

	cliSeederDir := cfg.Module.Dir + "/" + filepath.Clean(filepath.Join(seederDir, "..")) + "/cli/seeder"
	if _, err := os.Stat(cliSeederDir); err != nil {
		err = os.MkdirAll(cliSeederDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := os.WriteFile(cliSeederDir+"/seeder.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	return nil
}

func initMigration(cfg Config) error {
	migrationDir := fmt.Sprintf("%s/%s", cfg.Path, migrationsDir)

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
		"MigrationsSubDir":     migrationsDir,
		"MigrationsStructName": cases.Title(language.English, cases.Compact).String(migrationsDir),
	}

	var b bytes.Buffer
	err = t.Execute(&b, templateProps)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(migrationDir+"/migrations.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	t, err = template.New("migration-cli").Parse(string(templates.CliMigrationTemplate()))
	if err != nil {
		log.Fatalln(err)
	}

	templateProps = map[string]string{
		"MigrationModuleName": cfg.Module.Path + "/" + migrationDir,
		"MigrationSubDir":     migrationsDir,
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

	cliMigrationDir := cfg.Module.Dir + "/" + filepath.Clean(filepath.Join(migrationDir, "..")) + "/cli/migrator"
	if _, err := os.Stat(cliMigrationDir); err != nil {
		err = os.MkdirAll(cliMigrationDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := os.WriteFile(cliMigrationDir+"/migrator.go", b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	return nil
}
