package platform

func Slash(platform string) string {
	switch platform {
	case "linux":
		return "/"
	case "windows":
		return "\\"
	default:
		return "/"
	}
}
