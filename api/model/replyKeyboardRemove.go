package model

type ReplyKeyboardRemove struct {
	HideKeyboard bool `json:"hide_keyboard"`
	Selective    bool `json:"selective"`
}
