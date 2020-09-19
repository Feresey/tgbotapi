package tgapi

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if !flag.Parsed() {
		flag.Parse()
	}
	if testing.Short() {
		return
	}

	os.Exit(m.Run())
}

var ctx = context.Background()

func TestGetMe(t *testing.T) {
	me, _, err := api.GetMe(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, me)
}

func TestGetCommands(t *testing.T) {
	_, _, err := api.GetMyCommands(ctx)
	require.NoError(t, err)
}

const (
	TestToken               = "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I" //nolint:gosec // this is not my credentials
	ChatID                  = 76918703
	SupergroupChatID        = -1001120141283
	ReplyToMessageID        = 35
	ExistingPhotoFileID     = "AgADAgADw6cxG4zHKAkr42N7RwEN3IFShCoABHQwXEtVks4EH2wBAAEC"
	ExistingDocumentFileID  = "BQADAgADOQADjMcoCcioX1GrDvp3Ag"
	ExistingAudioFileID     = "BQADAgADRgADjMcoCdXg3lSIN49lAg"
	ExistingVoiceFileID     = "AwADAgADWQADjMcoCeul6r_q52IyAg"
	ExistingVideoFileID     = "BAADAgADZgADjMcoCav432kYe0FRAg"
	ExistingVideoNoteFileID = "DQADAgADdQAD70cQSUK41dLsRMqfAg"
	ExistingStickerFileID   = "BQADAgADcwADjMcoCbdl-6eB--YPAg"
)

func newInt(i int64) *int64      { return &i }
func newString(s string) *string { return &s }
func newBool(b bool) *bool       { return &b }

var api = New(TestToken)

func TestGetUpdates(t *testing.T) {
	_, _, err := api.GetUpdates(ctx, nil)
	require.NoError(t, err)
}

func TestSendWithMessage(t *testing.T) {
	msg := &SendMessageConfig{
		ChatID:    IntStr{Int: ChatID},
		Text:      "A test message from the test library in telegram-bot-api",
		ParseMode: newString(ParseModeMarkdown),
	}
	resp, _, err := api.SendMessage(ctx, msg)
	require.NoError(t, err)

	require.NotEmpty(t, resp.Text)
	require.Equal(t, msg.Text, *resp.Text)
}

func TestSendWithMessageReply(t *testing.T) {
	msg := &SendMessageConfig{
		ChatID:           IntStr{Int: ChatID},
		Text:             "A test message from the test library in telegram-bot-api",
		ReplyToMessageID: newInt(ReplyToMessageID),
	}
	resp, _, err := api.SendMessage(ctx, msg)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Text)
	require.Equal(t, msg.Text, *resp.Text)
	require.NotEmpty(t, resp.ReplyToMessage)
	require.Equal(t, msg.ReplyToMessageID, resp.ReplyToMessage.MessageID)
}

func TestSendWithMessageForward(t *testing.T) {
	msg := &ForwardMessageConfig{
		ChatID:     IntStr{Int: ChatID},
		FromChatID: IntStr{Int: ChatID},
		MessageID:  ReplyToMessageID,
	}
	_, _, err := api.ForwardMessage(ctx, msg)
	require.NoError(t, err)
}

func TestDeleteMessage(t *testing.T) {
	msg := &SendMessageConfig{
		ChatID:    IntStr{Int: ChatID},
		Text:      "A test message from the test library in telegram-bot-api",
		ParseMode: newString(ParseModeMarkdown),
	}
	message, _, err := api.SendMessage(ctx, msg)
	require.NoError(t, err)

	_, err = api.DeleteMessage(ctx, msg.ChatID, message.MessageID)
	require.NoError(t, err)
}

