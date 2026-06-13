package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main()  {

	//load .env file
	if err := godotenv.Load(); err != nil{
		log.Fatal(" error loading .env file")
		
	}

	// connect to database
	ctx := context.
	
}

godotenv.Load() // reads .env into environment variables