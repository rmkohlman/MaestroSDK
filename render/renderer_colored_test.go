package render

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/rmkohlman/MaestroSDK/colors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColoredRenderer_Name(t *testing.T) {
	r := NewColoredRenderer()
	assert.Equal(t, RendererColored, r.Name())
}

func TestColoredRenderer_SupportsColor(t *testing.T) {
	r := NewColoredRenderer()
	assert.True(t, r.SupportsColor())
}

func TestColoredRenderer_Render(t *testing.T) {
	r := NewColoredRenderer()

	t.Run("KeyValueData", func(t *testing.T) {
		var buf bytes.Buffer
		data := NewOrderedKeyValueData(
			KeyValue{Key: "Project", Value: "test"},
			KeyValue{Key: "Workspace", Value: "dev"},
		)

		err := r.Render(&buf, data, Options{})
		require.NoError(t, err)

		output := buf.String()
		// Content should be present (with ANSI codes)
		assert.Contains(t, output, "Project")
		assert.Contains(t, output, "test")
	})

	t.Run("TableData", func(t *testing.T) {
		var buf bytes.Buffer
		data := TableData{
			Headers: []string{"Name", "Status"},
			Rows: [][]string{
				{"proj1", "active"},
				{"proj2", "stopped"},
			},
		}

		err := r.Render(&buf, data, Options{})
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "Name")
		assert.Contains(t, output, "proj1")
	})

	t.Run("with title", func(t *testing.T) {
		var buf bytes.Buffer
		data := map[string]string{"key": "value"}

		err := r.Render(&buf, data, Options{Title: "Test Section"})
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "Test Section")
	})

	t.Run("empty state", func(t *testing.T) {
		var buf bytes.Buffer

		err := r.Render(&buf, nil, Options{
			Empty:        true,
			EmptyMessage: "No items found",
			EmptyHints:   []string{"Add an item"},
		})
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "No items found")
		assert.Contains(t, output, "Add an item")
	})
}

func TestColoredRenderer_RenderMessage(t *testing.T) {
	r := NewColoredRenderer()

	levels := []MessageLevel{
		LevelInfo,
		LevelSuccess,
		LevelWarning,
		LevelError,
		LevelProgress,
		LevelDebug,
	}

	for _, level := range levels {
		t.Run(string(level), func(t *testing.T) {
			var buf bytes.Buffer
			err := r.RenderMessage(&buf, Message{Level: level, Content: "test message"})
			require.NoError(t, err)
			assert.Contains(t, buf.String(), "test message")
		})
	}
}

func TestColoredRendererWithIcons(t *testing.T) {
	nerdIcons := NerdFontIcons()
	r := NewColoredRendererWithIcons(nerdIcons)

	var buf bytes.Buffer
	err := r.RenderMessage(&buf, Message{Level: LevelSuccess, Content: "done"})
	require.NoError(t, err)

	// Verify it renders (specific icon chars may not display in test output)
	assert.Contains(t, buf.String(), "done")
}

func TestIcons(t *testing.T) {
	t.Run("DefaultIcons", func(t *testing.T) {
		icons := DefaultIcons()
		assert.NotEmpty(t, icons.Success)
		assert.NotEmpty(t, icons.Warning)
		assert.NotEmpty(t, icons.Error)
		assert.NotEmpty(t, icons.Info)
		assert.NotEmpty(t, icons.Progress)
		assert.NotEmpty(t, icons.Bullet)
	})

	t.Run("NerdFontIcons", func(t *testing.T) {
		icons := NerdFontIcons()
		assert.NotEmpty(t, icons.Success)
		assert.NotEmpty(t, icons.Warning)
	})

	t.Run("PlainIcons", func(t *testing.T) {
		icons := PlainIcons()
		assert.Equal(t, "[OK]", icons.Success)
		assert.Equal(t, "[!]", icons.Warning)
		assert.Equal(t, "[X]", icons.Error)
	})
}

// ---------------------------------------------------------------------------
// Table styling tests (Issue #230)
// ---------------------------------------------------------------------------

func TestColoredRenderer_Table_HasBorders(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"NAME", "STATUS"},
		Rows: [][]string{
			{"alpha", "running"},
			{"bravo", "stopped"},
		},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	// Should contain box-drawing border characters
	assert.Contains(t, output, "┌", "should have top-left corner")
	assert.Contains(t, output, "┐", "should have top-right corner")
	assert.Contains(t, output, "└", "should have bottom-left corner")
	assert.Contains(t, output, "┘", "should have bottom-right corner")
	assert.Contains(t, output, "│", "should have vertical borders")
	assert.Contains(t, output, "─", "should have horizontal borders")
}

