package render

import (
	"context"
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// defaultTableStyles provides Catppuccin Mocha-inspired defaults when no
// TableStyleProvider is injected via context.
var defaultTableStyles = struct {
	header TableCellStyle
	cell   TableCellStyle
	border TableCellStyle
}{
	header: TableCellStyle{FG: "#CBA6F7", Bold: true},
	cell:   TableCellStyle{FG: "#CDD6F4"},
	border: TableCellStyle{FG: "#585B70"},
}

// PrettyRenderer produces bordered table output using charmbracelet/lipgloss/table.
// For non-table data it delegates to ColoredRenderer.
type PrettyRenderer struct {
	colored *ColoredRenderer
}

// NewPrettyRenderer creates a new PrettyRenderer.
func NewPrettyRenderer() *PrettyRenderer {
	return &PrettyRenderer{
		colored: NewColoredRenderer(),
	}
}

// Name returns the renderer identifier.
func (r *PrettyRenderer) Name() RendererName {
	return RendererPretty
}

// SupportsColor returns true — this renderer uses colors.
func (r *PrettyRenderer) SupportsColor() bool {
	return true
}

// Render outputs data; delegates to RenderWithContext with a background context.
func (r *PrettyRenderer) Render(w io.Writer, data any, opts Options) error {
	return r.RenderWithContext(context.Background(), w, data, opts)
}

// RenderWithContext outputs data with optional style injection via context.
func (r *PrettyRenderer) RenderWithContext(ctx context.Context, w io.Writer, data any, opts Options) error {
	switch v := data.(type) {
	case TableData:
		return r.renderBorderedTable(ctx, w, v)
	default:
		// Delegate non-table types to the ColoredRenderer
		return r.colored.RenderWithContext(ctx, w, data, opts)
	}
}

// RenderMessage delegates to the ColoredRenderer.
func (r *PrettyRenderer) RenderMessage(w io.Writer, msg Message) error {
	return r.colored.RenderMessage(w, msg)
}

// RenderMessageWithContext delegates to the ColoredRenderer.
func (r *PrettyRenderer) RenderMessageWithContext(ctx context.Context, w io.Writer, msg Message) error {
	return r.colored.RenderMessageWithContext(ctx, w, msg)
}

// resolveStyles returns header/cell/border styles from context or defaults.
func (r *PrettyRenderer) resolveStyles(ctx context.Context) (header, cell, border TableCellStyle) {
	if provider, ok := TableStyleProviderFromContext(ctx); ok {
		return provider.HeaderStyle(), provider.CellStyle(), provider.BorderStyle()
	}
	return defaultTableStyles.header, defaultTableStyles.cell, defaultTableStyles.border
}

// renderBorderedTable builds a lipgloss/table with box-drawing borders.
func (r *PrettyRenderer) renderBorderedTable(ctx context.Context, w io.Writer, td TableData) error {
	headerStyle, cellStyle, borderStyle := r.resolveStyles(ctx)

	// Pre-truncate cell values if constraints are provided.
	// lipgloss/table manages its own column widths, so we truncate the data
	// before handing it off.
	headers := td.Headers
	rows := td.Rows
	if len(td.Constraints) > 0 {
		// Truncate headers
		truncHeaders := make([]string, len(headers))
		for i, h := range headers {
			if i < len(td.Constraints) && td.Constraints[i].MaxWidth > 0 && len(h) > td.Constraints[i].MaxWidth {
				truncHeaders[i] = ApplyTruncation(h, td.Constraints[i].MaxWidth, td.Constraints[i].Truncate)
			} else {
				truncHeaders[i] = h
			}
		}
		headers = truncHeaders

		// Truncate row cells
		truncRows := make([][]string, len(rows))
		for ri, row := range rows {
			truncRow := make([]string, len(row))
			for ci, cell := range row {
				if ci < len(td.Constraints) && td.Constraints[ci].MaxWidth > 0 && len(cell) > td.Constraints[ci].MaxWidth {
					truncRow[ci] = ApplyTruncation(cell, td.Constraints[ci].MaxWidth, td.Constraints[ci].Truncate)
				} else {
					truncRow[ci] = cell
				}
			}
			truncRows[ri] = truncRow
		}
		rows = truncRows
	}

	// Build lipgloss styles
	hdrLipgloss := lipgloss.NewStyle().
		Foreground(lipgloss.Color(headerStyle.FG)).
		Bold(headerStyle.Bold).
		Padding(0, 1)

	cellLipgloss := lipgloss.NewStyle().
		Foreground(lipgloss.Color(cellStyle.FG)).
		Padding(0, 1)

	borderLipgloss := lipgloss.NewStyle().
		Foreground(lipgloss.Color(borderStyle.FG))

	// Create the table
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(borderLipgloss).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return hdrLipgloss
			}
			return cellLipgloss
		})

	fmt.Fprintln(w, t.Render())
	return nil
}
