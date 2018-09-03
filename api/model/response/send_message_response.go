package response

import "github.com/darkenery/gobot/api/model"

type SendMessageResponse struct {
	Ok     bool           `json:"ok"`
	Result *model.Message `json:"result"`
}
