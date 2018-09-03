package model

type Sticker struct {
	FileId   string     `json:"file_id"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Thumb    *PhotoSize `json:"thumb"`
	Emoji    string     `json:"emoji"`
	FileSize int        `json:"file_size"`
}
