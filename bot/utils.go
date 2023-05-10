package bot

// PrettyWord склоняет слово в зависимости от значения числа
// Пример: 1 сообщение, 2 сообщения, 5 сообщений
func PrettyWord(number int, s1, s2, s3 string) string {
	var titles = []string{s1, s2, s3}

	cases := []int{2, 0, 1, 1, 1, 2}

	switch {
	case number%100 > 4 && number%100 < 20:
		return titles[2]
	case number%10 < 5:
		return titles[cases[number%10]]
	default:
		return titles[cases[5]]
	}
}