func TestColoredRenderer_Table_AlternatingRowBackgrounds(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"NAME", "VALUE"},
		Rows: [][]string{
			{"row0", "even"},
			{"row1", "odd"},
			{"row2", "even"},
			{"row3", "odd"},
		},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Find data row lines (skip top border, header, separator)
	var dataLines []string
	for _, line := range lines {
		if strings.Contains(line, "row0") || strings.Contains(line, "row1") ||
			strings.Contains(line, "row2") || strings.Contains(line, "row3") {
			dataLines = append(dataLines, line)
		}
	}
	require.Len(t, dataLines, 4, "should have 4 data rows")

	// Without ColorProvider, default mode uses ANSI codes (not truecolor)
	// Odd rows (index 1, 3) should have ANSI background code
	assert.Contains(t, dataLines[1], "\x1b[100m", "odd row should have ANSI background color")
	assert.Contains(t, dataLines[3], "\x1b[100m", "odd row should have ANSI background color")

	// Even rows (index 0, 2) should NOT have background ANSI code
	assert.NotContains(t, dataLines[0], "\x1b[100m", "even row should not have background color")
	assert.NotContains(t, dataLines[2], "\x1b[100m", "even row should not have background color")

	// No truecolor codes should be present without a ColorProvider
	assert.NotContains(t, dataLines[1], "\x1b[48;2;", "default mode should not use truecolor")
	assert.NotContains(t, dataLines[3], "\x1b[48;2;", "default mode should not use truecolor")
}

func TestColoredRenderer_Table_HeaderHasBackground(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"NAME", "STATUS"},
		Rows:    [][]string{{"test", "ok"}},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Header line should contain ANSI background code (not truecolor in default mode)
	var headerLine string
	for _, line := range lines {
		if strings.Contains(line, "NAME") && strings.Contains(line, "STATUS") {
			headerLine = line
			break
		}
	}
	require.NotEmpty(t, headerLine, "should find header line")
	assert.Contains(t, headerLine, "\x1b[100m", "header should have ANSI background color")
	assert.NotContains(t, headerLine, "\x1b[48;2;", "default mode should not use truecolor for header")
}

func TestColoredRenderer_Table_NOCOLORDisablesBackgrounds(t *testing.T) {
	// Set NO_COLOR env var
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	r := NewColoredRenderer()
	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"NAME", "STATUS"},
		Rows: [][]string{
			{"alpha", "running"},
			{"bravo", "stopped"},
		},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	// NO_COLOR should prevent all color ANSI codes
	assert.NotContains(t, output, "\x1b[48;2;", "NO_COLOR should disable truecolor backgrounds")
	assert.NotContains(t, output, "\x1b[100m", "NO_COLOR should disable ANSI backgrounds")
	assert.NotContains(t, output, "\x1b[90m", "NO_COLOR should disable ANSI border colors")
	// Should still have border characters (structural, not color)
	assert.Contains(t, output, "│", "borders should still appear")
	assert.Contains(t, output, "─", "horizontal rules should still appear")
}

func TestColoredRenderer_Table_WithColorProvider(t *testing.T) {
	r := NewColoredRenderer()
	provider := colors.NewDefaultColorProvider()
	ctx := colors.WithProvider(context.Background(), provider)

	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"NAME", "VALUE"},
		Rows: [][]string{
			{"key1", "val1"},
			{"key2", "val2"},
		},
	}

	err := r.RenderWithContext(ctx, &buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	// Should still have borders and data
	assert.Contains(t, output, "key1")
	assert.Contains(t, output, "val2")
	assert.Contains(t, output, "│")

	// With a ColorProvider, should use truecolor (38;2; for border, 48;2; for backgrounds)
	assert.Contains(t, output, "38;2;", "with ColorProvider, borders should use truecolor")
	assert.Contains(t, output, "48;2;", "with ColorProvider, backgrounds should use truecolor")

	// Should NOT contain ANSI-only codes (those are for default mode)
	assert.NotContains(t, output, "\x1b[100m", "with ColorProvider, should not use ANSI fallback bg")
	assert.NotContains(t, output, "\x1b[90m", "with ColorProvider, should not use ANSI fallback border")
}

func TestColoredRenderer_Table_ContentPreserved(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"ECOSYSTEM", "DOMAINS", "APPS"},
		Rows: [][]string{
			{"prod", "3", "5"},
			{"staging", "2", "3"},
		},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	// All headers present
	assert.Contains(t, output, "ECOSYSTEM")
	assert.Contains(t, output, "DOMAINS")
	assert.Contains(t, output, "APPS")
	// All data present
	assert.Contains(t, output, "prod")
	assert.Contains(t, output, "staging")
	assert.Contains(t, output, "3")
	assert.Contains(t, output, "5")
}

