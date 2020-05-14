package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	httpPort := flag.Int("http.port", 4000, "HTTP Port to run server on")
	mongoDSN := flag.String("mongo.dsn", "localhost:27017", "DSN for mongoDB server")

	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appID := os.Getenv("PUSHER_APP_ID")
	appKey := os.Getenv("PUSHER_APP_KEY")
	appSecret := os.Getenv("PUSHER_APP_SECRET")
	appCluster := os.Getenv("PUSHER_APP_CLUSTER")
	appIsSecure := os.Getenv("PUSHER_APP_SECURE")
}
