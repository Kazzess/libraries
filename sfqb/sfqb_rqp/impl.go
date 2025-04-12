package sfqb_rqp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Kazzess/libraries/sfqb"
	"github.com/Kazzess/libraries/utils/random"
	rqp "github.com/timsolov/rest-query-parser"
)

type SFQB struct {
	*rqp.Query
	search sfqb.Search

	isCustomWhere   bool
	customWhere     string
	customWhereArgs []interface{}
}

func newRqpSQFB(query *rqp.Query) *SFQB {
	return &SFQB{Query: query}
}

func New() sfqb.SFQB {
	return newRqpSQFB(rqp.New())
}

func NewParseURL(query url.Values, spec sfqb.QuerySpecification) (sfqb.SFQB, error) {
	q := rqp.New().SetUrlQuery(query)

	for _, field := range spec.SingleFields {
		s := field.Name
		switch s {
		case "offset", "limit":
		default:
			s += ":" + field.DType.String()
		}

		if field.Required {
			s += ":required"
		}
		q = q.AddValidation(s, rqp.ValidationFunc(field.Validation))
	}

	for _, field := range spec.ArrayFields {
		q = q.AddValidation(field.Name+":"+field.DType.String(), rqp.In(field.Values...))
	}

	q = q.AddValidation("sort", rqp.In(spec.SortFields.Repr()...))

	for _, field := range spec.SortFields {
		q = q.AddSortBy(string(field), false)
	}

	res := newRqpSQFB(q)

	err := res.Query.Parse()
	if err != nil && !spec.IgnoreErrors {
		return nil, err
	}

	return res, nil
}

func (r *SFQB) SetSearch(search sfqb.Search) {
	r.search = search
}

func (r *SFQB) HasSearch() bool {
	return r.search != nil
}

func (r *SFQB) Search() sfqb.Search {
	return r.search
}

func (r *SFQB) SetCustomWhere(where sfqb.Sqlizer) error {
	if where == nil {
		r.isCustomWhere = false
		return nil
	}

	sql, args, err := where.ToSql()
	if err != nil {
		return err
	}

	r.customWhere = sql
	r.customWhereArgs = args

	r.isCustomWhere = true
	return nil
}

func (r *SFQB) AddFilter(fields ...sfqb.FilterField) {
	for _, field := range fields {
		r.Query.AddFilter(field.Name, rqp.Method(field.Method), field.Value)
	}
}

func (r *SFQB) AddRawFilter(condition string) {
	r.Query = r.Query.AddFilterRaw(condition)
}

func (r *SFQB) AddOrFields(fields ...sfqb.FilterField) {
	var addedTestFilter bool
	var testFilterName string

	r.Query.AddORFilters(func(q *rqp.Query) {
		for _, field := range fields {
			q.AddFilter(field.Name, rqp.Method(field.Method), field.Value)
		}
		if len(fields) == 1 {
			// this is because method AddORFilters in rest-query-parser doesn't support adding a single filter
			addedTestFilter = true
			testFilterName = fields[0].Name + random.RandomString(10)
			q.AddFilter(testFilterName, rqp.EQ, "test")
		}
	})

	if addedTestFilter {
		r.Query.RemoveFilter(testFilterName)
	}
}

func (r *SFQB) SetLimit(limit int) {
	r.Query.SetLimit(limit)
}

func (r *SFQB) SetOffset(offset int) {
	r.Query.SetOffset(offset)
}

func (r *SFQB) AddSortBy(field string, desc bool) {
	r.Query.AddSortBy(field, desc)
}

func (r *SFQB) RemoveSort(field string) {
	sorts := make([]rqp.Sort, len(r.Query.Sorts))

	for _, s := range r.Query.Sorts {
		if s.By == field {
			continue
		}

		sorts = append(sorts, rqp.Sort{
			By:   s.By,
			Desc: s.Desc,
		})
	}

	r.Query.Sorts = sorts
}

func (r *SFQB) RemoveFilter(names ...string) {
	for _, field := range names {
		r.Query.RemoveFilter(field)
	}
}

func (r *SFQB) HasFilter(name string) bool {
	return r.Query.HaveFilter(name)
}

func (r *SFQB) GetFilter(name string) (sfqb.FilterField, error) {
	filter, err := r.Query.GetFilter(name)
	if err != nil {
		return sfqb.FilterField{}, err
	}
	return sfqb.FilterField{
		Name:   filter.Name,
		Method: sfqb.Method(filter.Method),
		Value:  filter.Value,
	}, nil
}

func (r *SFQB) Order() string {
	return r.Query.Order()
}

func (r *SFQB) Limit() int {
	return r.Query.Limit
}

func (r *SFQB) Offset() int {
	return r.Query.Offset
}

func (r *SFQB) ReplaceFilterName(from, to string) {
	r.Query.ReplaceNames(rqp.Replacer{from: to})
}

func (r *SFQB) AllSorts() []sfqb.Sort {
	sorts := r.Query.Sorts

	res := make([]sfqb.Sort, len(sorts))

	for i, sort := range sorts {
		res[i] = sfqb.Sort{
			By:   sort.By,
			Desc: sort.Desc,
		}
	}

	return res
}

func (r *SFQB) AllFilters() []*sfqb.FilterField {
	rawFilters := r.Query.Filters

	res := make([]*sfqb.FilterField, len(rawFilters))

	for i, filter := range rawFilters {
		res[i] = &sfqb.FilterField{
			Name:   filter.Name,
			Method: sfqb.Method(filter.Method),
			Value:  filter.Value,
		}
	}

	return res
}

func (r *SFQB) String() string {
	var strBuilder strings.Builder

	strBuilder.WriteString(r.Where())
	for _, arg := range r.Args() {
		switch v := arg.(type) {
		case string:
			strBuilder.WriteString(v)
		case float64:
			strBuilder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		case int:
			strBuilder.WriteString(strconv.Itoa(v))
		case uint:
			strBuilder.WriteString(strconv.FormatUint(uint64(v), 10))
		default:
			strBuilder.WriteString(fmt.Sprintf("%v", arg))
		}
	}

	strBuilder.WriteString(r.Order())

	for _, sort := range r.AllSorts() {
		strBuilder.WriteString(sort.By)
		if sort.Desc {
			strBuilder.WriteString("DESC")
		} else {
			strBuilder.WriteString("ASC")
		}
	}

	if r.HasSearch() {
		for _, field := range r.search.Fields() {
			strBuilder.WriteString(field)
		}

		strBuilder.WriteString(r.search.Term())
	}

	strBuilder.WriteString(strconv.Itoa(r.Query.Limit))
	strBuilder.WriteString(strconv.Itoa(r.Query.Offset))

	return strBuilder.String()
}

func (r *SFQB) Where() string {
	if r.isCustomWhere {
		return r.customWhere
	}

	return r.Query.Where()
}

func (r *SFQB) SQL(table string) string {
	return r.Query.SQL(table)
}

func (r *SFQB) Args() []interface{} {
	if r.isCustomWhere {
		return r.customWhereArgs
	}

	return r.Query.Args()
}

// ReplaceFilterValue replaces the filter value from the given key.
func (r *SFQB) ReplaceFilterValue(key string, to interface{}) {
	for i, filter := range r.Query.Filters {
		if filter.Name == key {
			r.Query.Filters[i].Value = to
		}
	}
}
