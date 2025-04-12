package sfqb

type Search interface {
	Fields() []string
	Term() string
}

type SearchData struct {
	term   string
	fields []string
}

func NewSearch(query string, fields ...string) Search {
	return &SearchData{term: query, fields: fields}
}

func (s *SearchData) Fields() []string {
	return s.fields
}

func (s *SearchData) Term() string {
	return s.term
}
