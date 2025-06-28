package response

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Result {
	return Result{
		Code:    200,
		Message: "操作成功",
		Data:    data,
	}
}

func Error(code int, message string) Result {
	return Result{
		Code:    code,
		Message: message,
	}
}

func BadRequest(message string) Result {
	return Error(400, message)
}

func Unauthorized(message string) Result {
	return Error(401, message)
}

func InternalError(message string) Result {
	return Error(500, message)
}
