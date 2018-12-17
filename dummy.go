package gounity

var (
	typeNameDUMMY   = "dummy"
	typeFieldsDUMMY = "fields"
)

type DUMMY struct {
	Resource
	Id   string
	Name string
}

type DUMMYOperator interface {
	genDUMMYOperator
}
