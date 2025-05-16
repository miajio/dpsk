package dto

func NewBaseResponse(code int, message string, data any, err error) map[string]any {
	base := map[string]any{
		"code":    code,
		"message": message,
	}

	if data != nil {
		base["data"] = data
	}

	if err != nil {
		base["error"] = err.Error()
	}
	return base
}
