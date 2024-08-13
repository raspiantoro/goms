package gomsloader

import (
	"os"
	"reflect"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/raspiantoro/gormseeder"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type gomsBase interface {
	Path() string
}

type Seeds interface {
	gomsBase
}

type Migrations interface {
	gomsBase
}

type metaCollection []meta

func (f metaCollection) getSeeders(value reflect.Value) []*gormseeder.Seed {
	seeders := []*gormseeder.Seed{}

	for _, info := range f {
		seeders = append(seeders, &gormseeder.Seed{
			Key:      info.key,
			Name:     info.name,
			Seed:     intoSeederFunc(value.MethodByName("Up" + info.baseFuncName)),
			Rollback: intoSeederFunc(value.MethodByName("Down" + info.baseFuncName)),
		})
	}

	return seeders
}

func (f metaCollection) getMigrations(value reflect.Value) []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{}

	for _, info := range f {
		migrations = append(migrations, &gormigrate.Migration{
			ID:       info.key,
			Migrate:  intoMigrationUpFunc(value.MethodByName("Up" + info.baseFuncName)),
			Rollback: intoMigrationDownFunc(value.MethodByName("Down" + info.baseFuncName)),
		})
	}

	return migrations
}

type meta struct {
	key          string
	name         string
	baseFuncName string
}

func parseFilename(filename string) meta {
	nameSplit := strings.Split(filename, "_")
	key := nameSplit[1]
	name := strings.Replace(strings.Join(nameSplit[2:], " "), ".go", "", -1)

	var baseFuncName string

	for _, n := range nameSplit[2:] {
		baseFuncName += cases.Title(language.English, cases.Compact).String(n)
	}

	baseFuncName = strings.Replace(baseFuncName, ".go", "", -1)

	return meta{
		key:          key,
		name:         name,
		baseFuncName: baseFuncName,
	}
}

func load(g gomsBase) metaCollection {

	collection := metaCollection{}
	items, _ := os.ReadDir(g.Path())

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if item.Name() == "seeds.go" || item.Name() == "migrations.go" {
			continue
		}

		collection = append(collection, parseFilename(item.Name()))
	}

	return collection
}

func LoadSeeder(seeds Seeds) []*gormseeder.Seed {
	// seeders := []*gormseeder.Seed{}
	// items, _ := os.ReadDir(seeds.Path())

	// for _, item := range items {
	// 	if item.IsDir() {
	// 		continue
	// 	}

	// 	if item.Name() == "seeds.go" || item.Name() == "migrations.go" {
	// 		continue
	// 	}

	// 	key, name, baseFuncName := getAttrs(item.Name())

	// 	seedFuncName := "Seed" + baseFuncName
	// 	rollbackFuncName := "Rollback" + baseFuncName

	// 	value := reflect.ValueOf(seeds)

	// 	seeder := &gormseeder.Seed{
	// 		Key:      key,
	// 		Name:     name,
	// 		Seed:     intoSeederFunc(value.MethodByName(seedFuncName)),
	// 		Rollback: intoSeederFunc(value.MethodByName(rollbackFuncName)),
	// 	}

	// 	seeders = append(seeders, seeder)

	// }
	value := reflect.ValueOf(seeds)
	collection := load(seeds)

	return collection.getSeeders(value)
}

func LoadMigration(migrations Migrations) []*gormigrate.Migration {
	value := reflect.ValueOf(migrations)
	collection := load(migrations)

	return collection.getMigrations(value)
}

func intoSeederFunc(value reflect.Value) gormseeder.SeederFunc {
	f := value.Interface().(func(*gorm.DB) error)
	return gormseeder.SeederFunc(f)
}

func intoMigrationUpFunc(value reflect.Value) gormigrate.MigrateFunc {
	f := value.Interface().(func(*gorm.DB) error)
	return gormigrate.MigrateFunc(f)
}

func intoMigrationDownFunc(value reflect.Value) gormigrate.RollbackFunc {
	f := value.Interface().(func(*gorm.DB) error)
	return gormigrate.RollbackFunc(f)
}
