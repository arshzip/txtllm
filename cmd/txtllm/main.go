package main

import (
	"fmt"
	"log"

	"github.com/arshzip/txtllm/internal/dns"
	"github.com/arshzip/txtllm/internal/llm"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to txtllm!")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	client, err := llm.NewOpenRouterClient()
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	dns.Start(client)
}
