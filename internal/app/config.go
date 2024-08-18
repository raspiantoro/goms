package app

const (
	defaultPath      = "db"
	migrationBaseDir = "migration"
	migrationSubDir  = "migrations"
	seederBaseDir    = "seeder"
	seederSubDir     = "seeds"

	migrationNode = "Migrations"
	seederNode    = "Seeds"
)

type Config struct {
	Name     string
	Filename string
	Path     string
	Module   Module
}
