package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

type PostData struct {
	Prompt  string
	PhotoID string
	Image   string
}

var postStore sync.Map

func InitTgBot() {
	b, err := gotgbot.NewBot(config.BotToken, nil)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("genpost", genpost))
	dispatcher.AddHandler(handlers.NewCommand("post", post))
	dispatcher.AddHandler(handlers.NewCommand("help", help))
	dispatcher.AddHandler(handlers.NewCommand("ai", geminiai))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("proceed."), proceedcbk))
	dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix("post."), postcbk))
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)
	updater.Idle()
}

func start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hey there, I'm %s! I can help you post to LinkedIn and process text with Gemini AI.", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	return err
}

func help(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hey there, I'm %s! I can help you post to LinkedIn and process text with Gemini AI.\n\nCommands:\n/genpost - Generate post\n/post - Post to LinkedIn\n/ai - Process text with Gemini AI", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
	})
	return err
}

func geminiai(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	args := ctx.Args()
	if ctx.EffectiveSender.Id() != config.OwnerID {
		_, err := msg.Reply(b, "You are not authorized to use this command", nil)
		return err
	}
	query := strings.Join(args[1:], " ")

	if query == "" {
		_, err := msg.Reply(b, "Send text to process", nil)
		return err
	}

	airesponse, err := ProcessGemini(query)
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("Failed to process text: %v", err), nil)
		return err
	}

	_, err = msg.Reply(b, airesponse, nil)
	return err
}

func genpost(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	args := ctx.Args()
	if ctx.EffectiveSender.Id() != config.OwnerID {
		_, err := msg.Reply(b, "You are not authorized to use this command", nil)
		return err
	}
	query := strings.Join(args[1:], " ")

	if query == "" {
		_, err := msg.Reply(b, "Send repo URL", nil)
		return err
	}

	var photoID, imagePath string
	if msg.ReplyToMessage != nil && len(msg.ReplyToMessage.Photo) > 0 {
		photoID = msg.ReplyToMessage.Photo[len(msg.ReplyToMessage.Photo)-1].FileId
		file, err := b.GetFile(photoID, nil)
		if err != nil {
			return err
		}
		imagePath, err = DownloadFile(file.URL(b, nil))
		if err != nil {
			return err
		}
	}

	uniqueID := fmt.Sprintf("%d", time.Now().UnixNano())
	prompt := GenPrompt(query)
	postStore.Store(uniqueID, PostData{
		Prompt:  prompt,
		PhotoID: photoID,
		Image:   imagePath,
	})

	_, err := msg.Reply(b, prompt, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Proceed", CallbackData: "proceed." + uniqueID},
			}},
		},
	})
	return err
}

func post(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	text := strings.Join(ctx.Args()[1:], " ")
	if ctx.EffectiveSender.Id() != config.OwnerID {
		_, err := msg.Reply(b, "You are not authorized to use this command", nil)
		return err
	}
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.Text != "" {
		text = msg.ReplyToMessage.Text
	}
	var photoID, imagePath string
	if len(msg.ReplyToMessage.Photo) > 0 {
		photoID = msg.ReplyToMessage.Photo[len(msg.ReplyToMessage.Photo)-1].FileId
		file, err := b.GetFile(photoID, nil)
		if err != nil {
			return err
		}
		imagePath, err = DownloadFile(file.URL(b, nil))
		if err != nil {
			return err
		}
	}

	if imagePath != "" {
		uploadUrl, asset, err := RegisterImageUpload()
		if err != nil {
			return err
		}
		defer os.Remove(imagePath)
		err = UploadImage(uploadUrl, imagePath)
		if err != nil {
			return err
		}
		link, err := PostToLinkedInWithImage(text, asset)
		if err != nil {
			return err
		}
		_, err = msg.Reply(b, fmt.Sprintf("Posted to LinkedIn: %s", link), nil)
		return err
	} else {
		link, err := PostToLinkedIn(text)
		if err != nil {
			return err
		}
		_, err = msg.Reply(b, fmt.Sprintf("Posted to LinkedIn: %s", link), nil)
		return err
	}
}

func proceedcbk(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	args := strings.Split(query.Data, ".")
	uniqueID := args[1]

	rawData, ok := postStore.Load(uniqueID)
	if !ok {
		_, _, err := query.Message.EditText(b, "Expired", nil)
		return err
	}

	postData := rawData.(PostData)
	query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Generating post...", ShowAlert: true})
	airesponse, err := ProcessGemini(postData.Prompt)
	if err != nil {
		return err
	}

	postStore.Store(uniqueID, PostData{
		Prompt:  airesponse,
		PhotoID: postData.PhotoID,
		Image:   postData.Image,
	})

	_, _, err = query.Message.EditText(b, airesponse, &gotgbot.EditMessageTextOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Post to LinkedIn", CallbackData: "post." + uniqueID},
			}},
		},
	})
	return err
}

func postcbk(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.Update.CallbackQuery
	args := strings.Split(query.Data, ".")
	uniqueID := args[1]

	rawData, ok := postStore.LoadAndDelete(uniqueID)
	if !ok {
		_, _, err := query.Message.EditText(b, "Expired", nil)
		return err
	}

	postData := rawData.(PostData)
	var link string
	var err error
	if postData.Image != "" {
		uploadUrl, asset, err := RegisterImageUpload()
		if err != nil {
			_, _, err := query.Message.EditText(b, fmt.Sprintf("Failed to upload image: %v", err), nil)
			return err
		}

		err = UploadImage(uploadUrl, postData.Image)
		defer os.Remove(postData.Image)
		if err != nil {
			_, _, err := query.Message.EditText(b, fmt.Sprintf("Failed to upload image: %v", err), nil)
			return err
		}

		link, err = PostToLinkedInWithImage(postData.Prompt, asset)
		if err != nil {
			_, _, err := query.Message.EditText(b, fmt.Sprintf("Failed to post to LinkedIn: %v", err), nil)
			return err
		}
	} else {
		link, err = PostToLinkedIn(postData.Prompt)
		if err != nil {
			_, _, err := query.Message.EditText(b, fmt.Sprintf("Failed to post to LinkedIn: %v", err), nil)
			return err
		}
	}

	query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{Text: "Posted to LinkedIn", ShowAlert: true})
	_, _, err = query.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
			{Text: "Link", Url: link},
		}},
	}})
	return err
}
