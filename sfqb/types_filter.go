package sfqb

type FilterField struct {
	Name   string
	Method Method
	Value  interface{}
}

func NewFilterField(name string, method Method, value interface{}) FilterField {
	return FilterField{Name: name, Method: method, Value: value}
}
