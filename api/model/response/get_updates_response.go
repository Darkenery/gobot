package response

import "github.com/darkenery/gobot/api/model"

type GetUpdatesResponse struct {
	Ok     bool            `json:"ok"`
	Result []*model.Update `json:"result"`
}
