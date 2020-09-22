package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	vkapi "github.com/himidori/golang-vk-api"
	"log"
	"strings"
	"time"
)

const (
	tgUpdateTimeout = 60
)

var AccountList []Account
var AccountInUsing = map[string]bool{}

func getTelegramMessage(bot *tgbotapi.BotAPI) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = tgUpdateTimeout

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch {
		case strings.HasPrefix(update.Message.Text, "/go"):
			routeGo(update, bot)
		case strings.HasPrefix(update.Message.Text, "/start"):
			routeStart(update, bot)
		case strings.HasPrefix(update.Message.Text, "/delete"):
			routeDeleteAccount(update, bot)
		default:
			routeDefault(update, bot)
		}

	}
}

func runWatchingGoroutine(currentAccount Account, bot *tgbotapi.BotAPI) {

	if currentAccount.VkLogin == "" || currentAccount.VkPass == "" {
		fmt.Println("Env vk-data isn't filled")
		return
	}

	if res, ok := AccountInUsing[currentAccount.VkLogin]; res == true && ok {
		return
	}

	log.Printf("Authorized on account %s", currentAccount.VkLogin)
	client, err := vkapi.NewVKClient(vkapi.DeviceIPhone, currentAccount.VkLogin, currentAccount.VkPass, false)

	if err != nil {
		fmt.Println("Error vk connecting")
		return
	}

	var cancelSignal = new(bool)
	*cancelSignal = false

	AccountInUsing[currentAccount.VkLogin] = true

	client.AddLongpollCallback("msgin", func(m *vkapi.LongPollMessage) {
		if _, ok := AccountInUsing[currentAccount.VkLogin]; !ok {
			*cancelSignal = true
		}

		body := CollectMessage(client, m)
		msg := tgbotapi.NewMessage(currentAccount.TgId, body)
		_, err = bot.Send(msg)
	})

	client.ListenLongPollServer(cancelSignal)
}

func main() {
	bot, _ := tgbotapi.NewBotAPI(TelegramToken)
	AccountList, _ = loadAccounts()

	for _, account := range AccountList {
		go runWatchingGoroutine(account, bot)
	}

	go getTelegramMessage(bot)

	for {
		<-time.After(time.Second * 1)
	}
}
