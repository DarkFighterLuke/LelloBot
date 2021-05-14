package main

import (
	"encoding/json"
	"fmt"
	"github.com/NicoNex/echotron"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type bot struct {
	chatId int64
	echotron.Api
	roundNegazione int
}

const (
	botLogsFolder  = "/LelloBotData/logs/"
	botAudioFolder = "/LelloBotData/audio/"
)

var TOKEN = os.Getenv("LelloBot")
var logsFolder string
var audioFolder string

func newBot(chatId int64) echotron.Bot {
	return &bot{
		chatId,
		echotron.NewApi(TOKEN),
		0,
	}
}

func (b *bot) makeButtons(buttonsText []string, callbacksData []string, layout int) ([]byte, error) {
	if layout != 1 && layout != 2 {
		return nil, fmt.Errorf("wrong layout")
	}
	if len(buttonsText) != len(callbacksData) {
		return nil, fmt.Errorf("different text and data length")
	}

	buttons := make([]echotron.InlineButton, 0)
	for i, v := range buttonsText {
		buttons = append(buttons, b.InlineKbdBtn(v, "", callbacksData[i]))
	}

	keys := make([]echotron.InlineKbdRow, 0)
	switch layout {
	case 1:
		for i := 0; i < len(buttons); i++ {
			keys = append(keys, echotron.InlineKbdRow{buttons[i]})
		}
		break
	case 2:
		for i := 0; i < len(buttons); i += 2 {
			if i+1 < len(buttons) {
				keys = append(keys, echotron.InlineKbdRow{buttons[i], buttons[i+1]})
			} else {
				keys = append(keys, echotron.InlineKbdRow{buttons[i]})
			}
		}
		break
	}

	inlineKMarkup := b.InlineKbdMarkup(keys...)
	return inlineKMarkup, nil
}

func initFolders() {
	currentPath, _ := os.Getwd()

	logsFolder = currentPath + botLogsFolder
	_ = os.MkdirAll(logsFolder, 0755)

	audioFolder = currentPath + botAudioFolder
	_ = os.MkdirAll(audioFolder, 0755)
}

func main() {
	initFolders()

	dsp := echotron.NewDispatcher(TOKEN, newBot)
	dsp.ListenWebhook("https://hiddenfile.tk:443/bot/LelloBot", 40989)
}

func (b *bot) Update(update *echotron.Update) {
	b.logUser(update, logsFolder)
	if update.Message != nil {
		messageTextLower := strings.ToLower(update.Message.Text)
		if messageTextLower == "/start" {
			b.sendStart(update.Message)
		} else if messageTextLower == "/credits" {
			b.sendCredits(update)
		} else if strings.Contains(messageTextLower, "cant") && strings.Contains(messageTextLower, "canzone") {
			b.sendLelloSong(update.Message)
		} else if strings.Contains(messageTextLower, "stai pieno") {
			b.roundNegazione = 1
			b.sendLelloNegazioneSbronza(update.Message)
		} else if strings.Contains(messageTextLower, "dove") && strings.Contains(messageTextLower, "vai") &&
			strings.Contains(messageTextLower, "le") || strings.Contains(messageTextLower, "lè") ||
			strings.Contains(messageTextLower, "lé") {
			b.sendLelloTypicalExpression(update.Message, 22)
		} else if strings.Contains(messageTextLower, "angela") {
			b.sendLelloTypicalExpression(update.Message, 24)
		} else if strings.Contains(messageTextLower, "ubriac") || strings.Contains(messageTextLower, "mbriac") {
			b.sendLelloTypicalExpression(update.Message, 12)
		} else if b.roundNegazione != 0 && (strings.Contains(messageTextLower, "sì") ||
			strings.Contains(messageTextLower, "si")) {
			b.sendLelloNegazioneSbronza(update.Message)
		} else if b.roundNegazione != 0 && strings.Contains(messageTextLower, "no") {
			b.sendLelloTypicalExpression(update.Message, 7)
		} else if strings.Contains(messageTextLower, "lello") || strings.Contains(messageTextLower, "lè") ||
			strings.Contains(messageTextLower, "lé") {
			b.sendLelloTypicalExpression(update.Message, -1)

		} else if update.Message.Chat.Type == "private" {
			b.privateTalkWithLello(update.Message)
		}
	} else if update.Message == nil && update.CallbackQuery != nil {
		if update.CallbackQuery.Data == "credits" {
			b.sendCredits(update)
		}
	}

}

func (b *bot) sendCredits(update *echotron.Update) {
	var chatId int64
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		chatId = update.Message.Chat.ID
	}

	b.SendMessage("🤖 Bot creato da @GiovanniRanaTortello\n😺 GitHub: https://github.com/DarkFighterLuke\n"+
		"\n🌐 Proudly hosted on Raspberry Pi 3\n"+
		"\nContribuisci anche tu alla linguistica di LelloBot su GitHub o contattando il creatore!\n"+
		"N.B. Questo bot è satirico e non intende offendere chi rappresenta. "+
		"Ti auguriamo una pronta guarigione Lello.", chatId, echotron.PARSE_HTML)
	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(update.CallbackQuery.ID, "Crediti", false)
	}
}

