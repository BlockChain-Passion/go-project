package main

import (
	"fmt"
	"log"

	"github.com/BlockChain-Passion/go-project/pkg/logger"
	"github.com/lpernett/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Env file missing", err)
	}
}

type gateway struct {
	Logger logger.ILogger
}

func main() {
	fmt.Println("Om namah shivay")
}
