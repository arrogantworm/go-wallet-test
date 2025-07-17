package handler

type (
	ErrorRes struct {
		Message string `json:"error"`
	}

	SuccessRes struct {
		Message string `json:"success"`
	}
)
