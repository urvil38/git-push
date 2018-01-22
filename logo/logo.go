package logo

// Wrap returns colored ascii text .Now it only use cyan color but you can change it
func Wrap(s string) string {
	return format() + s + unformat()
}

func format() string {
	return "\x1b[1;36m"
}

func unformat() string {
	return "\x1b[0m"
}