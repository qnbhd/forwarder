package main

import (
	"encoding/json"
	"errors"
	"fmt"
	vkapi "github.com/himidori/golang-vk-api"
	"strings"
)

func getVoiceAttachment(m *vkapi.LongPollMessage) (string, error) {
	res, ok := m.Attachments["attachments"]
	if ok {
		var result map[string]interface{}
		_ = json.Unmarshal([]byte(res[1:len(res)-1]), &result)
		if result["audio_message"] != nil {
			u := result["audio_message"].(map[string]interface{})
			return u["link_ogg"].(string), nil
		}
	}
	return "", errors.New("no voice att")
}

func CollectMessage(client *vkapi.VKClient, m *vkapi.LongPollMessage) string {
	response, _ := client.UsersGet([]int{int(m.UserID)})
	sender := response[0].FirstName + " " + response[0].LastName

	m.Body = strings.Replace(m.Body, "<br>", "\n", -1)

	if strings.Trim(sender, " ") == "DELETED" {
		sender = m.Title
	}

	body := fmt.Sprintf("Новое сообщение от %s:\n%s\n", sender, m.Body)
	voiceAtt, err := getVoiceAttachment(m)

	if err == nil {
		body += voiceAtt + "\n"
	}

	body = strings.ReplaceAll(body, "<br>", "\n")
	body = strings.ReplaceAll(body, "&quot;", "")

	return body
}
