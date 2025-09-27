package response


// ApiResponse — универсальный тип ответа
type ApiResponse[T any] struct {
	// Статус ответа
	Success bool `json:"success"`

	// Сообщение (комментарий) об ошибке
	Message string `json:"message,omitempty"`

	// Данные ответа
	Data *T `json:"data,omitempty"`

	// Детали ошибки
	Details string `json:"details,omitempty"`
}

func NewBadResponse[T any](message, details string) ApiResponse[T] {
	return ApiResponse[T]{
		Success: false,
		Message: message,
		Details: details,
	}
}

func NewSuccessResponse[T any](data *T) ApiResponse[T] {
	return ApiResponse[T]{
		Success: true,
		Data: data,
	}
}