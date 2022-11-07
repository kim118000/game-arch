module github.com/kim118000/gate

go 1.19

require (
	github.com/BurntSushi/toml v1.2.1
	github.com/kim118000/core v0.0.0
	github.com/kim118000/protocol v0.0.0
	google.golang.org/protobuf v1.28.1
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/kim118000/core => ../core
	github.com/kim118000/protocol => ../protocol
)
