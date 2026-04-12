package render

import (
	"fmt"
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
	headers     []string
	rows        [][]string
	constraints []ColumnConstraint
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

// SetConstraints sets column constraints for the table.
// The number of constraints should match the number of headers.
func (tb *TableBuilder) SetConstraints(constraints ...ColumnConstraint) *TableBuilder {
	tb.constraints = constraints
	return tb
}

// Build returns the constructed TableData.
func (tb *TableBuilder) Build() TableData {
	return TableData{
		Headers:     tb.headers,
		Rows:        tb.rows,
		Constraints: tb.constraints,
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

// TruncateMiddle shortens s by preserving the start and end of the string and
// replacing the middle with "...". This is useful for URLs and file paths where
// both the prefix (protocol/host) and suffix (filename) carry meaning.
//
// Example: "git@gitlab.ana.shawcable.net:access-network-automation/beansng/beans-ray-actorkit.git"
//
//	→ "git@gitlab.ana...actorkit.git"  (with maxLen=30)
func TruncateMiddle(s string, maxLen int) string {
	if maxLen < 5 || len(s) <= maxLen {
		return s
	}
	// Split the budget evenly between start and end, minus the 3-char "..."
	half := (maxLen - 3) / 2
	// When maxLen-3 is odd, give the extra char to the end portion
	endLen := maxLen - 3 - half
	return s[:half] + "..." + s[len(s)-endLen:]
}

// ApplyTruncation dispatches to the appropriate truncation function based on
// the given strategy. If strategy is TruncateNone or maxLen <= 0, the original
// string is returned unchanged.
func ApplyTruncation(s string, maxLen int, strategy TruncateStrategy) string {
	if maxLen <= 0 {
		return s
	}
	switch strategy {
	case TruncEnd:
		return Truncate(s, maxLen)
	case TruncMiddle:
		return TruncateMiddle(s, maxLen)
	default:
		// TruncNone — return as-is
		return s
	}
}

// hexToRGB parses a hex color string (e.g., "#7aa2f7" or "7aa2f7") into
// its red, green, blue components. Returns ok=false if parsing fails.
func hexToRGB(hex string) (r, g, b int, ok bool) {
	if len(hex) == 0 {
		return 0, 0, 0, false
	}
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return 0, 0, 0, false
	}
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return 0, 0, 0, false
	}
	return r, g, b, true
}

// clamp restricts an integer to the 0–255 range for ANSI color values.
func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
