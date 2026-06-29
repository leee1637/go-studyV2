package domain

type Teacher struct {
	ID          int      `json:"id"`
	FIO         string   `json:"fio"`
	PhoneNumber string   `json:"phone_number"`
	Group       []string `json:"group"`
}
