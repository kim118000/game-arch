package log

type LogConfig struct {
	Level int
	ShowGid bool
	OutFormat string
	FileName   string
	MaxSize int
	MaxAge int
	MaxBackups int
	Compress bool
}