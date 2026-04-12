package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTableBuilder(t *testing.T) {
	tb := NewTableBuilder("NAME", "TYPE", "STATUS")
	td := tb.Build()

	assert.Equal(t, []string{"NAME", "TYPE", "STATUS"}, td.Headers)
	assert.Empty(t, td.Rows)
}

func TestTableBuilder_AddRow(t *testing.T) {
	tb := NewTableBuilder("NAME", "VALUE")
	tb.AddRow("foo", "bar")
	tb.AddRow("baz", "qux")

	td := tb.Build()
	assert.Len(t, td.Rows, 2)
	assert.Equal(t, []string{"foo", "bar"}, td.Rows[0])
	assert.Equal(t, []string{"baz", "qux"}, td.Rows[1])
}

func TestTableBuilder_Fluent(t *testing.T) {
	td := NewTableBuilder("A", "B").
		AddRow("1", "2").
		AddRow("3", "4").
		Build()

	assert.Len(t, td.Headers, 2)
	assert.Len(t, td.Rows, 2)
}

func TestTableBuilder_Len(t *testing.T) {
	tb := NewTableBuilder("X")
	assert.Equal(t, 0, tb.Len())
	tb.AddRow("a")
	assert.Equal(t, 1, tb.Len())
	tb.AddRow("b")
	assert.Equal(t, 2, tb.Len())
}

func TestTableBuilder_EmptyHeaders(t *testing.T) {
	tb := NewTableBuilder()
	td := tb.Build()
	assert.Empty(t, td.Headers)
	assert.Empty(t, td.Rows)
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string unchanged", "hello", 10, "hello"},
		{"exact length unchanged", "hello", 5, "hello"},
		{"long string truncated", "hello world", 8, "hello..."},
		{"empty string", "", 10, ""},
		{"maxLen less than 4", "hello", 3, "hello"},
		{"maxLen exactly 4", "hello", 4, "h..."},
		{"one char over", "abcdef", 5, "ab..."},
		{"unicode safe", "description text here", 15, "description ..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name   string
		hex    string
		wantR  int
		wantG  int
		wantB  int
		wantOK bool
	}{
		{"with hash", "#7aa2f7", 122, 162, 247, true},
		{"without hash", "7aa2f7", 122, 162, 247, true},
		{"black", "#000000", 0, 0, 0, true},
		{"white", "#ffffff", 255, 255, 255, true},
		{"empty string", "", 0, 0, 0, false},
		{"too short", "#fff", 0, 0, 0, false},
		{"invalid hex", "#gggggg", 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, ok := hexToRGB(tt.hex)
			assert.Equal(t, tt.wantOK, ok)
			if ok {
				assert.Equal(t, tt.wantR, r)
				assert.Equal(t, tt.wantG, g)
				assert.Equal(t, tt.wantB, b)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{"negative", -10, 0},
		{"zero", 0, 0},
		{"mid range", 128, 128},
		{"max", 255, 255},
		{"over max", 300, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, clamp(tt.in))
		})
	}
}

func TestTruncateMiddle(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			"short string unchanged",
			"hello", 10, "hello",
		},
		{
			"exact length unchanged",
			"hello", 5, "hello",
		},
		{
			"long URL truncated",
			"git@gitlab.ana.shawcable.net:access-network-automation/beansng/beans-ray-actorkit.git",
			30,
			// half=(30-3)/2=13, endLen=30-3-13=14
			"git@gitlab.an...y-actorkit.git",
		},
		{
			"minimum maxLen (5)",
			"abcdefghij", 5, "a...j",
		},
		{
			"maxLen less than 5 returns original",
			"abcdefghij", 4, "abcdefghij",
		},
		{
			"empty string",
			"", 10, "",
		},
		{
			"even split",
			"abcdefghijkl", 10, "abc...ijkl",
		},
		{
			"odd split gives extra char to end",
			"abcdefghijklm", 10, "abc...jklm",
		},
		{
			"maxLen equals 0",
			"hello world", 0, "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateMiddle(tt.input, tt.maxLen)
			assert.Equal(t, tt.want, got)
			// Verify truncated result does not exceed maxLen (when truncation applies)
			if tt.maxLen >= 5 && len(tt.input) > tt.maxLen {
				assert.LessOrEqual(t, len(got), tt.maxLen,
					"truncated result should not exceed maxLen")
			}
		})
	}
}

func TestApplyTruncation(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz"

	tests := []struct {
		name     string
		input    string
		maxLen   int
		strategy TruncateStrategy
		want     string
	}{
		{
			"TruncNone returns original",
			long, 10, TruncNone, long,
		},
		{
			"TruncEnd delegates to Truncate",
			long, 10, TruncEnd, Truncate(long, 10),
		},
		{
			"TruncMiddle delegates to TruncateMiddle",
			long, 10, TruncMiddle, TruncateMiddle(long, 10),
		},
		{
			"maxLen 0 returns original regardless of strategy",
			long, 0, TruncEnd, long,
		},
		{
			"negative maxLen returns original",
			long, -5, TruncMiddle, long,
		},
		{
			"short string unaffected by TruncEnd",
			"hi", 10, TruncEnd, "hi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyTruncation(tt.input, tt.maxLen, tt.strategy)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTableBuilder_SetConstraints(t *testing.T) {
	td := NewTableBuilder("NAME", "URL").
		SetConstraints(
			ColumnConstraint{MaxWidth: 20, Truncate: TruncEnd},
			ColumnConstraint{MaxWidth: 40, Truncate: TruncMiddle},
		).
		AddRow("short", "also-short").
		Build()

	assert.Len(t, td.Constraints, 2)
	assert.Equal(t, 20, td.Constraints[0].MaxWidth)
	assert.Equal(t, TruncEnd, td.Constraints[0].Truncate)
	assert.Equal(t, 40, td.Constraints[1].MaxWidth)
	assert.Equal(t, TruncMiddle, td.Constraints[1].Truncate)
}
