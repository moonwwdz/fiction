package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var urlT = "https://api.telegram.org/bot%s/sendMessage"

type msg struct {
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
	Dn        bool   `json:"disable_notification"`
}

func HttpPost(bot, chatId, cont string) bool {
	urlT = fmt.Sprintf(urlT, bot)
	data := msg{ChatId: chatId,
		Text:      cont,
		ParseMode: "HTML",
		Dn:        false, //提示声音 true:关闭 false:开启
	}
	dataJson, _ := json.Marshal(data)

	resp, err := http.Post(urlT, "application/json;charset=utf-8", bytes.NewBuffer(dataJson))
	if err != nil {
		log.Println("Connect Telegram Failed")
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Telegram data error")
	}
	log.Printf("%v", string(body))
	return true
}
