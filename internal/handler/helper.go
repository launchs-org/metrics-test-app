package handler

func errorResp(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}
