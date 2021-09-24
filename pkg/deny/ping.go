package deny

func Ping(config []byte, debug bool) error {
	return Item(config, debug, denyPingShell)
}
