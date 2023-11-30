package shared

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	MongoUri                  string
	MongoDatabaseName         string
	MongoEventsCollectionName string
	GithubEventsUrl           string
	Role                      string
	ApiPort                   int
	WriteDBWorkers            int
	GithubRepoWorkers         int
)

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println("Successfully load .env file.")
}

func init() {
	initEnv()
	var err error
	Role = os.Getenv("ROLE")

	MongoUri = os.Getenv("MONGODB_URI")
	if MongoUri == "" {
		log.Fatal("MONGODB_URI cant be empty.")
	}

	MongoDatabaseName = os.Getenv("MONGODB_DATABASE_NAME")
	if MongoDatabaseName == "" {
		log.Fatal("MONGODB_DATABASE_NAME cant be empty.")
	}

	MongoEventsCollectionName = os.Getenv("MONGODB_EVENTS_COLLECTION_NAME")
	if MongoEventsCollectionName == "" {
		log.Fatal("MONGODB_EVENTS_COLLECTION_NAME cant be empty.")
	}

	GithubEventsUrl = os.Getenv("GITHUB_EVENTS_URL")
	if GithubEventsUrl == "" && Role != "API" {
		log.Fatal("GITHUB_EVENTS_URL cant be empty.")
	}

	ApiPort, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Printf("API_PORT set to default, port %v", DefaultApiPort)
		ApiPort = DefaultApiPort
	}

	WriteDBWorkers, err = strconv.Atoi(os.Getenv("WRITE_DB_WORKERS"))
	if err != nil && Role != ApiRoleName {
		log.Fatal(err)
	}
	GithubRepoWorkers, err = strconv.Atoi(os.Getenv("GITHUB_REPO_WORKERS"))
	if err != nil && Role != ApiRoleName {
		log.Fatal(err)
	}
}