func (b *bot) sendStart(message *echotron.Message) {
	msg := `<b>Hai svegliato Lello!</b>
Piacere di conoscerti, %s!
Io sono Lello.
Beh mo' parliamo un poco...
Se non ti rispondo è che ho preso sonno.
`
	msg = fmt.Sprintf(msg, message.User.FirstName)

	buttonText := []string{"Credits 🌟"}
	buttonCallback := []string{"credits"}
	buttons, err := b.makeButtons(buttonText, buttonCallback, 1)
	if err != nil {
		log.Println("Error creating buttons:", err)
	}

	b.SendMessageWithKeyboard(msg, message.Chat.ID, buttons, echotron.PARSE_HTML)
}

func (b *bot) sendLelloTypicalExpression(message *echotron.Message, n int) {
	if n < 0 {
		n = rand.Intn(25)
	}

	switch n {
	case 0:
		msg := "La Ciocia Ciola"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 1:
		msg := "No!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 2:
		msg := "Non sto!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 3:
		msg := "Questo me ne sbatto un cazzo!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 4:
		msg := "Dillo a loro non a me"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 5:
		msg := "O-O-Ommeladai"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 6:
		msg := "Ommelaprendo"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 7:
		msg := "Ma vedi che ahhhhh!!!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 8:
		msg := "*TSK*"
		b.SendMessage("...", message.Chat.ID)
		b.SendMessage("...", message.Chat.ID)
		b.SendMessage(msg, message.Chat.ID)
		break
	case 9:
		msg := "*manifestazione di dissenso*"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 10:
		stickerId := "CAACAgQAAxkBAAMdYDA2FjribEzLMIiSPRdH8cgSyZUAAg8AAw9Q-xDrE0-TbJNyXh4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 11:
		stickerId := "CAACAgQAAxkBAAMhYDA2b34gHAWywAl8zBk7FRwoCHYAAg4AAw9Q-xDrcu-a7OMXmB4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 12:
		if b.roundNegazione == 0 {
			b.sendLelloNegazioneSbronza(message)
		}
		break
	case 13:
		msg := "Seh, o mo'!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 14:
		msg := "A me mi serve casa"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 15:
		// Sends "A_me_mi_serve_casa.mp3"
		fileId := "AwACAgQAAxkBAANFYDLbXlZJL6zn8aSTiFthxi6w-IQAAswJAAJKqZlRmdxNX9F1ZrseBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 16:
		// Sends "Non_sto_umbriacato.mp3"
		if b.roundNegazione == 0 {
			fileId := "AwACAgQAAxkBAANJYDLbpwMJPry0s8U6Dr8xODde34oAAs4JAAJKqZlRsKmcOp2e-h8eBA"
			b.SendVoiceByID(fileId, "", message.Chat.ID)
			b.roundNegazione = 1
		}
		break
	case 17:
		// Sends "Questo_me_ne_sbatto_un_cazzo.mp3"
		fileId := "AwACAgQAAxkBAANPYDLc10yHO-wCU7ZTUIvzieEI0PsAAtoJAAJKqZlRmLjCyKj7vfkeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 18:
		// Sends "Seh_o_mo.mp3"
		fileId := "AwACAgQAAxkBAANDYDLbEtnWKgABGNPvYeVG0NcwduP6AALKCQACSqmZUbRhIw7Bw2xgHgQ"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 19:
		// Sends "Tsk.mp3"
		fileId := "AwACAgQAAxkBAANLYDLb5tw-7xJ8CgI8dGq9VQQPlxIAAtIJAAJKqZlRcjAR10XbKcMeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 20:
		// Sends "Ma_vedi_che_ahhhhh.mp3"
		fileId := "AwACAgQAAxkBAANHYDLbhQOHO4T8dLbU2SNAPTtrCJ8AAs0JAAJKqZlRtMYwKuvppgUeBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 21:
		// Sends "Ommeladai_Ommelaprendo.mp3"
		fileId := "AwACAgQAAxkBAANNYDLcOHXHKaeSyNVEkuxI2QYbmTEAAtQJAAJKqZlRCmVf0kSJem4eBA"
		b.SendVoiceByID(fileId, "", message.Chat.ID)
		break
	case 22:
		msg := "GIRONZOLAANDOOO!"
		b.SendMessage(msg, message.Chat.ID)
		break
	case 23:
		stickerId := "CAACAgQAAxkBAAMrYDEccZiw8v9nGXddfaFyBETyjJUAAgoAAw9Q-xALVTzbP3nVux4E"
		b.SendStickerByID(stickerId, message.Chat.ID)
		break
	case 24:
		msg := "La voglio sbaciucchiare tutta con quelle sue guanciotte che ha"
		b.SendMessage(msg, message.Chat.ID)
		break
	}

	if b.roundNegazione != 0 && (n != 12 && n != 16) {
		b.roundNegazione = 0
	}
}

