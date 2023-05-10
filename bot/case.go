package bot

import (
	"fmt"
	"strings"
	"vkbot/state"
	"vkbot/vk"
)

// CaseCommand отправляет кнопки выбора регистра текста
func CaseCommand(e Event) {
	state.FuncState.Add(e.msg.PeerID, Case)

	kb := vk.CreateKeyboard(true, false, [][]vk.Button{
		{
			vk.AddButton("Верхний", CaseUp, vk.KbPrimary),
			vk.AddButton("Нижний", CaseDown, vk.KbDefault),
		},
		{
			vk.AddButton("Вернуться в начало", Close, vk.KbNegative),
		},
	})

	_ = e.c.SendMessage(e.msg.PeerID, "Хорошо, выберите необходимый регистр\n\nВерхний регистр: \"Всем привет!\" → \"ВСЕМ ПРИВЕТ!\"\nНижний регистр: \"ВсЕМ ПРИвЕТ!\" → \"всем привет!\"", kb.String(), "", 0)
}

// CaseSelect обрабатывает выбор регистра
func CaseSelect(e Event, value int) {
	state.FuncState.Add(e.msg.PeerID, value)

	kb := vk.CreateKeyboard(false, false, [][]vk.Button{
		{
			vk.AddButton("Завершить", Close, vk.KbPrimary),
		},
	})

	text := "Отлично! Вы выбрали перевод в %s регистр. Присылайте мне любой текст, а я буду переводить его в %s регистр.\n\nДля переключения на другой инструмент нажмите на кнопку \"Завершить\" "

	if value == CaseUp {
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "верхний", "верхний"), kb.String(), "", 0)
	} else if value == CaseDown {
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "нижний", "нижний"), kb.String(), "", 0)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}

// CaseHandler обрабатывает текст, если state пользователя равен CaseUp или CaseDown
func CaseHandler(e Event, state int) {
	if state == CaseUp {
		_ = e.c.SendMessage(e.msg.PeerID, strings.ToUpper(e.msg.Text), "", "", e.msg.ConversationMessageID)
	} else if state == CaseDown {
		_ = e.c.SendMessage(e.msg.PeerID, strings.ToLower(e.msg.Text), "", "", e.msg.ConversationMessageID)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}
