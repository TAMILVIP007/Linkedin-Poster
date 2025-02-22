# Telegram Bot for LinkedIn Posting with AI Integration

This project is a Telegram bot that helps users generate and post content to LinkedIn. The bot can:

- Generate AI-based responses using Gemini AI.
- Post AI-generated content or user-provided content (including images) to LinkedIn.
- Interact through Telegram inline commands.

The bot is designed to streamline content creation for LinkedIn by automating the process using AI and offering an easy way to post directly from Telegram.

## Features

- **AI Response Generation**: Generates LinkedIn post content using Gemini AI based on user prompts.
- **LinkedIn Integration**: Posts text and images directly to LinkedIn.
- **Image Uploads**: Handles LinkedIn image uploads seamlessly.
- **Interactive**: Provides inline buttons to confirm and proceed with posting.

---

## Prerequisites

1. **Telegram Bot Token**: You need to create a bot on Telegram and get a bot token from [BotFather](https://t.me/BotFather).
2. **LinkedIn API Credentials**: You'll need API keys from LinkedIn to enable posting.
3. **Gemini AI API Key**: Required to generate AI-based responses.

### How to Get LinkedIn API Credentials

1. Go to [LinkedIn Developer Portal](https://developer.linkedin.com/).
2. Create an application.
3. Once the application is created, navigate to the "Auth" section and note the following:
   - **Client ID**
   - **Client Secret**
4. Set up the permissions (`rw_organization_admin`, `w_member_social`, etc.) for posting to LinkedIn.
5. Use LinkedIn's OAuth 2.0 flow to get an access token for API interactions. This token will be used to post on behalf of users.



## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/TAMILVIP007/Linkedin-Poster
   cd Linkedin-Poster
   ```

2. **Install Dependencies**:

   Ensure you have Go installed on your system. Then install the required packages:

   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables**:

   Create a `.env` file with your configuration:

   ```bash
   BOT_TOKEN=<your-telegram-bot-token>
   LINKEDIN_CLIENT_ID=<your-linkedin-client-id>
   LINKEDIN_CLIENT_SECRET=<your-linkedin-client-secret>
   LINKEDIN_ACCESS_TOKEN=<linkedin-access-token>
   GEMINI_API_KEY=<your-gemini-api-key>
   ```

---

## Running the Bot

### Linux / macOS

1. Build and run the bot:

   ```bash
   go run main.go
   ```

2. Alternatively, build an executable:

   ```bash
   go build -o Linkedin-Poster
   ./Linkedin-Poster
   ```

### Windows

1. Open the command prompt or PowerShell and navigate to the project directory.
2. Run the bot:

   ```bash
   go run main.go
   ```

3. Or build an executable:

   ```bash
   go build -o Linkedin-Poster.exe
   .\Linkedin-Poster.exe
   ```

---

## Hosting Options

### Option 1: Local Hosting

Run the bot on your local machine by executing the steps in the "Running the Bot" section. Make sure your machine is always running to keep the bot online.

### Option 2: Cloud Hosting (e.g., AWS EC2, DigitalOcean, Heroku)

1. **AWS EC2 / DigitalOcean**:
   - Set up a virtual machine (VM).
   - SSH into your server and follow the installation and running steps above.

2. **Heroku**:
   - Set up a new Heroku project.
   - Use `heroku cli` to push your code.
   - Add environment variables in the Heroku dashboard.

### Option 3: Docker

1. **Build the Docker Image**:

   ```bash
   docker build -t Linkedin-Poster .
   ```

2. **Run the Docker Container**:

   ```bash
   docker run -d --name Linkedin-Poster -e BOT_TOKEN=<bot-token> -e LINKEDIN_CLIENT_ID=<client-id> -e LINKEDIN_CLIENT_SECRET=<client-secret> -e LINKEDIN_ACCESS_TOKEN=<access-token> -e GEMINI_API_KEY=<gemini-api-key> Linkedin-Poster
   ```

---

## Commands

- `/start`: Starts the bot and provides a welcome message.
- `/genpost <prompt>`: Generates AI-based content using the provided prompt.
- `/post <text>`: Posts a text or an image (if replying to an image) directly to LinkedIn.

---

## License

This project is licensed under the MIT License.
