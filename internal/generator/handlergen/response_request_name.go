package handlergen

func GetHandlerRequestName(handlerName string) string {
	return handlerName + "Request"
}

func GetHandlerResponseName(handlerName string) string {
	return handlerName + "Response"
}
