package queryify

type Table struct {
	Schema string
	Name   string
	Alias  string
	PKName string
}

func (t *Table) String() string {
	return t.Schema + "." + t.Name
}

func (t *Table) From() string {
	return t.Schema + "." + t.Name + " AS " + t.Alias
}

func (t *Table) JoinTable(alias, field string) string {
	return t.From() + " ON " + t.Alias + "." + t.PKName + "=" + alias + "." + field
}

func (t *Table) JoinTablePK(jt *Table) string {
	if jt == nil {
		return ""
	}

	return t.From() + " ON " + t.Alias + "." + t.PKName + "=" + jt.Alias + "." + jt.PKName
}

func NewTable(schema, name, alias, pk string) *Table {
	return &Table{
		Schema: schema,
		Name:   name,
		Alias:  alias,
		PKName: pk,
	}
}
