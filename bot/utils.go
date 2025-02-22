package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenPrompt(repoLink string) string {
	resp, err := http.Get("https://raw.githubusercontent.com/" + strings.TrimPrefix(repoLink, "https://github.com/") + "/main/README.md")
	if err != nil {
		log.Fatalf("Error fetching README: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	lines := strings.Split(string(body), "\n")
	description := "A GitHub project demonstrating innovative use of modern technologies."
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		description = strings.TrimSpace(line)
		break
	}

	projectTitle := strings.Split(repoLink, "/")[len(strings.Split(repoLink, "/"))-1]

	return fmt.Sprintf(`
Generate a professional LinkedIn post for the project below. Follow this structure:

1. Simpler headline including project name on my behalf
2. Short project overview (2-3 sentences)
3. Key features/functionality (bullet points)
4. Technologies used (if mentioned)
5. Call-to-action for contributions or testing

Project: %s
Description: %s
Repository: %s

Include relevant emojis in headings and sections, clean formatting, and 3-5 hashtags.
`, projectTitle, description, repoLink)
}

func PostToLinkedIn(postText string) (string, error) {
	url := "https://api.linkedin.com/v2/ugcPosts"
	payload := map[string]interface{}{
		"author":         fmt.Sprintf("urn:li:person:%s", config.AuthorId),
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": postText,
				},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LinkedInToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to post to LinkedIn: %s", body)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return fmt.Sprintf("https://www.linkedin.com/feed/update/%s", result["id"]), nil
}

func RegisterImageUpload() (string, string, error) {
	url := "https://api.linkedin.com/v2/assets?action=registerUpload"
	payload := map[string]interface{}{
		"registerUploadRequest": map[string]interface{}{
			"owner":   fmt.Sprintf("urn:li:person:%s", config.AuthorId),
			"recipes": []string{"urn:li:digitalmediaRecipe:feedshare-image"},
			"serviceRelationships": []map[string]string{
				{
					"identifier":       "urn:li:userGeneratedContent",
					"relationshipType": "OWNER",
				},
			},
			"supportedUploadMechanism": []string{"SYNCHRONOUS_UPLOAD"},
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LinkedInToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	uploadUrl := result["value"].(map[string]interface{})["uploadMechanism"].(map[string]interface{})["com.linkedin.digitalmedia.uploading.MediaUploadHttpRequest"].(map[string]interface{})["uploadUrl"].(string)
	asset := result["value"].(map[string]interface{})["asset"].(string)

	return uploadUrl, asset, nil
}

func UploadImage(uploadUrl, imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, _ := http.NewRequest("POST", uploadUrl, file)
	req.Header.Set("Content-Type", "image/jpeg")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload image: %s", body)
	}

	return nil
}

func PostToLinkedInWithImage(postText, asset string) (string, error) {
	url := "https://api.linkedin.com/v2/ugcPosts"
	payload := map[string]interface{}{
		"author":         fmt.Sprintf("urn:li:person:%s", config.AuthorId),
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": postText,
				},
				"shareMediaCategory": "IMAGE",
				"media": []map[string]interface{}{
					{
						"status":      "READY",
						"description": map[string]interface{}{"text": "Image description"},
						"media":       asset,
						"title":       map[string]interface{}{"text": "Image Title"},
					},
				},
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.LinkedInToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to post to LinkedIn: %s", body)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return fmt.Sprintf("https://www.linkedin.com/feed/update/%s", result["id"]), nil
}

func ProcessGemini(text string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GeminiApiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	req := []genai.Part{genai.Text(text)}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, req...)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates found in response")
	}
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts found in response")
	}

	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
}

func DownloadFile(url string) (string, error) {
	downloadsDir := "downloads"
	if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
		err := os.Mkdir(downloadsDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create downloads directory: %v", err)
		}
	}
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	filePath := fmt.Sprintf("%s/%s", downloadsDir, fileName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return filePath, nil
}
