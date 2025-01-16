package contract

type Plugin interface {
	Check() (result bool, err error)
}
