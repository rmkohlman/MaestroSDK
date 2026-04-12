package render

import (
	"context"
	"fmt"
	"io"
	"strings"
)

// CompactRenderer is like ColoredRenderer but more condensed.
// Good for smaller terminals or when you want less visual noise.
type CompactRenderer struct {
	*ColoredRenderer
}

// NewCompactRenderer creates a new compact renderer
func NewCompactRenderer() *CompactRenderer {
	return &CompactRenderer{
		ColoredRenderer: NewColoredRenderer(),
	}
}

// Name returns the renderer identifier
func (r *CompactRenderer) Name() RendererName {
	return RendererCompact
}

// Render outputs data in compact format
func (r *CompactRenderer) Render(w io.Writer, data any, opts Options) error {
	// CompactRenderer delegates to ColoredRenderer - no context needed for compatibility
	return r.RenderWithContext(context.Background(), w, data, opts)
}

// RenderWithContext outputs data in compact format with context for theming
func (r *CompactRenderer) RenderWithContext(ctx context.Context, w io.Writer, data any, opts Options) error {
	styles := r.getStyles(ctx)

	// Handle empty state
	if opts.Empty {
		if opts.EmptyMessage != "" {
			fmt.Fprintln(w, styles.muted.Render(opts.EmptyMessage))
		}
		return nil
	}

	// Compact title
	if opts.Title != "" {
		fmt.Fprintln(w, styles.title.Render("▸ "+opts.Title))
	}

	// Render based on data type
	switch v := data.(type) {
	case KeyValueData:
		return r.renderCompactKeyValueWithStyles(w, v, styles)
	case TableData:
		return r.renderCompactTableWithStyles(w, v, styles)
	case ListData:
		return r.renderCompactListWithStyles(w, v, styles)
	default:
		return r.ColoredRenderer.RenderWithContext(ctx, w, data, opts)
	}
}

func (r *CompactRenderer) renderCompactKeyValue(w io.Writer, kv KeyValueData) error {
	return r.renderCompactKeyValueWithStyles(w, kv, r.styles)
}

func (r *CompactRenderer) renderCompactKeyValueWithStyles(w io.Writer, kv KeyValueData, styles styles) error {
	for _, pair := range kv.Pairs {
		fmt.Fprintf(w, "%s: %s\n",
			styles.muted.Render(pair.Key),
			pair.Value)
	}
	return nil
}

func (r *CompactRenderer) renderCompactTable(w io.Writer, t TableData) error {
	return r.renderCompactTableWithStyles(w, t, r.styles)
}

func (r *CompactRenderer) renderCompactTableWithStyles(w io.Writer, t TableData, styles styles) error {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, styles.muted.Render("(empty)"))
		return nil
	}

	// Calculate column widths using visible display width (strips ANSI codes).
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = displayWidth(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) {
				if dw := displayWidth(cell); dw > widths[i] {
					widths[i] = dw
				}
			}
		}
	}

	// Apply column constraints if provided.
	if len(t.Constraints) > 0 {
		for i, c := range t.Constraints {
			if i >= len(widths) {
				break
			}
			if c.MaxWidth > 0 && widths[i] > c.MaxWidth {
				widths[i] = c.MaxWidth
			}
			if c.MinWidth > 0 && widths[i] < c.MinWidth {
				widths[i] = c.MinWidth
			}
		}
	}

	// Headers - muted style, no separator
	var headerParts []string
	for i, h := range t.Headers {
		cell := h
		if i < len(t.Constraints) && t.Constraints[i].MaxWidth > 0 && displayWidth(h) > t.Constraints[i].MaxWidth {
			cell = ApplyTruncation(h, t.Constraints[i].MaxWidth, t.Constraints[i].Truncate)
		}
		headerParts = append(headerParts, styles.muted.Render(padToWidth(cell, widths[i])))
	}
	fmt.Fprintln(w, strings.Join(headerParts, " "))

	// Rows - tighter spacing, use padToWidth for ANSI-safe alignment.
	for _, row := range t.Rows {
		var cellParts []string
		for i, cell := range row {
			if i < len(widths) {
				truncated := cell
				if i < len(t.Constraints) && t.Constraints[i].MaxWidth > 0 && displayWidth(cell) > t.Constraints[i].MaxWidth {
					truncated = ApplyTruncation(stripANSI(cell), t.Constraints[i].MaxWidth, t.Constraints[i].Truncate)
				}
				cellParts = append(cellParts, padToWidth(truncated, widths[i]))
			}
		}
		fmt.Fprintln(w, strings.Join(cellParts, " "))
	}

	return nil
}

func (r *CompactRenderer) renderCompactList(w io.Writer, list ListData) error {
	return r.renderCompactListWithStyles(w, list, r.styles)
}

func (r *CompactRenderer) renderCompactListWithStyles(w io.Writer, list ListData, styles styles) error {
	for _, item := range list.Items {
		fmt.Fprintf(w, "  - %s\n", item)
	}
	return nil
}

// RenderMessage delegates to the embedded ColoredRenderer
func (r *CompactRenderer) RenderMessage(w io.Writer, msg Message) error {
	return r.ColoredRenderer.RenderMessage(w, msg)
}

// RenderMessageWithContext delegates to the embedded ColoredRenderer
func (r *CompactRenderer) RenderMessageWithContext(ctx context.Context, w io.Writer, msg Message) error {
	return r.ColoredRenderer.RenderMessageWithContext(ctx, w, msg)
}

// Note: Registration is handled centrally by factory.go init()
