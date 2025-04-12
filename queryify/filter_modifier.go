package queryify

import (
	"fmt"
	"strings"

	"github.com/Kazzess/libraries/sfqb"
)

const length = 2
const sep = "."

// ReplaceTableToAlias replaces the table name to alias in the filters.
func ReplaceTableToAlias(filters sfqb.SFQB, tables ...*Table) {
	if filters == nil {
		return
	}

	allFilters := filters.AllFilters()
	for i := range allFilters {
		fltr := allFilters[i]

		replaceToAlias(filters, fltr.Name, true, tables...)
	}

	sorts := filters.AllSorts()
	for i := range sorts {
		srt := sorts[i]

		replaceToAlias(filters, srt.By, false, tables...)
	}
}

// replaceToAlias replaces the table name to alias in the filters.
func replaceToAlias(filters sfqb.SFQB, name string, filter bool, tables ...*Table) {
	split := strings.Split(name, sep)

	if len(split) != length {
		if filter {
			filters.RemoveFilter(name)
		} else {
			filters.RemoveSort(name)
		}

		return
	}

	entity := split[0]
	field := split[1]

	for _, table := range tables {
		if table.Name == entity {
			filters.ReplaceFilterName(name, table.Alias+"."+field)
			break
		}
	}

}

// ApplySearchFilters applies the search filters to the filters.
/*
   - textFormat: mask for the field name.
	   - Example for Postgres: "%s::text"
   - formatter: mask for the value.
	   - Example for Postgres: "%%%s%%"
*/
func ApplySearchFilters(filters sfqb.SFQB, textFormat, formatter string) {
	if filters == nil || !filters.HasSearch() {
		return
	}

	search := filters.Search()

	searchFilters := make([]sfqb.FilterField, 0)

	for _, field := range search.Fields() {
		fieldNameFormatted := fmt.Sprintf(textFormat, field)
		fieldValueFormatted := fmt.Sprintf(formatter, search.Term())

		searchFilters = append(searchFilters, sfqb.FilterField{
			Name:   fieldNameFormatted,
			Method: sfqb.ILIKE,
			Value:  fieldValueFormatted,
		})
	}

	filters.AddOrFields(searchFilters...)
}

// ReplaceFilterLike replaces the filter value to like format.
func ReplaceFilterLike(filters sfqb.SFQB, formatter string) {
	if filters == nil {
		return
	}

	for _, filter := range filters.AllFilters() {

		switch filter.Method {
		case sfqb.ILIKE:
			val, ok := filter.Value.(string)
			if !ok {
				continue
			}

			fieldValueFormatted := fmt.Sprintf(formatter, val)

			filters.ReplaceFilterValue(filter.Name, fieldValueFormatted)
		case sfqb.LIKE:
			val, ok := filter.Value.(string)
			if !ok {
				continue
			}

			fieldValueFormatted := fmt.Sprintf(formatter, val)

			filters.ReplaceFilterValue(filter.Name, fieldValueFormatted)

		}
	}
}
