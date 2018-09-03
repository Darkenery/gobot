package model

type Update struct {
	UpdateId      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	EditedMessage *Message       `json:"edited_message"`
	InlineQuery   *InlineQuery   `json:"inline_query"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}
