package queryify

import (
	"errors"

	pb_common_filter "github.com/Kazzess/contracts/gen/go/common/filter/v1"
	"github.com/Kazzess/libraries/sfqb"
	"github.com/Kazzess/libraries/sfqb/sfqb_rqp"
)

const (
	minSearchSymbolsVal = 3
)

var ErrSearchMinSymbols = errors.New("not enough search symbols")

type Paginator interface {
	GetPagination() *pb_common_filter.Pagination
}

type Searcher interface {
	GetSearch() string
}

type Sorter interface {
	GetSort() *pb_common_filter.Sort
}

type FiltersConfig struct {
	paginator Paginator
	searcher  Searcher
	sorter    Sorter

	searchFields     []string
	minSearchSymbols int // this number from Policy

	defaultLimit int
	maxLimit     int
}

type FiltersOption func(*FiltersConfig)

// WithMinSearchSymbols sets the minimum search symbols for FiltersConfig.
func WithMinSearchSymbols(val int) FiltersOption {
	return func(f *FiltersConfig) {
		f.minSearchSymbols = val
	}
}

// WithPaginator sets the paginator and limits for FiltersConfig.
func WithPaginator(p Paginator, defaultLimit, maxLimit int) FiltersOption {
	return func(f *FiltersConfig) {
		f.paginator = p
		f.defaultLimit = defaultLimit
		f.maxLimit = maxLimit
	}
}

// WithSearcher sets the searcher for FiltersConfig.
func WithSearcher(s Searcher) FiltersOption {
	return func(f *FiltersConfig) {
		f.searcher = s
	}
}

// WithSorter sets the sorter for FiltersConfig.
func WithSorter(s Sorter) FiltersOption {
	return func(f *FiltersConfig) {
		f.sorter = s
	}
}

// WithSearchFields sets the search fields for FiltersConfig.
func WithSearchFields(val []string) FiltersOption {
	return func(f *FiltersConfig) {
		f.searchFields = val
	}
}

func NewFilters(options ...FiltersOption) (sfqb.SFQB, error) {
	config := &FiltersConfig{
		minSearchSymbols: minSearchSymbolsVal,
	}

	// Apply the options to the config
	for _, option := range options {
		option(config)
	}

	filters := sfqb_rqp.New()

	if config.paginator != nil && config.paginator.GetPagination() != nil {
		filters.SetLimit(int(config.paginator.GetPagination().GetLimit()))
		filters.SetOffset(int(config.paginator.GetPagination().GetOffset()))
	}

	if filters.Limit() == 0 {
		filters.SetLimit(config.defaultLimit)
	}

	if filters.Limit() > config.maxLimit {
		filters.SetLimit(config.maxLimit)
	}

	if config.sorter != nil && config.sorter.GetSort() != nil {
		filters.AddSortBy(config.sorter.GetSort().GetField(), config.sorter.GetSort().GetDesc())
	}

	if config.searcher != nil && config.searchFields != nil && config.searcher.GetSearch() != "" {
		if len(config.searcher.GetSearch()) < config.minSearchSymbols {
			return nil, ErrSearchMinSymbols
		}

		filters.SetSearch(sfqb.NewSearch(config.searcher.GetSearch(), config.searchFields...))
	}

	return filters, nil
}
