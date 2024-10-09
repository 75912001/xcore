package parameters

type defaultParameters struct {
	parameters []interface{}
}

func NewDefaultParameters() IParameters {
	return &defaultParameters{
		parameters: make([]interface{}, 0),
	}
}

func (p *defaultParameters) Set(parameters ...interface{}) {
	p.reset()
	for _, v := range parameters {
		p.parameters = append(p.parameters, v)
	}
}

func (p *defaultParameters) Get() []interface{} {
	return p.parameters
}

func (p *defaultParameters) Add(parameters ...interface{}) {
	for _, v := range parameters {
		p.parameters = append(p.parameters, v)
	}
}

func (p *defaultParameters) reset() {
	p.parameters = make([]interface{}, 0)
}
