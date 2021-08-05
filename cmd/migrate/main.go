package main

import (
	"Tradeasy/config"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pressly/goose"
	"log"
	"os"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "../../migration", "directory with migration files")
	err error
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}
	command := args[0] // command like up, down

	dbstring := config.DbURL(config.BuildConfig())

	//assign connection to global *gorm.DB variable DB
	config.DB,err = gorm.Open("mysql", dbstring)
	if err != nil {
		log.Fatalf("Gorm: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := config.DB.Close(); err != nil {
			log.Fatalf("Gorm: failed to close DB: %v\n", err)
		}
	}()


	// running goose commands
	db,err_ := goose.OpenDBWithDriver("mysql",dbstring)
	if err_ != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}
	defer func() {
		if err_ := db.Close(); err_ != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()


	if err_ := goose.Run(command, db, *dir); err_ != nil {
		log.Fatalf("goose %v: %v", command, err)
	}

}
func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status
Options:
`
	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)
