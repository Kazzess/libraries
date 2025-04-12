package sfqb

const NULL = "NULL"

type Method string

// Compare methods:
var (
	EQ     Method = "EQ"
	NE     Method = "NE"
	GT     Method = "GT"
	LT     Method = "LT"
	GTE    Method = "GTE"
	LTE    Method = "LTE"
	LIKE   Method = "LIKE"
	ILIKE  Method = "ILIKE"
	NLIKE  Method = "NLIKE"
	NILIKE Method = "NILIKE"
	IS     Method = "IS"
	NOT    Method = "NOT"
	IN     Method = "IN"
	NIN    Method = "NIN"
	Range  Method = "RANGE"
	raw    Method = "raw" // internal usage
)

var (
	translateMethods map[Method]string = map[Method]string{
		EQ:     "=",
		NE:     "!=",
		GT:     ">",
		LT:     "<",
		GTE:    ">=",
		LTE:    "<=",
		LIKE:   "LIKE",
		ILIKE:  "ILIKE",
		NLIKE:  "NOT LIKE",
		NILIKE: "NOT ILIKE",
		IS:     "IS",
		NOT:    "IS NOT",
		IN:     "IN",
		NIN:    "NOT IN",
		Range:  "RANGE",
	}
)
