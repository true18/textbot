package bot

import (
	"log"
	"os"
	"strconv"
	"vkbot/state"
	"vkbot/vk"
)

type Event struct {
	c             *vk.Client
	inputMessages chan *vk.Message
	msg           *vk.Message
}

type Payload struct {
	Type string
}

// Константы состояния работы с пользователем
const (
	Close = iota + 1
	Case
	CaseUp
	CaseDown
	Reverse
	ReverseSymbols
	ReverseWords
	Count
	CountSymbols
	CountWords
	CleanVowelsConsonants
	CleanVowels
	CleanConsonants
)

// Запуск longpoll и обработка сообщений
func Start() {
	token := os.Getenv("BOT_TOKEN")

	c := vk.NewClient(token)

	chMessages := make(chan *vk.Message, 100)

	lp := vk.NewLongPoll(token)

	log.Println("Запуск longpoll")
	go lp.ListenMessages(chMessages)

	for msg := range chMessages {
		log.Printf("%d: %s (payload: %s)\n", msg.PeerID, msg.Text, msg.Payload)
		Handle(Event{c, chMessages, msg})
	}
}

// Обработка сообщения
func Handle(e Event) {
	// Клавиатура выбора режима работы
	kb := vk.CreateKeyboard(false, true, [][]vk.Button{
		{
			vk.AddButton("Поменять регистр текста", Case, vk.KbPrimary),
		},
		{
			vk.AddButton("Отзеркалить текст", Reverse, vk.KbDefault),
		},
		{
			vk.AddButton("Подсчитать текст", Count, vk.KbPrimary),
		},
		{
			vk.AddButton("Убрать согласные или гласные", CleanVowelsConsonants, vk.KbDefault),
		},
	})

	commandInt, _ := strconv.Atoi(e.msg.GetPayload())

	switch commandInt {
	case Case:
		CaseCommand(e)
	case CaseUp, CaseDown:
		CaseSelect(e, commandInt)
	case Reverse:
		ReverseCommand(e)
	case ReverseSymbols, ReverseWords:
		ReverseSelect(e, commandInt)
	case Count:
		CountCommand(e)
	case CountSymbols, CountWords:
		CountSelect(e, commandInt)
	case CleanVowelsConsonants:
		CleanVowelsConsonantsCommand(e)
	case CleanVowels, CleanConsonants:
		CleanVowelsConsonantsSelect(e, commandInt)
	case Close:
		state.FuncState.Pop(e.msg.PeerID)
		e.c.SendMessage(e.msg.PeerID, "Выберите необходимое действие", kb.String(), "", 0)
	}

	if e.msg.GetPayload() != "" {
		return
	}

	v, ok := state.FuncState.In(e.msg.PeerID)

	switch v {
	case Case, Count, Reverse, CleanVowelsConsonants:
		e.c.SendMessage(e.msg.PeerID, "Сначала необходимо выбрать режим обработки текста", "", "", 0)
	case CaseUp, CaseDown:
		CaseHandler(e, v)
	case ReverseSymbols, ReverseWords:
		ReverseHandler(e, v)
	case CountSymbols, CountWords:
		CountHandler(e, v)
	case CleanVowels, CleanConsonants:
		CleanVowelsConsonantsHandler(e, v)
	}

	if ok {
		return
	}

	// Приветственное сообщение
	e.c.SendSticker(e.msg.PeerID, 86032)

	e.c.SendMessage(e.msg.PeerID, "Я чат-бот для работы с текстом.\n\nДанный чат-бот написан в рамках задания стажировки VK на вакансию \"Go-разработчик\".\nЗадание выполнил: [zinnatullin|Чингиз Зиннатуллин]\n\n Выберите необходимое действие", kb.String(), "", 0)
}
