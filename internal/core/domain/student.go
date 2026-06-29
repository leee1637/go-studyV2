package domain

type Student struct {
	ID          int     `json"id"`
	FIO         string  `json"fio"`
	Group       string  `json"group"`
	PhoneNumber *string `json"phone_number"`
}
