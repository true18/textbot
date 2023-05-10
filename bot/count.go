package bot

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"vkbot/state"
	"vkbot/vk"
)

// CountCommand отправляет кнопки выбора регистра текста
func CountCommand(e Event) {
	state.FuncState.Add(e.msg.PeerID, Count)

	kb := vk.CreateKeyboard(true, false, [][]vk.Button{
		{
			vk.AddButton("Количество символов", CountSymbols, vk.KbPrimary),
		},
		{
			vk.AddButton("Количество слов", CountWords, vk.KbDefault),
		},
		{
			vk.AddButton("Вернуться в начало", Close, vk.KbNegative),
		},
	})

	_ = e.c.SendMessage(e.msg.PeerID, "Хорошо, выберите режим подсчёта слов:\n\nПодсчёт символов: \"Всем привет!\" → \"В тексте 12 символов\"\nПодсчёт слов: \"Всем привет!\" → \"В тексте 2 слова\"", kb.String(), "", 0)
}

// CountSelect обрабатывает подсчёта количества символов или слов
func CountSelect(e Event, value int) {
	state.FuncState.Add(e.msg.PeerID, value)

	kb := vk.CreateKeyboard(false, false, [][]vk.Button{
		{
			vk.AddButton("Завершить", Close, vk.KbPrimary),
		},
	})

	text := "Отлично! Вы выбрали подсчёт %s в тексте. Присылайте мне любой текст, а я считать количество %s в тексте.\n\nДля переключения на другой инструмент нажмите на кнопку \"Завершить\" "

	if value == CountSymbols {
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "символов", "символов"), kb.String(), "", 0)
	} else if value == CountWords {
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "слов", "слов"), kb.String(), "", 0)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}

// CountHandler обрабатывает текст, если state пользователя равен CountSymbols или CountWords
func CountHandler(e Event, state int) {
	if state == CountSymbols {
		_ = e.c.SendMessage(e.msg.PeerID,
			fmt.Sprintf("В тексте %d %s", utf8.RuneCountInString(e.msg.Text), PrettyWord(utf8.RuneCountInString(e.msg.Text), "символ", "символа", "символов")), "", "", e.msg.ConversationMessageID)
	} else if state == CountWords {
		_ = e.c.SendMessage(e.msg.PeerID,
			fmt.Sprintf("В тексте %d %s", len(strings.Fields(e.msg.Text)), PrettyWord(len(strings.Fields(e.msg.Text)), "слово", "слова", "слов")), "", "", e.msg.ConversationMessageID)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}
