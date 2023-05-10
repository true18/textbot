package bot

import (
	"fmt"
	"strings"
	"vkbot/state"
	"vkbot/vk"
)

// ReverseCommand отправляет кнопки выбора зеркалирования текста
func ReverseCommand(e Event) {
	state.FuncState.Add(e.msg.PeerID, Reverse)

	kb := vk.CreateKeyboard(true, false, [][]vk.Button{
		{
			vk.AddButton("Зеркалить символы", ReverseSymbols, vk.KbPrimary),
		},
		{
			vk.AddButton("Зеркалить слова", ReverseWords, vk.KbDefault),
		},
		{
			vk.AddButton("Вернуться в начало", Close, vk.KbNegative),
		},
	})

	_ = e.c.SendMessage(e.msg.PeerID, "Хорошо, выберите режим зеркалирования текста:\n\nРеверсить символы: \"Всем привет!\" → \"!тевирп месВ\""+
		"\nРеверсить слова: \"Всем привет!\" → \"привет! Всем\"", kb.String(), "", 0)
}

// ReverseSelect обрабатывает выбор зеркалирования текста
func ReverseSelect(e Event, value int) {
	state.FuncState.Add(e.msg.PeerID, value)

	kb := vk.CreateKeyboard(false, false, [][]vk.Button{
		{
			vk.AddButton("Завершить", Close, vk.KbPrimary),
		},
	})

	text := "Отлично! Вы выбрали зеркалирование %s. Присылайте мне любой текст, а я буду зеркалить %s.\n\nДля переключения на другой инструмент нажмите на кнопку \"Завершить\" "

	if value == ReverseSymbols {

		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "символов", "символы"), kb.String(), "", 0)
	} else if value == ReverseWords {
		state.FuncState.Add(e.msg.PeerID, ReverseWords)
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "слов", "слова"), kb.String(), "", 0)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}

// ReverseHandler обрабатывает текст, если state пользователя равен ReverseSymbols или ReverseWords
func ReverseHandler(e Event, state int) {
	if state == ReverseSymbols {
		_ = e.c.SendMessage(e.msg.PeerID, ReverseStringSymbols(e.msg.Text), "", "", e.msg.ConversationMessageID)
	} else if state == ReverseWords {
		_ = e.c.SendMessage(e.msg.PeerID, ReverseStringWords(e.msg.Text), "", "", e.msg.ConversationMessageID)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}

// ReverseStringSymbols зеркалит символы
func ReverseStringSymbols(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// ReverseStringWords зеркалит слова
func ReverseStringWords(s string) string {
	words := strings.Fields(s)

	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	return strings.Join(words, " ")
}
