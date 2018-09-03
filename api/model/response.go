package model

type Response struct {
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}
