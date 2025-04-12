package sfqb

type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}

type SFQB interface {
	AddFilter(fields ...FilterField)
	AddRawFilter(condition string)
	AddOrFields(fields ...FilterField)

	SetLimit(limit int)
	SetOffset(offset int)

	RemoveSort(field string)
	AddSortBy(field string, desc bool)
	AllSorts() []Sort

	SetSearch(search Search)
	HasSearch() bool
	Search() Search

	Limit() int
	Offset() int
	Order() string

	RemoveFilter(names ...string)
	GetFilter(name string) (FilterField, error)
	HasFilter(name string) bool
	AllFilters() []*FilterField
	String() string
	ReplaceFilterName(from, to string)
	ReplaceFilterValue(key string, to interface{})

	SetCustomWhere(where Sqlizer) error
	Where() string
	SQL(table string) string
	Args() []interface{}
}
