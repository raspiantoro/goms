package app

const (
	defaultPath   = "db"
	migrationsDir = "migrations"
	seedersDir    = "seeds"

	migrationNode = "Migrations"
	seederNode    = "Seeds"
)

type Config struct {
	Name     string
	Filename string
	Path     string
	Module   Module
}
