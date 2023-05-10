package bot

import (
	"fmt"
	"strings"
	"vkbot/state"
	"vkbot/vk"
)

// CleanVowelsConsonantsCommand отправляет кнопки выбора удаления гласных или согласных
func CleanVowelsConsonantsCommand(e Event) {
	state.FuncState.Add(e.msg.PeerID, CleanVowelsConsonants)

	kb := vk.CreateKeyboard(true, false, [][]vk.Button{
		{
			vk.AddButton("Удалить гласные", CleanVowels, vk.KbPrimary),
		},
		{
			vk.AddButton("Удалить согласные", CleanConsonants, vk.KbDefault),
		},
		{
			vk.AddButton("Вернуться в начало", Close, vk.KbNegative),
		},
	})

	_ = e.c.SendMessage(e.msg.PeerID, "Хорошо, выберите, какие буквы необходимо удалить"+
		"\n\nВерхний регистр: \"Всем привет!\" → \"Всм првт!\"\nНижний регистр: \"Всем привет!\" → \"е ие!\""+
		"\n\nОбратите внимание: поддерживается русский и английский алфавит.", kb.String(), "", 0)
}

// CleanVowelsConsonantsSelect обрабатывает выбор удаления букв
func CleanVowelsConsonantsSelect(e Event, value int) {
	state.FuncState.Add(e.msg.PeerID, value)

	kb := vk.CreateKeyboard(false, false, [][]vk.Button{
		{
			vk.AddButton("Завершить", Close, vk.KbPrimary),
		},
	})

	text := "Отлично! Вы выбрали удаление %s букв. Присылайте мне любой текст, а я буду удалять из него %s буквы." +
		"\n\nДля переключения на другой инструмент нажмите на кнопку \"Завершить\" "

	if value == CleanVowels {
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "гласных", "гласные"), kb.String(), "", 0)
	} else if value == CleanConsonants {
		_ = e.c.SendMessage(e.msg.PeerID, fmt.Sprintf(text, "согласных", "согласные"), kb.String(), "", 0)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}

// CleanVowelsConsonantsHandler обрабатывает текст, если state пользователя
// равен CleanVowels или CleanConsonants
func CleanVowelsConsonantsHandler(e Event, state int) {
	if state == CleanVowels {
		result := RemoveVowels(e.msg.Text)
		if result == "" {
			result = "Получился пустой результат"
		}
		_ = e.c.SendMessage(e.msg.PeerID, result, "", "", e.msg.ConversationMessageID)
	} else if state == CleanConsonants {
		result := RemoveConsonants(e.msg.Text)
		if result == "" {
			result = "Получился пустой результат"
		}

		_ = e.c.SendMessage(e.msg.PeerID, result, "", "", e.msg.ConversationMessageID)
	} else {
		_ = e.c.SendMessage(e.msg.PeerID, "Произошло не совсем то, что я ожидал 🤔", "", "", 0)
	}
}

// RemoveVowels удаляет из строки гласные буквы
func RemoveVowels(input string) string {
	vowels := []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U", "а", "е", "ё", "и", "о",
		"у", "ы", "э", "ю", "я", "А", "Е", "Ё", "И", "О", "У", "Ы", "Э", "Ю", "Я"}

	for _, vowel := range vowels {
		input = strings.ReplaceAll(input, vowel, "")
	}

	return input
}

// RemoveConsonants удаляет из строки согласные буквы
func RemoveConsonants(input string) string {
	consonants := []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p",
		"q", "r", "s", "t", "v", "w", "x", "y", "z",
		"B", "C", "D", "F", "G", "H", "J", "K", "L", "M", "N", "P",
		"Q", "R", "S", "T", "V", "W", "X", "Y", "Z",
		"б", "в", "г", "д", "ж", "з", "й", "к", "л", "м", "н", "п", "р", "с", "т",
		"ф", "х", "ц", "ч", "ш", "щ", "Б", "В", "Г", "Д", "Ж", "З", "Й", "К", "Л", "М", "Н", "П", "Р", "С", "Т",
		"Ф", "Х", "Ц", "Ч", "Ш", "Щ"}

	for _, consonant := range consonants {
		input = strings.ReplaceAll(input, consonant, "")
	}

	return input
}
