package vk

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)
import "encoding/json"

type LongPoll struct {
	*Client
	Key    string
	Server string
	TS     string
}

type Attachment struct {
	Type string `json:"type"`
}

type getLongPollServerData struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type Message struct {
	FromID                int          `json:"from_id"`
	Out                   int          `json:"out"`
	PeerID                int          `json:"peer_id"`
	ConversationMessageID int          `json:"conversation_message_id"`
	Text                  string       `json:"text"`
	Payload               string       `json:"Payload"`
	Attachments           []Attachment `json:"attachments"`
}

type longPollEvent struct {
	Message Message `json:"message"`
}

type longPollUpdate struct {
	longPollEvent `json:"object"`
	Type          string `json:"type"`
	GroupId       int    `json:"group_id"`
}

type longPollResponse struct {
	Ts      string           `json:"ts"`
	Updates []longPollUpdate `json:"updates"`
	Failed  int              `json:"failed"`
}

type getLongpollServerResponse struct {
	Response getLongPollServerData `json:"response"`
}

// NewLongPoll возвращает объект для работы с longpoll ВКонтакте
func NewLongPoll(token string) *LongPoll {
	return &LongPoll{Client: NewClient(token)}
}

// Инициализирует параметры работы с longpoll сервером
func (lp *LongPoll) initVKParams() {
	jsonR := lp.Request("groups.getLongPollServer", map[string]string{"group_id": strconv.Itoa(lp.ClientID)})
	response := getLongpollServerResponse{}
	err := json.Unmarshal(jsonR, &response)
	CheckError(err)

	lp.Key = response.Response.Key
	lp.Server = response.Response.Server
	lp.TS = response.Response.Ts
}

// Запрашивает новые обновления у longpoll-сервера
func (lp *LongPoll) getEvents() (longPollResponse, error) {
	url := fmt.Sprintf("%s?act=a_check&key=%s&ts=%s&wait=25", lp.Server, lp.Key, lp.TS)
	r, err := http.Get(url)

	if err != nil {
		return longPollResponse{}, err
	}

	defer r.Body.Close()

	answer, err := ioutil.ReadAll(r.Body)
	log.Println(string(answer))
	CheckError(err)

	response := longPollResponse{}
	err = json.Unmarshal(answer, &response)
	CheckError(err)

	return response, nil
}

// ListenMessages прослушивание longpoll-сервера
func (lp *LongPoll) ListenMessages(inputMessages chan<- *Message) {
	lp.initVKParams()

	for {
		response, err := lp.getEvents()
		if response.Failed != 0 || response.Ts == "" || err != nil {
			log.Println("Restart longpoll")
			lp.initVKParams()
			continue
		}

		lp.TS = response.Ts
		for _, event := range response.Updates {
			if event.Type != "message_new" || event.Message.Out == 1 || event.Message.FromID < 0 {
				continue
			}

			inputMessages <- &event.Message
		}
	}
}
