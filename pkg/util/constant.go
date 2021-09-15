package util

const (
	Localhost        = "localhost"
	Success          = "success"
	Fail             = "fail"
	AlreadyCloseCode = 3
	ConfigFile       = "config.yaml"
)

type LogLevel int8

const (
	Debug LogLevel = iota
	Info
	Error
)
