package ports

type Tgb interface {
	Send()
	Empty()
	Handler()
}
