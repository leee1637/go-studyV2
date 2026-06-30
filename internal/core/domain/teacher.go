package domain

type Teacher struct {
	ID          int      `json:"id"`
	FIO         string   `json:"fio"`
	GroupName   []string `json:"group"`
	PhoneNumber string   `json:"phone_number"`
}
