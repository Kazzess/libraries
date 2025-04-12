package sfqb_rqp

import (
	"testing"

	"github.com/Kazzess/libraries/sfqb"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

func TestReplaceFilterValue(t *testing.T) {
	cases := []struct {
		title    string
		name     string
		value    interface{}
		newValue interface{}
	}{
		{
			title:    "string to string",
			name:     "foo",
			value:    "bar",
			newValue: "buzz",
		},
		{
			title:    "string to int",
			name:     "foo",
			value:    100,
			newValue: "buzz",
		},
		{
			title:    "bool to string",
			name:     "foo",
			value:    true,
			newValue: "buzz",
		},
		{
			title:    "float64 to int",
			name:     "foo",
			value:    100.34,
			newValue: 100,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			fltrs := New()
			fltrs.AddFilter(sfqb.NewFilterField(c.name, sfqb.EQ, c.value))
			filter, _ := fltrs.GetFilter(c.name)
			assert.Equal(t, filter.Value, c.value)

			fltrs.ReplaceFilterValue(c.name, c.newValue)
			filter, _ = fltrs.GetFilter(c.name)
			assert.Equal(t, filter.Value, c.newValue)
		})
	}
}

func TestAddOrFields(t *testing.T) {
	fltrs := New()

	fltrs.AddOrFields([]sfqb.FilterField{
		{Name: "foo", Method: sfqb.EQ, Value: "bar"},
	}...)

	assert.Equal(t, 1, len(fltrs.AllFilters()))

	fltrs2 := New()

	fltrs2.AddOrFields([]sfqb.FilterField{
		{Name: "foo", Method: sfqb.EQ, Value: "bar"},
		{Name: "foo1", Method: sfqb.EQ, Value: "bar1"},
	}...)

	assert.Equal(t, 2, len(fltrs2.AllFilters()))
}

func TestCustomWhere(t *testing.T) {
	var cond sfqb.Sqlizer
	cond = squirrel.And{
		squirrel.Eq{"foo": "bar"},
		squirrel.Or{
			squirrel.Eq{"foo1": "bar1"},
			squirrel.Eq{"foo2": "bar2"},
		},
		squirrel.Or{
			squirrel.Eq{"foo1": "bar1"},
			squirrel.Eq{"foo2": "bar2"},
		},
	}

	fltrs := New()

	fltrs.AddFilter(sfqb.NewFilterField("ffff", sfqb.EQ, "aaaa"))

	err := fltrs.SetCustomWhere(cond)
	if err != nil {
		panic(err)
	}

	where := fltrs.Where()
	if where != "(foo = ? AND (foo1 = ? OR foo2 = ?) AND (foo1 = ? OR foo2 = ?))" {
		panic(where)
	}

	args := fltrs.Args()
	for i, v := range args {
		if i == 0 && v != "bar" {
			panic(v)
		}
		if i == 1 && v != "bar1" {
			panic(v)
		}
		if i == 2 && v != "bar2" {
			panic(v)
		}
		if i == 3 && v != "bar1" {
			panic(v)
		}
		if i == 4 && v != "bar2" {
			panic(v)
		}
	}
}

func TestStringM(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		filters  func() sfqb.SFQB
	}{
		{
			name:     "filters only and",
			expected: "foo = ? AND buz = ?barbaz00",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				return fltrs
			},
		},
		{
			name:     "custom where same filters",
			expected: "(foo = ? AND buz = ?)barbaz00",
			filters: func() sfqb.SFQB {
				fltrs := New()
				err := fltrs.SetCustomWhere(squirrel.And{squirrel.Eq{"foo": "bar"}, squirrel.Eq{"buz": "baz"}})
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return fltrs
			},
		},
		{
			name:     "filters and and or",
			expected: "(foo = ? OR buz = ?)barbaz00",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddOrFields(
					sfqb.NewFilterField("foo", sfqb.EQ, "bar"),
					sfqb.NewFilterField("buz", sfqb.EQ, "baz"),
				)
				return fltrs
			},
		},
		{
			name:     "filters and sorts",
			expected: "foo = ? AND buz = ?barbazfoo, bar DESCfooASCbarDESC00",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				return fltrs
			},
		},
		{
			name:     "filters and sorts and limit",
			expected: "foo = ? AND buz = ?barbazfoo, bar DESCfooASCbarDESC100",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				fltrs.SetLimit(10)
				return fltrs
			},
		},
		{
			name:     "filters and sorts and offset",
			expected: "foo = ? AND buz = ?barbazfoo, bar DESCfooASCbarDESC020",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				fltrs.SetOffset(20)
				return fltrs
			},
		},
		{
			name:     "filters and sorts and limit and offset",
			expected: "foo = ? AND buz = ?barbazfoo, bar DESCfooASCbarDESC1020",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				fltrs.SetLimit(10)
				fltrs.SetOffset(20)
				return fltrs
			},
		},
		{
			name:     "sorts and limit and offset and custom_where and search",
			expected: "(foo = ? AND buz = ?)barbazfoo, bar DESCfooASCbarDESC1020",
			filters: func() sfqb.SFQB {
				fltrs := New()
				err := fltrs.SetCustomWhere(squirrel.And{squirrel.Eq{"foo": "bar"}, squirrel.Eq{"buz": "baz"}})
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				fltrs.SetLimit(10)
				fltrs.SetOffset(20)
				return fltrs
			},
		},
		{
			name:     "filters and sorts and limit and offset and search and raw_filter",
			expected: "foo = ? AND buz = ?barbazfoo, bar DESCfooASCbarDESCatmyhorsemyhorseisamazinglook1020",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				fltrs.SetLimit(10)
				fltrs.SetOffset(20)
				fltrs.SetSearch(sfqb.NewSearch("look", "at", "my", "horse", "my", "horse", "is", "amazing"))
				return fltrs
			},
		},
		{
			name:     "filters and sorts and limit and offset and custom_where and search and raw_filter",
			expected: "foo = ? AND buz = ? AND yakzupzup=123barbazfoo, bar DESCfooASCbarDESCatmyhorsemyhorseisamazinglook1020",
			filters: func() sfqb.SFQB {
				fltrs := New()
				fltrs.AddFilter(sfqb.NewFilterField("foo", sfqb.EQ, "bar"))
				fltrs.AddFilter(sfqb.NewFilterField("buz", sfqb.EQ, "baz"))
				fltrs.AddRawFilter("yakzupzup=123")
				fltrs.AddSortBy("foo", false)
				fltrs.AddSortBy("bar", true)
				fltrs.SetLimit(10)
				fltrs.SetOffset(20)
				fltrs.SetSearch(sfqb.NewSearch("look", "at", "my", "horse", "my", "horse", "is", "amazing"))
				return fltrs
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.filters().String()
			assert.Equalf(t, c.expected, result, "expected %s, got %s", c.expected, result)
		})
	}
}
