package models

type Order struct {
	Id         string `json:"id"`
	UserID     string `json:"user_id"`
	Number     string `json:"number"`
	Status     string `json:"status"`
	UploadedAt string `json:"uploaded_at"`
	Accrual    int    `json:"accrual"`
}

func (o *Order) Validate() []error {
	errorSlice := []error{}
	errorSlice = append(errorSlice, o.validateNumber()...)
	return errorSlice
}

func (o *Order) validateNumber() []error {
	errorSlice := []error{}
	return errorSlice
}
