package dispatcher

type EmailDispatcher interface {
	Send(to, fio, token string) error
}
