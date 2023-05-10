package vk

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type getGroupIDData struct {
	ID int `json:"id"`
}

type getGroupIDResponse struct {
	Response []getGroupIDData `json:"response"`
}

// GetGroupID возвращает ClientID сообщества, от имени
// которого делается запрос
func (c *Client) GetGroupID() int {
	jsonR := c.Request("groups.getById", nil)
	response := getGroupIDResponse{}
	err := json.Unmarshal(jsonR, &response)
	CheckError(err)

	if len(response.Response) == 0 {
		panic(fmt.Sprintf("GetGroupID() return 0 \n%s\n%#v\n", string(jsonR), c))
	}

	return response.Response[0].ID
}

// SendSticker отправляет стикер, используя его ClientID
func (c *Client) SendSticker(peerID int, stickerID int) []byte {
	var params = make(map[string]string)
	params["peer_id"] = strconv.Itoa(peerID)
	params["sticker_id"] = strconv.Itoa(stickerID)
	params["random_id"] = "0"

	return c.Request("messages.send", params)
}

type ContentSource struct {
	Type                  string `json:"type"`
	OwnerID               int    `json:"owner_id"`
	PeerID                int    `json:"peer_id"`
	ConversationMessageID int    `json:"conversation_message_id"`
}

// SendMessage отправляет сообщение
func (c *Client) SendMessage(peerID int, message, keyboard, attachment string, source int) []byte {
	var params = make(map[string]string)

	params["message"] = message
	params["keyboard"] = keyboard
	params["attachment"] = attachment
	params["peer_id"] = strconv.Itoa(peerID)

	if source != 0 {
		contentSource, err := json.Marshal(ContentSource{Type: "message", OwnerID: -c.ClientID, PeerID: peerID, ConversationMessageID: source})
		CheckError(err)

		params["content_source"] = string(contentSource)
	}

	params["random_id"] = "0"

	return c.Request("messages.send", params)
}
