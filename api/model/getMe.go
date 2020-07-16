package model

type GetMeResponse struct {
	Ok     bool `json:"ok"`
	Result *User `json:"result"`
}
