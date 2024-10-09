package parameters

type IParameters interface {
	Set(parameters ...interface{})
	Get() []interface{}
	Add(parameter ...interface{})
}
