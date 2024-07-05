package service

type Car struct {
	ID             int    `json:"id"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}
