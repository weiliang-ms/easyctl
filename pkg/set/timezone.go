package set

const setTimezoneShell = "ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime"

func Timezone(config []byte, debug bool) error {
	return Config(config, debug, setTimezoneShell)
}
