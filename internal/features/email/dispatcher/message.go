package dispatcher

type EmailMessage struct {
	To    string `json:"to"`
	FIO   string `json:"fio"`
	Token string `json:"token"`
	Retry int    `json:"retry"` //повторения
}
