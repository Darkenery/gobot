package api

import (
	"github.com/darkenery/gobot/api/httputils"
	"github.com/darkenery/gobot/api/model"
	"github.com/darkenery/gobot/api/model/response"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var (
	httpError  = errors.New("Http status is not OK")
	parseError = errors.New("Can't parse Telegram response")
)

type BotApi struct {
	url string
	r   *httputils.HttpRequest
}

func NewBotApi(url, token string, timeout, keepAlive, handshakeTimeout time.Duration) *BotApi {
	return &BotApi{
		url: url + token,
		r: httputils.NewRequest(
			timeout,
			keepAlive,
			handshakeTimeout),
	}
}

func (b *BotApi) GetMe() (user *model.User, err error) {
	url := b.url + "/getMe"
	httpStatus, data, err := b.r.SendJSON(
		url,
		"GET",
		nil,
		model.GetMeResponse{},
	)

	if err != nil {
		return
	}

	if httpStatus != http.StatusOK {
		return nil, httpError
	}

	if response, ok := data.(*model.GetMeResponse); ok {
		return response.Result, nil
	}

	return nil, parseError
}

func (b *BotApi) SendMessage(chatId int64, replyToMessageId int, text string) (message *model.Message, err error) {
	url := b.url + "/sendMessage"

	sendMessageRequest := make(map[string]interface{})
	sendMessageRequest["chat_id"] = chatId
	sendMessageRequest["reply_to_message_id"] = replyToMessageId
	sendMessageRequest["text"] = text

	httpStatus, data, err := b.r.SendJSON(
		url,
		"POST",
		sendMessageRequest,
		response.SendMessageResponse{},
	)

	if err != nil {
		return
	}

	if httpStatus != http.StatusOK {
		return nil, httpError
	}

	if sendMessageResponse, ok := data.(*response.SendMessageResponse); ok {
		return sendMessageResponse.Result, nil
	}

	return nil, parseError
}

func (b *BotApi) GetUpdates(offset, limit, timeout int, allowedUpdates []string) (updates []*model.Update, err error) {
	url := b.url + "/getUpdates"

	getUpdatesRequest := make(map[string]interface{})
	getUpdatesRequest["offset"] = offset
	getUpdatesRequest["limit"] = limit
	getUpdatesRequest["timeout"] = timeout
	getUpdatesRequest["allowed_updates"] = allowedUpdates

	httpStatus, data, err := b.r.SendJSON(
		url,
		"GET",
		getUpdatesRequest,
		response.GetUpdatesResponse{},
	)

	if err != nil {
		return
	}

	if httpStatus != http.StatusOK {
		return nil, httpError
	}

	if getUpdatesResponse, ok := data.(*response.GetUpdatesResponse); ok {
		return getUpdatesResponse.Result, nil
	}

	return nil, parseError
}
