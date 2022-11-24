package config

type ServerBase struct {
	ServerName string
	ServerPort uint32
	ServerId   int32
}

type TemplateBase struct {
	Time  uint64 `json:"time"`
	VTime string `json:"vtime"`
}
