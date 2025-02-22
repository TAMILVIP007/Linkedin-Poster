package main

import (
	"Linkedin-Poster/bot"
	"log"
)

func main() {
	if err := bot.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}
	bot.InitTgBot()
}
