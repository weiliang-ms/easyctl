package strings

import "runtime"

func Slash() string {
	switch runtime.GOOS {
	case "linux":
		return "/"
	case "windows":
		return "\\"
	default:
		return "/"
	}
}
