package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/radifan9/social-media-backend/internal/configs"
	"github.com/radifan9/social-media-backend/internal/routers"
)

func main() {
	log.Println("--- --- Social Media --- ---")

	// Load environment variables
	if err := godotenv.Load(""); err != nil {
		log.Println("failed to load environment variables\nCause: ", err.Error())
		return
	}

	// PostgreSQL DB Initialization
	db, err := configs.InitDB()
	if err != nil {
		log.Println("failed to connect to database\nCause: ", err.Error())
		return
	}
	defer db.Close()

	// Test DB Connection
	if err := configs.TestDB(db); err != nil {
		log.Println("ping to DB failed\nCause: ", err.Error())
		return
	}
	log.Println("✅ PostgreSQL connected.")

	// Redis Initialization
	rdb := configs.InitRDB()
	defer rdb.Close()

	// Test Redis Connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Println("failed to ping redis database\nCause: ", err.Error())
		return
	}
	log.Println("✅ Successfully connect & ping to rdb!")

	// Engine Gin Initialization
	router := routers.InitRouter(db, rdb)
	router.Run(":8080")

	// Flow of the program
	// client => (router => handler => repo => handler) => client
}