func (b *bot) sendLelloNegazioneSbronza(message *echotron.Message) {
	if b.roundNegazione == 0 {
		b.roundNegazione++
		msg := "Non sto umbriacato"
		b.SendMessageReply(msg, message.Chat.ID, message.ID)
	} else if b.roundNegazione < 3 {
		b.roundNegazione++
		msg := "No!"
		b.SendMessageReply(msg, message.Chat.ID, message.ID)
	} else if b.roundNegazione == 3 {
		b.roundNegazione++
		msg := "No! Non sto!"
		b.SendMessageReply(msg, message.Chat.ID, message.ID)
	} else if b.roundNegazione == 4 {
		if n := rand.Intn(13); n <= 3 {
			b.sendLelloTypicalExpression(message, 3)
		} else if n <= 4 {
			b.sendLelloTypicalExpression(message, 4)
		} else if n <= 7 {
			b.sendLelloTypicalExpression(message, 7)
		} else if n <= 8 {
			b.sendLelloTypicalExpression(message, 8)
		} else if n <= 9 {
			b.sendLelloTypicalExpression(message, 9)
		} else if n <= 13 {
			b.sendLelloTypicalExpression(message, 13)
		}
	}
}

func (b *bot) privateTalkWithLello(message *echotron.Message) {
	n := rand.Float32()
	if n < 0.75 {
		b.sendLelloTypicalExpression(message, -1)
	}
}

func (b *bot) sendLelloSong(message *echotron.Message) {
	n := rand.Intn(2)

	switch n {
	case 0:
		fileId := "AwACAgQAAxkBAAMPYDGWbkKgs6VXuzrYXhR6n5jO2j8AAngKAAL9zJBRSwOokr3_fVMeBA"
		b.SendVoiceByID(fileId, "By Dix (@michele.di.croce)", message.Chat.ID)
		break
	case 1:
		fileId := "AwACAgQAAxkBAANRYDLfOf5hsdP2q7ZOIsGX1jFaqzMAAtsJAAJKqZlRUgSukIAU98ceBA"
		b.SendVoiceByID(fileId, "By Vincenzo Santoliquido", message.Chat.ID)
		break
	}
}

func (b *bot) logUser(update *echotron.Update, folder string) {
	data, err := json.Marshal(update)
	if err != nil {
		log.Println("Error marshaling logs: ", err)
		return
	}

	var filename string

	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message.Chat.Type == "private" {
			if update.CallbackQuery.Message.Chat.Username == "" {
				filename = folder + update.CallbackQuery.Message.Chat.FirstName + "_" + update.CallbackQuery.Message.Chat.LastName + ".txt"
			} else {
				filename = folder + update.CallbackQuery.Message.Chat.Username + ".txt"
			}
		} else {
			filename = folder + update.CallbackQuery.Message.Chat.Title + ".txt"
		}

	} else if update.Message != nil {
		if update.Message.Chat.Type == "private" {
			if update.Message.Chat.Username == "" {
				filename = folder + update.Message.Chat.FirstName + "_" + update.Message.Chat.LastName + ".txt"
			} else {
				filename = folder + update.Message.Chat.Username + ".txt"
			}
		} else {
			filename = folder + update.Message.Chat.Title + ".txt"
		}

	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	dataString := time.Now().Format("2006-01-02T15:04:05") + string(data[:])
	_, err = f.WriteString(dataString + "\n")
	if err != nil {
		log.Println(err)
		return
	}
	err = f.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
