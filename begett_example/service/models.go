package service

type Company struct {
	UUID string `json:"uuid"`
	Name string	`json:"name"`
}

type Employee struct {
	UUID string `json:"uuid"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Company *Company `json:"company"`
}