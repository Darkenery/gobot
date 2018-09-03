package model

type ChatMember struct {
	User   *User  `json:"user"`
	Status string `json:"status"`
}
