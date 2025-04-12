package sfqb

type QuerySpecification struct {
	IgnoreErrors bool
	SingleFields []QuerySingleField
	ArrayFields  []QueryArrayField
	SortFields   QuerySortFields
}

type QueryValidationFunc func(value interface{}) error

type QuerySortField string
type QuerySortFields []QuerySortField

func NewQuerySortFields(fields ...QuerySortField) QuerySortFields {
	return fields
}

func (sf QuerySortFields) Repr() []interface{} {
	res := make([]interface{}, len(sf))
	for i, field := range sf {
		res[i] = string(field)
	}
	return res
}

type FieldType string

const (
	StringField FieldType = "string"
	IntField    FieldType = "int"
	BoolField   FieldType = "bool"
)

func (d FieldType) String() string {
	switch d {
	case StringField, IntField, BoolField:
		return string(d)
	}

	return ""
}

type QuerySingleField struct {
	Name       string
	DType      FieldType
	Validation QueryValidationFunc
	Required   bool
}

// NewDateQuerySingleField add date single query params field
func NewDateQuerySingleField(name string, dType FieldType, format string) QuerySingleField {
	return QuerySingleField{DType: dType, Name: name, Validation: DateFormat(format)}
}

// NewQuerySingleField add single query params field. If no validation require - use nil.
func NewQuerySingleField(name string, dType FieldType, validation QueryValidationFunc) QuerySingleField {
	return QuerySingleField{DType: dType, Name: name, Validation: validation}
}

// NewRequiredQuerySingleField add single query params field. If no validation require - use nil.
func NewRequiredQuerySingleField(name string, dType FieldType, validation QueryValidationFunc) QuerySingleField {
	qsf := NewQuerySingleField(name, dType, validation)
	qsf.Required = true
	return qsf
}

type QueryArrayField struct {
	Name   string
	DType  FieldType
	Values []interface{}
}

func NewQueryArrayField(name string, dType FieldType, values ...interface{}) QueryArrayField {
	return QueryArrayField{DType: dType, Name: name, Values: values}
}
