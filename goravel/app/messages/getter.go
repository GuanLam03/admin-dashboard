package messages

func GetError(key string) string {
	if msg, ok := ErrorMessages[key]; ok {
		return msg
	}
	return key
}

func GetSuccess(key string) string {
	if msg, ok := SuccessfulMessages[key]; ok {
		return msg
	}
	return key
}