func TestPin(t *testing.T) {
	msg := &SendMessageConfig{
		ChatID:    IntStr{Int: SupergroupChatID},
		Text:      "A test message from the test library in telegram-bot-api",
		ParseMode: newString(ParseModeMarkdown),
	}

	message, _, err := api.SendMessage(ctx, msg)
	require.NoError(t, err)

	t.Run("pin", func(t *testing.T) {
		_, err = api.PinChatMessage(ctx, &PinChatMessageConfig{
			ChatID:    msg.ChatID,
			MessageID: message.MessageID,
		})
		require.NoError(t, err)
	})

	t.Run("unpin", func(t *testing.T) {
		_, err := api.UnpinChatMessage(ctx, msg.ChatID)
		require.NoError(t, err)
	})
}

// func TestSendWithNewPhoto(t *testing.T) {
// 	msg := NewPhotoUpload(ChatID, "tests/image.jpg")
// 	msg.Caption = "Test"
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewPhotoWithFileBytes(t *testing.T) {
// 	data, _ := ioutil.ReadFile("tests/image.jpg")
// 	b := FileBytes{Name: "image.jpg", Bytes: data}

// 	msg := NewPhotoUpload(ChatID, b)
// 	msg.Caption = "Test"
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewPhotoWithFileReader(t *testing.T) {
// 	f, _ := os.Open("tests/image.jpg")
// 	reader := FileReader{Name: "image.jpg", Reader: f, Size: -1}

// 	msg := NewPhotoUpload(ChatID, reader)
// 	msg.Caption = "Test"
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewPhotoReply(t *testing.T) {
// 	msg := NewPhotoUpload(ChatID, "tests/image.jpg")
// 	msg.ReplyToMessageID = ReplyToMessageID

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingPhoto(t *testing.T) {
// 	msg := NewPhotoShare(ChatID, ExistingPhotoFileID)
// 	msg.Caption = "Test"
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewDocument(t *testing.T) {
// 	msg := NewDocumentUpload(ChatID, "tests/image.jpg")
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingDocument(t *testing.T) {
// 	msg := NewDocumentShare(ChatID, ExistingDocumentFileID)
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewAudio(t *testing.T) {
// 	msg := NewAudioUpload(ChatID, "tests/audio.mp3")
// 	msg.Title = "TEST"
// 	msg.Duration = 10
// 	msg.Performer = "TEST"
// 	msg.MimeType = "audio/mpeg"
// 	msg.FileSize = 688
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingAudio(t *testing.T) {
// 	msg := NewAudioShare(ChatID, ExistingAudioFileID)
// 	msg.Title = "TEST"
// 	msg.Duration = 10
// 	msg.Performer = "TEST"

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewVoice(t *testing.T) {
// 	msg := NewVoiceUpload(ChatID, "tests/voice.ogg")
// 	msg.Duration = 10
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingVoice(t *testing.T) {
// 	msg := NewVoiceShare(ChatID, ExistingVoiceFileID)
// 	msg.Duration = 10
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithContact(t *testing.T) {
// 	contact := NewContact(ChatID, "5551234567", "Test")

// 	if resp, _, err := api.SendMessage(contact); err != nil {
// 		t.Error(err)
// 		t.Fail()
// 	}
// }

// func TestSendWithLocation(t *testing.T) {
// 	resp, _, err := api.SendMessage(NewLocation(ChatID, 40, 40))

// 	require.NoError(t, err)
// }

// func TestSendWithVenue(t *testing.T) {
// 	venue := NewVenue(ChatID, "A Test Location", "123 Test Street", 40, 40)

// 	if resp, _, err := api.SendMessage(venue); err != nil {
// 		t.Error(err)
// 		t.Fail()
// 	}
// }

// func TestSendWithNewVideo(t *testing.T) {
// 	msg := NewVideoUpload(ChatID, "tests/video.mp4")
// 	msg.Duration = 10
// 	msg.Caption = "TEST"

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingVideo(t *testing.T) {
// 	msg := NewVideoShare(ChatID, ExistingVideoFileID)
// 	msg.Duration = 10
// 	msg.Caption = "TEST"

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewVideoNote(t *testing.T) {
// 	msg := NewVideoNoteUpload(ChatID, 240, "tests/videonote.mp4")
// 	msg.Duration = 10

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingVideoNote(t *testing.T) {
// 	msg := NewVideoNoteShare(ChatID, 240, ExistingVideoNoteFileID)
// 	msg.Duration = 10

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewSticker(t *testing.T) {
// 	msg := NewStickerUpload(ChatID, "tests/image.jpg")

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingSticker(t *testing.T) {
// 	msg := NewStickerShare(ChatID, ExistingStickerFileID)

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithNewStickerAndKeyboardHide(t *testing.T) {
// 	msg := NewStickerUpload(ChatID, "tests/image.jpg")
// 	msg.ReplyMarkup = ReplyKeyboardRemove{
// 		RemoveKeyboard: true,
// 		Selective:      false,
// 	}
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithExistingStickerAndKeyboardHide(t *testing.T) {
// 	msg := NewStickerShare(ChatID, ExistingStickerFileID)
// 	msg.ReplyMarkup = ReplyKeyboardRemove{
// 		RemoveKeyboard: true,
// 		Selective:      false,
// 	}

// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)
// }

// func TestSendWithDice(t *testing.T) {
// 	msg := NewDice(ChatID)
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)

// }

// func TestSendWithDiceWithEmoji(t *testing.T) {
// 	msg := NewDiceWithEmoji(ChatID, "🏀")
// 	resp, _, err := api.SendMessage(msg)

// 	require.NoError(t, err)

// }

// func TestGetFile(t *testing.T) {
// 	file := FileConfig{FileID: ExistingPhotoFileID}

// 	_, err := api.GetFile(file)

// 	require.NoError(t, err)
// }

// func TestSendChatConfig(t *testing.T) {
// 	resp, _, err := api.SendMessage(NewChatAction(ChatID, ChatTyping))

// 	require.NoError(t, err)
// }

// func TestSendEditMessage(t *testing.T) {
// 	msg, err := api.SendMessage(NewMessage(ChatID, "Testing editing."))
// 	require.NoError(t, err)

// 	edit := EditMessageTextConfig{
// 		BaseEdit: BaseEdit{
// 			ChatID:    ChatID,
// 			MessageID: msg.MessageID,
// 		},
// 		Text: "Updated text.",
// 	}

// 	_, err = api.SendMessage(edit)
// 	require.NoError(t, err)
// }

// func TestGetUserProfilePhotos(t *testing.T) {
// 	_, err := api.GetUserProfilePhotos(NewUserProfilePhotos(ChatID))
// 	require.NoError(t, err)
// }

// func TestSetWebhookWithCert(t *testing.T) {
// 	time.Sleep(time.Second * 2)

// 	_, err := api.RemoveWebhook()
// 	require.NoError(t, err)

// 	wh := NewWebhookWithCert("https://example.com/tgbotapi-test/"+api.Token, "tests/cert.pem")
// 	_, err = api.SetWebhook(wh)
// 	require.NoError(t, err)
// 	_, err = api.GetWebhookInfo()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = api.RemoveWebhook()
// 	require.NoError(t, err)
// }

// func TestSetWebhookWithoutCert(t *testing.T) {
// 	time.Sleep(time.Second * 2)

// 	_, err := api.RemoveWebhook()
// 	require.NoError(t, err)

// 	wh := NewWebhook("https://example.com/tgbotapi-test/" + api.Token)
// 	_, err = api.SetWebhook(wh)
// 	require.NoError(t, err)
// 	info, err := api.GetWebhookInfo()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if info.MaxConnections == 0 {
// 		t.Errorf("Expected maximum connections to be greater than 0")
// 	}
// 	if info.LastErrorDate != 0 {
// 		t.Errorf("[Telegram callback failed]%s", info.LastErrorMessage)
// 	}
// 	_, err = api.RemoveWebhook()
// 	require.NoError(t, err)
// }

// func TestUpdatesChan(t *testing.T) {
// 	var ucfg UpdateConfig = NewUpdate(0)
// 	ucfg.Timeout = 60
// 	_, err := api.GetUpdatesChan(ucfg)

// 	require.NoError(t, err)
// }

// func TestSendWithMediaGroup(t *testing.T) {
// 	cfg := NewMediaGroup(ChatID, []interface{}{
// 		NewInputMediaPhoto("https://i.imgur.com/unQLJIb.jpg"),
// 		NewInputMediaPhoto("https://i.imgur.com/J5qweNZ.jpg"),
// 		NewInputMediaVideo("https://i.imgur.com/F6RmI24.mp4"),
// 	})
// 	resp, _, err := api.SendMessage(cfg)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

/*

// func ExampleNewBotAPI() {
// 	bot, err := New("MyAwesomeBotToken")
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	log.Printf("Authorized on account %s", api.Self.UserName)

// 	u := NewUpdate(0)
// 	u.Timeout = 60

// 	updates, err := api.GetUpdatesChan(u)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Optional: wait for updates and clear them if you don't want to handle
// 	// a large backlog of old messages
// 	time.Sleep(time.Millisecond * 500)
// 	updates.Clear()

// 	for update := range updates {
// 		if update.Message == nil {
// 			continue
// 		}

// 		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

// 		msg := NewMessage(update.Message.Chat.ID, update.Message.Text)
// 		msg.ReplyToMessageID = update.Message.MessageID

// 		if resp, _, err := api.SendMessage(msg); err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func ExampleNewWebhook() {
// 	bot, err := New("MyAwesomeBotToken")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Authorized on account %s", api.Self.UserName)

// 	_, err = api.SetWebhook(NewWebhookWithCert("https://www.google.com:8443/"+api.Token, "cert.pem"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	info, err := api.GetWebhookInfo()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if info.LastErrorDate != 0 {
// 		log.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
// 	}
// 	updates := api.ListenForWebhook("/" + api.Token)
// 	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

// 	for update := range updates {
// 		log.Printf("%+v\n", update)
// 	}
// }

// func ExampleBotAPI_SetWebhook() {
// 	bot, err := New("MyAwesomeBotToken")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Authorized on account %s", api.Self.UserName)

// 	_, err = api.SetWebhook(NewWebhookWithCert("https://www.google.com:8443/"+api.Token, "cert.pem"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	info, err := api.GetWebhookInfo()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if info.LastErrorDate != 0 {
// 		log.Printf("[Telegram callback failed]%s", info.LastErrorMessage)
// 	}

// 	http.HandleFunc("/"+api.Token, func(w http.ResponseWriter, r *http.Request) {
// 		update, err := api.HandleUpdate(r)
// 		if err != nil {
// 			log.Printf("%+v\n", err.Error())
// 		} else {
// 			log.Printf("%+v\n", *update)
// 		}
// 	})

// 	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
// }

// func ExampleBotAPI_AnswerInlineQuery() {
// 	bot, err := New("MyAwesomeBotToken") // create new bot
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	log.Printf("Authorized on account %s", api.Self.UserName)

// 	u := NewUpdate(0)
// 	u.Timeout = 60

// 	updates, err := api.GetUpdatesChan(u)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for update := range updates {
// 		if update.InlineQuery == nil { // if no inline query, ignore it
// 			continue
// 		}

// 		article := NewInlineQueryResultArticle(update.InlineQuery.ID, "Echo", update.InlineQuery.Query)
// 		article.Description = update.InlineQuery.Query

// 		inlineConf := InlineConfig{
// 			InlineQueryID: update.InlineQuery.ID,
// 			IsPersonal:    true,
// 			CacheTime:     0,
// 			Results:       []interface{}{article},
// 		}

// 		if _, err := api.AnswerInlineQuery(inlineConf); err != nil {
// 			log.Println(err)
// 		}
// 	}
// }
*/
