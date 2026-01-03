package main

import (
	"expence_management/api"
	"expence_management/memory"
	"expence_management/pgsql"
	"expence_management/repo"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	EnvFilePath = ".env"
	DefaultDB   = "postgres"
)

var (
	repoFlag string
	options  *Config
)

// init function run before main function without need to be called
// here, we read configurations from config file and load environment variables
func init() {
	flag.StringVar(&repoFlag, "d", DefaultDB, "Which db?")
	flag.Parse()
	options = GetOptions()

	err := godotenv.Load(EnvFilePath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	var dbrepo *repo.DBRepository

	if repoFlag == "mem" {
		memrepo, err := memory.NewMemoryRepo()
		if err != nil {
			log.Fatalf("Error in connecting to memrepo: %s", err)
		}

		dbrepo = repo.NewDBResository(memrepo)
	} else {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Europe/Helsinki",
			options.Server.Host,
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			options.Database.Port,
			options.Database.SSLMode,
		)

		pgsqlRepo, err := pgsql.NewPgsqlRepo(dsn)
		if err != nil {
			log.Fatalf("Error in connecting to postgres: %s", err)
		}

		dbrepo = repo.NewDBResository(pgsqlRepo)
	}

	h := api.NewHandler(dbrepo)

	initDB(h)

	router := makeRouts(h)

	// Start server
	strport := strconv.Itoa(options.Server.Port)
	err := router.RunTLS(":"+strport, options.Server.Tls_Cert_Path, options.Server.Tls_Key_Path)
	if err != nil {
		log.Printf("Error in starting http server: %s", err)
	}
}
