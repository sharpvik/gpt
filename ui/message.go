package ui

func message(prefix, text string) string {
	return " " + prefix + "  " + text + "\n\n"
}

func aiMessage(text string) string {
	return message("🤖", text)
}

func humanMessage(text string) string {
	return message("👾", text)
}

func errorMessage(err error) string {
	return message("🚨", err.Error())
}
