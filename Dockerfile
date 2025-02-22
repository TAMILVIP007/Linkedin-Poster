# Use the official Go image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files first to download dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o Linkedin-Poster .

# Expose port (if necessary, depending on how your bot communicates)
# EXPOSE 8080

# Set environment variables (you can override these when running the container)
ENV BOT_TOKEN=<your-telegram-bot-token>
ENV LINKEDIN_CLIENT_ID=<your-linkedin-client-id>
ENV LINKEDIN_CLIENT_SECRET=<your-linkedin-client-secret>
ENV LINKEDIN_ACCESS_TOKEN=<linkedin-access-token>
ENV GEMINI_API_KEY=<your-gemini-api-key>

# Run the application
CMD ["./telegram-linkedin-bot"]
