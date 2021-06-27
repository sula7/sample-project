package service

func (api *APIv1) httpRespUnsuccessful(message string) *Response {
	return &Response{Success: false, Message: message}
}

func (api *APIv1) httpRespSuccessful(data interface{}) *Response {
	return &Response{Success: true, Message: "OK", Data: data}
}
