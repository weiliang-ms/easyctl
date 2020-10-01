package sys

func packageRedhatService(sysVersion string, name string, action string, forever bool) string {

	if sysVersion == "6" {
		return redhat6ManageService(name, action, forever)
	} else if sysVersion == "7" {
		return redhat7ManageService(name, action, forever)
	} else {
		return ""
	}
}

func redhat6ManageService(name string, action string, forever bool) string {
	if name == "firewalld" {
		name = "iptables"
	}
	if action == stop {
		return redhat6StopService(name, forever)
	} else {
		return redhat6StartService(name, forever)
	}
}

func redhat6StopService(name string, forever bool) string {

	if forever {
		return chkconfig + " " + name + " " + off
	} else {
		return service + " " + name + " " + stop
	}
}
func redhat6StartService(name string, forever bool) string {
	if forever {
		return chkconfig + " " + name + " " + on
	} else {
		return service + " " + name + " " + start
	}
}

func redhat7ManageService(name string, action string, forever bool) string {
	if action == stop {
		return redhat7StopService(name, forever)
	} else {
		return redhat7StartService(name, forever)
	}
}

func redhat7StopService(name string, forever bool) string {
	if forever {
		return systemctl + " " + disable + " " + name
	} else {
		return systemctl + " " + stop + " " + name
	}
}
func redhat7StartService(name string, forever bool) string {
	if forever {
		return systemctl + " " + enable + " " + name
	} else {
		return systemctl + " " + start + " " + name + " --now"
	}
}
