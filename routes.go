package main

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	vkapi "github.com/himidori/golang-vk-api"
	"strings"
)

func routeGo(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	words := strings.Split(update.Message.Text, " ")

	if len(words) != 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Уупс! Кажется, я не смог прочитать эти данные 😔")
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	accountInstance := Account{
		VkLogin: words[1], VkPass: words[2],
		TgId: update.Message.Chat.ID,
	}

	_, err := vkapi.NewVKClient(vkapi.DeviceIPhone, accountInstance.VkLogin, accountInstance.VkPass, false)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Уупс! Кажется, эти данные не корректны (я не смог войти в аккаунт) 😔")
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
		return
	}

	for _, account := range AccountList {
		if account.VkLogin == accountInstance.VkLogin &&
			account.TgId == accountInstance.TgId {
			return
		}
	}

	AccountList = append(AccountList, accountInstance)
	AccountInUsing[accountInstance.VkLogin] = true
	_ = saveAccounts(AccountList)

	go runWatchingGoroutine(accountInstance, bot)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Супер! Теперь можешь ждать уведомлений 🤩")
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

func routeDeleteAccount(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	words := strings.Split(update.Message.Text, " ")

	switch len(words) {
	case 1:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			`Чтобы удалить аккаунт из отслеживаемых просто введите /delete и логин аккаунта от ВКонтакте!`)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	case 2:
		login := words[1]

		exists := false
		targetIdx := -1
		for idx, account := range AccountList {
			if account.VkLogin == login && account.TgId == update.Message.Chat.ID {
				targetIdx = idx
				exists = true
				break
			}
		}

		if !exists {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хей! Этот аккаунт не отслеживается!")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
			return
		}

		AccountList = append(AccountList[:targetIdx], AccountList[targetIdx+1:]...)
		_ = saveAccounts(AccountList)
		delete(AccountInUsing, login)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хорошо, этот аккаунт я удалил!")
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Уупс! Я не смог прочитать эти данные 😔")
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msg)
	}

}

func routeStart(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf(
			`Привет, %s 😊
Для начала нашей работы с тобой введи команду /go и напиши свои данные от аккаунта ВКонтакте.
Пример: /go vasyaLogin vasyaPassword
				`, update.Message.Chat.FirstName))
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}

func routeDefault(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не понял тебя 🙁")
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
