package app_param

type (
	Image struct {
		FileType string `json:"file_type" form:"file_type"`
		Channel  string `json:"channel" form:"channel"`
		ID       int    `json:"id" form:"id"`
	}
	Video struct {
		FileType string `json:"file_type" form:"file_type"`
		Channel  string `json:"channel" form:"channel"`
		ID       int    `json:"id" form:"id"`
	}
	Music struct {
		FileType string `json:"file_type" form:"file_type"`
		Channel  string `json:"channel" form:"channel"`
		ID       int    `json:"id" form:"id"`
	}
)
