package app

const (
	defaultPath      = "db"
	migrationBaseDir = "migrator"
	migrationSubDir  = "migrations"
	seederBaseDir    = "seeder"
	seederSubDir     = "seeds"
)

type Config struct {
	Name     string
	Filename string
	Path     string
	// GomsDir  string
	Module Module
}
