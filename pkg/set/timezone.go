package set

import "log"

func (ac *Actuator) TimeZone() {
	log.Println("配置时区...")
	ac.parseServerList().setTimeZoneCmd().execute("配置时区为上海时区")
}

func (ac *Actuator) setTimeZoneCmd() *Actuator {
	ac.Cmd = "\\cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -R"
	return ac
}
