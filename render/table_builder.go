package render

// TableBuilder provides a fluent API for constructing TableData.
// It eliminates the boilerplate of manually building header/row slices
// that is duplicated across CLI output functions.
//
// Example:
//
//	tb := render.NewTableBuilder("NAME", "TYPE", "DESCRIPTION")
//	for _, item := range items {
//	    tb.AddRow(item.Name, item.Type, item.Desc)
//	}
//	return render.OutputWith(format, tb.Build(), render.Options{Type: render.TypeTable})
type TableBuilder struct {
	headers []string
	rows    [][]string
}

// NewTableBuilder creates a new TableBuilder with the given column headers.
func NewTableBuilder(headers ...string) *TableBuilder {
	return &TableBuilder{
		headers: headers,
		rows:    make([][]string, 0),
	}
}

// AddRow appends a row of values to the table.
// The number of values should match the number of headers.
func (tb *TableBuilder) AddRow(values ...string) *TableBuilder {
	tb.rows = append(tb.rows, values)
	return tb
}

// Build returns the constructed TableData.
func (tb *TableBuilder) Build() TableData {
	return TableData{
		Headers: tb.headers,
		Rows:    tb.rows,
	}
}

// Len returns the number of rows added so far.
func (tb *TableBuilder) Len() int {
	return len(tb.rows)
}

// Truncate truncates s to maxLen by keeping the first (maxLen-3) characters
// and suffixing with "...". If len(s) <= maxLen the original string is returned.
// This is a common pattern for description columns in table output.
func Truncate(s string, maxLen int) string {
	if maxLen < 4 {
		return s
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
