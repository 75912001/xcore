package parameters

type IParameters interface {
	Set(parameters ...interface{})
	Get() []interface{}
}

type parameters struct {
	parameters []interface{}
}

func NewParameters() IParameters {
	return &parameters{
		parameters: make([]interface{}, 0),
	}
}

func (p *parameters) Set(parameters ...interface{}) {
	p.reset()
	for _, v := range parameters {
		p.parameters = append(p.parameters, v)
	}
}

func (p *parameters) Get() []interface{} {
	return p.parameters
}

func (p *parameters) reset() {
	p.parameters = make([]interface{}, 0)
}
