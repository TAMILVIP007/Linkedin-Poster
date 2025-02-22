package bot

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	LinkedInToken string
	OwnerID       int64
	BotToken      string
	GeminiApiKey  string
	AuthorId      string
}

var config *Config

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}
	ownerId, _ := strconv.ParseInt(os.Getenv("OWNER_ID"), 10, 64)
	config = &Config{
		LinkedInToken: os.Getenv("LINKEDIN_TOKEN"),
		OwnerID:       ownerId,
		BotToken:      os.Getenv("BOT_TOKEN"),
		AuthorId:      os.Getenv("AUTHOR_ID"),
		GeminiApiKey:  os.Getenv("GEMINI_API_KEY"),
	}

	return nil
}
