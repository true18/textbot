package vk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiVersion = "5.131"

const requestURL string = "https://api.vk.com/method/%s?access_token=%s&v=" + apiVersion

const sleepTime = time.Second / 20

func getRequestUrl(method, token string) string {
	return fmt.Sprintf(requestURL, method, token)
}

type Client struct {
	token       string
	ClientID    int
	lastRequest time.Time
}

// NewClient создаёт новую сессию работы с API
func NewClient(token string) *Client {
	c := &Client{token: token}
	c.ClientID = c.GetGroupID()

	return c
}

// Request принимает метод и параметры к API ВКонтакте
func (c *Client) Request(method string, params map[string]string) []byte {
	// Оптимизация отправки запросов
	left := time.Since(c.lastRequest)
	if left < sleepTime {
		time.Sleep(sleepTime - left)
	}
	c.lastRequest = time.Now()

	rURL := getRequestUrl(method, c.token)

	data := url.Values{}
	for key, value := range params {
		data.Set(key, value)
	}
	reader := strings.NewReader(data.Encode())

	r, err := http.Post(rURL, "application/x-www-form-urlencoded", reader)
	CheckError(err)
	defer r.Body.Close()

	binAnswer, err := ioutil.ReadAll(r.Body)
	CheckError(err)

	if strings.Contains(string(binAnswer), "err") {
		fmt.Println(string(binAnswer))
		fmt.Println(rURL)
		fmt.Printf("%#v\n\n", params)
	}
	return binAnswer
}
