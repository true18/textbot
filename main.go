package main

import (
	"fmt"
	"log"
	"time"
	"vkbot/bot"
)

const maxRetries = 5
const delay = 1 * time.Second

// Токен для чат-бота нужно задаётся переменной BOT_TOKEN в окружении
func main() {
	log.Println("Запуск чат-бота")

	// Хоть и systemctl будет перезапускать скрипт самостоятельно при краше,
	// на каждую панику нет смысла завершать программу, поэтому используем recover
	for i := 0; i < maxRetries; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in main from panic in Start", r)
				}
			}()

			bot.Start()
		}()

		time.Sleep(delay)
	}
}