func TestColoredRenderer_Table_DefaultUsesANSICodes(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer
	data := TableData{
		Headers: []string{"COL1", "COL2"},
		Rows: [][]string{
			{"a", "b"},
			{"c", "d"},
			{"e", "f"},
		},
	}

	// Render without ColorProvider (default mode)
	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()

	// Default mode should use standard ANSI codes, NOT truecolor
	assert.NotContains(t, output, "\x1b[48;2;", "default mode must not use truecolor backgrounds")
	assert.NotContains(t, output, "\x1b[38;2;", "default mode must not use truecolor foregrounds")

	// Should contain ANSI border color code
	assert.Contains(t, output, "\x1b[90m", "default borders should use bright black fg (ANSI 90)")

	// Header and odd rows should use bright black background
	assert.Contains(t, output, "\x1b[100m", "default header/odd rows should use bright black bg (ANSI 100)")

	// Should still have box-drawing characters and data
	assert.Contains(t, output, "┌")
	assert.Contains(t, output, "│")
	assert.Contains(t, output, "COL1")
	assert.Contains(t, output, "c")
	assert.Contains(t, output, "f")
}

// ---------------------------------------------------------------------------
// Column constraint tests (Issue #258)
// ---------------------------------------------------------------------------

func TestColoredRenderer_Table_ColumnConstraints(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer

	longURL := "git@gitlab.ana.shawcable.net:access-network-automation/beansng/beans-ray-actorkit.git"
	data := TableData{
		Headers: []string{"NAME", "URL", "STATUS"},
		Rows: [][]string{
			{"short-name", longURL, "active"},
			{"another", "https://github.com/org/repo.git", "synced"},
		},
		Constraints: []ColumnConstraint{
			{MaxWidth: 15, Truncate: TruncEnd},
			{MaxWidth: 30, Truncate: TruncMiddle},
			{}, // no constraint on STATUS
		},
	}

	// Ensure NO_COLOR so we can inspect raw output easier
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()

	// The long URL should be middle-truncated (not appear in full)
	assert.NotContains(t, output, longURL, "full URL should not appear — should be truncated")

	// The truncated URL should contain "..." in the middle
	assert.Contains(t, output, "...", "truncated URL should contain ellipsis")

	// STATUS column should be unaffected
	assert.Contains(t, output, "active")
	assert.Contains(t, output, "synced")

	// Verify the URL truncation result fits in 30 chars
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		stripped := stripANSI(line)
		if strings.Contains(stripped, "git@gitlab") {
			// The cell containing the URL should have been truncated
			assert.NotContains(t, stripped, "beansng", "middle of URL should be removed")
			break
		}
	}
}

func TestColoredRenderer_Table_ConstraintsBackwardCompatible(t *testing.T) {
	r := NewColoredRenderer()

	// Render WITHOUT constraints
	var buf1 bytes.Buffer
	data1 := TableData{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"hello", "world"}},
	}

	// Render WITH nil constraints (same thing)
	var buf2 bytes.Buffer
	data2 := TableData{
		Headers:     []string{"A", "B"},
		Rows:        [][]string{{"hello", "world"}},
		Constraints: nil,
	}

	// Ensure NO_COLOR for deterministic comparison
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	err1 := r.Render(&buf1, data1, Options{Type: TypeTable})
	require.NoError(t, err1)

	err2 := r.Render(&buf2, data2, Options{Type: TypeTable})
	require.NoError(t, err2)

	// Output should be identical
	assert.Equal(t, buf1.String(), buf2.String(),
		"nil Constraints should produce identical output to missing Constraints")
}

func TestColoredRenderer_Table_MinWidthConstraint(t *testing.T) {
	r := NewColoredRenderer()
	var buf bytes.Buffer

	data := TableData{
		Headers: []string{"X", "Y"},
		Rows:    [][]string{{"a", "b"}},
		Constraints: []ColumnConstraint{
			{MinWidth: 20}, // force column X to be at least 20 wide
			{},
		},
	}

	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	err := r.Render(&buf, data, Options{Type: TypeTable})
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")

	// Find the header separator line — its segments tell us column widths.
	// The top border line uses ─ characters; count them for the first column.
	for _, line := range lines {
		stripped := stripANSI(line)
		if strings.Contains(stripped, "─") && strings.Contains(stripped, "┌") {
			// The first segment between ┌ and ┬ should be at least 20 + 2 (padding) = 22 chars of ─
			parts := strings.Split(stripped, "┬")
			if len(parts) >= 1 {
				firstSeg := strings.TrimPrefix(parts[0], "┌")
				// Count ─ characters
				dashCount := strings.Count(firstSeg, "─")
				assert.GreaterOrEqual(t, dashCount, 22,
					"first column should be at least 20 wide (+ 2 padding)")
			}
			break
		}
	}
}
