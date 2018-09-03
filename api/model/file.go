package model

type File struct {
	Result struct {
		FileId   string `json:"file_id"`
		FileSize int    `json:"file_size"`
		FilePath string `json:"file_path"`
	} `json:"result"`
}
