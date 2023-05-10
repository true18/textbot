package vk

import (
	"log"
	"strconv"
)
import "encoding/json"

const (
	KbNegative = "negative"
	KbPositive = "positive"
	KbDefault  = "default"
	KbPrimary  = "primary"
)

type Payload struct {
	Command string `json:"button"`
}

type Action struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label"`
}

type Button struct {
	Action `json:"action"`
	Color  string `json:"color"`
}

// SetPayload устанавливает значение полезной нагрузки в command
func (bi *Action) SetPayload(pl interface{}) {
	p, ok := pl.(string)
	if !ok {
		p = strconv.Itoa(pl.(int))
	}

	j, err := json.Marshal(Payload{p})
	if err != nil {
		panic(err)
	}

	bi.Payload = string(j)
}

// GetPayload возвращает значение полезной нагрузки в command
func (msg *Message) GetPayload() string {
	if msg.Payload == "" {
		return ""
	}

	data := Payload{}
	err := json.Unmarshal([]byte(msg.Payload), &data)
	if err != nil {
		return ""
	}
	return data.Command
}

// AddButton возвращает объект Button с заданными данными
func AddButton(name string, payload interface{}, color string) Button {
	action := Action{
		Type:  "text",
		Label: name,
	}
	action.SetPayload(payload)

	return Button{
		Action: action,
		Color:  color,
	}
}

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Inline  bool       `json:"inline"`
	Buttons [][]Button `json:"buttons"`
	Cache   string     `json:"-"`
}

// CreateKeyboard создаёт и возвращает объект клавиатуры
func CreateKeyboard(inline bool, oneTime bool, buttonsGrid [][]Button) *Keyboard {
	return &Keyboard{Inline: inline, OneTime: oneTime, Buttons: buttonsGrid}
}

// Возвращает клавиатуру в формате Json
func (kb *Keyboard) String() string {
	if kb.Cache == "" {
		data, err := json.Marshal(kb)
		if err != nil {
			log.Fatalln(data)
		}
		kb.Cache = string(data)
	}
	return kb.Cache
}
