package main

func ResponseOk() map[string]interface{} {
	return map[string]interface{}{
		"message": "ok",
		"error":   "",
		"data":    nil,
	}
}

func ResponseErr(err string) map[string]interface{} {
	return map[string]interface{}{
		"message": "error",
		"error":   err,
		"data":    nil,
	}
}

func ResponseData(data any) map[string]interface{} {
	return map[string]interface{}{
		"message": "ok",
		"error":   "",
		"data":    data,
	}
}
