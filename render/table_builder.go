package render

import (
	"regexp"
	"unicode/utf8"
)

// ansiEscapeRE matches ANSI SGR escape sequences (e.g. \x1b[38;2;0;0;0m, \x1b[0m).
// These sequences control text color and styling but have zero display width.
var ansiEscapeRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI removes all ANSI SGR escape sequences from s, returning only the
// visible characters. This is used for accurate column-width calculation in
// table renderers when cells contain colored text.
func stripANSI(s string) string {
	return ansiEscapeRE.ReplaceAllString(s, "")
}

// displayWidth returns the number of visible characters in s after stripping
// ANSI escape sequences. It uses utf8.RuneCountInString so that multi-byte
// Unicode characters (like ●) are counted as one character each.
func displayWidth(s string) int {
	return utf8.RuneCountInString(stripANSI(s))
}

// padToWidth pads cell with trailing spaces so that the visible portion
// occupies exactly targetWidth characters. The returned string keeps the
// original ANSI codes intact but adds the right amount of trailing space.
func padToWidth(cell string, targetWidth int) string {
	visible := displayWidth(cell)
	if visible >= targetWidth {
		return cell
	}
	pad := targetWidth - visible
	buf := make([]byte, len(cell)+pad)
	copy(buf, cell)
	for i := len(cell); i < len(cell)+pad; i++ {
		buf[i] = ' '
	}
	return string(buf)
}

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
