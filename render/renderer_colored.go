package render

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rmkohlman/MaestroSDK/colors"
)

// Icons holds the icon set for colored output
type Icons struct {
	Success  string
	Warning  string
	Error    string
	Info     string
	Progress string
	Bullet   string
	Section  string
}

// DefaultIcons returns Unicode icons
func DefaultIcons() Icons {
	return Icons{
		Success:  "✓",
		Warning:  "⚠",
		Error:    "✗",
		Info:     "ℹ",
		Progress: "→",
		Bullet:   "•",
		Section:  "▌",
	}
}

// NerdFontIcons returns Nerd Font icons
func NerdFontIcons() Icons {
	return Icons{
		Success:  "\uf00c", // nf-fa-check
		Warning:  "\uf071", // nf-fa-exclamation_triangle
		Error:    "\uf00d", // nf-fa-times
		Info:     "\uf05a", // nf-fa-info_circle
		Progress: "\uf061", // nf-fa-arrow_right
		Bullet:   "\uf111", // nf-fa-circle
		Section:  "\ue0b0", // powerline
	}
}

// PlainIcons returns ASCII-only icons
func PlainIcons() Icons {
	return Icons{
		Success:  "[OK]",
		Warning:  "[!]",
		Error:    "[X]",
		Info:     "[i]",
		Progress: "->",
		Bullet:   "*",
		Section:  "|",
	}
}

// styles holds lipgloss styles for colored output
type styles struct {
	success   lipgloss.Style
	warning   lipgloss.Style
	errStyle  lipgloss.Style
	info      lipgloss.Style
	muted     lipgloss.Style
	header    lipgloss.Style
	title     lipgloss.Style
	key       lipgloss.Style
	value     lipgloss.Style
	highlight lipgloss.Style
}

func defaultStyles() styles {
	return styles{
		success:   lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")),
		warning:   lipgloss.NewStyle().Foreground(lipgloss.Color("#F9E2AF")),
		errStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8")),
		info:      lipgloss.NewStyle().Foreground(lipgloss.Color("#89B4FA")),
		muted:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7086")),
		header:    lipgloss.NewStyle().Foreground(lipgloss.Color("#CBA6F7")).Bold(true),
		title:     lipgloss.NewStyle().Foreground(lipgloss.Color("#CBA6F7")).Bold(true),
		key:       lipgloss.NewStyle().Foreground(lipgloss.Color("#89DCEB")),
		value:     lipgloss.NewStyle().Foreground(lipgloss.Color("#F5E0DC")),
		highlight: lipgloss.NewStyle().Foreground(lipgloss.Color("#FAB387")).Bold(true),
	}
}

// stylesFromProvider creates styles from a ColorProvider
func stylesFromProvider(provider colors.ColorProvider) styles {
	return styles{
		success:   lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Success())),
		warning:   lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Warning())),
		errStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Error())),
		info:      lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Info())),
		muted:     lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Muted())),
		header:    lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Accent())).Bold(true),
		title:     lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Primary())).Bold(true),
		key:       lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Secondary())),
		value:     lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Foreground())),
		highlight: lipgloss.NewStyle().Foreground(lipgloss.Color(provider.Accent())).Bold(true),
	}
}

// ColoredRenderer outputs richly formatted text with colors and icons.
// This is the default renderer for interactive terminal use.
type ColoredRenderer struct {
	icons  Icons
	styles styles
}

// NewColoredRenderer creates a new colored renderer
func NewColoredRenderer() *ColoredRenderer {
	return &ColoredRenderer{
		icons:  DefaultIcons(),
		styles: defaultStyles(),
	}
}

// NewColoredRendererWithIcons creates a colored renderer with custom icons
func NewColoredRendererWithIcons(icons Icons) *ColoredRenderer {
	return &ColoredRenderer{
		icons:  icons,
		styles: defaultStyles(),
	}
}

// Name returns the renderer identifier
func (r *ColoredRenderer) Name() RendererName {
	return RendererColored
}

// SupportsColor returns true - this renderer uses colors
func (r *ColoredRenderer) SupportsColor() bool {
	return true
}

// getStyles returns styles from ColorProvider if available in context, otherwise defaults
func (r *ColoredRenderer) getStyles(ctx context.Context) styles {
	if provider, ok := colors.FromContext(ctx); ok {
		return stylesFromProvider(provider)
	}
	return r.styles
}

// Render outputs data with colors and formatting
func (r *ColoredRenderer) Render(w io.Writer, data any, opts Options) error {
	// Use background context for backward compatibility
	return r.RenderWithContext(context.Background(), w, data, opts)
}

// RenderWithContext outputs data with colors and formatting using context for theming
func (r *ColoredRenderer) RenderWithContext(ctx context.Context, w io.Writer, data any, opts Options) error {
	styles := r.getStyles(ctx)

	// Handle empty state
	if opts.Empty {
		if opts.EmptyMessage != "" {
			r.RenderMessageWithContext(ctx, w, Message{Level: LevelInfo, Content: opts.EmptyMessage})
		}
		if len(opts.EmptyHints) > 0 {
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Set context with:")
			for _, hint := range opts.EmptyHints {
				fmt.Fprintf(w, "  %s %s\n", styles.info.Render(r.icons.Bullet), hint)
			}
		}
		return nil
	}

	// Render title if present
	if opts.Title != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, styles.title.Render(r.icons.Section+" "+opts.Title))
		fmt.Fprintln(w)
	}

	// Render based on type hint or data type
	switch v := data.(type) {
	case KeyValueData:
		return r.renderKeyValueWithStyles(w, v, styles)
	case TableData:
		return r.renderTableWithStyles(ctx, w, v, styles)
	case ListData:
		return r.renderListWithStyles(w, v, styles)
	case []string:
		return r.renderListWithStyles(w, ListData{Items: v}, styles)
	case map[string]string:
		kv := NewKeyValueData(v)
		return r.renderKeyValueWithStyles(w, kv, styles)
	case map[string]interface{}:
		kv := mapToKeyValueData(v)
		return r.renderKeyValueWithStyles(w, kv, styles)
	default:
		// For other types, just print
		fmt.Fprintf(w, "%v\n", data)
	}

	return nil
}

func (r *ColoredRenderer) renderKeyValueWithStyles(w io.Writer, kv KeyValueData, styles styles) error {
	for _, pair := range kv.Pairs {
		key := styles.key.Render(pair.Key + ":")
		value := styles.value.Render(pair.Value)
		fmt.Fprintf(w, "%s %s\n", key, value)
	}
	return nil
}

func (r *ColoredRenderer) renderTableWithStyles(ctx context.Context, w io.Writer, t TableData, styles styles) error {
	if len(t.Rows) == 0 {
		fmt.Fprintln(w, styles.muted.Render("No data"))
		return nil
	}

	// Determine table color scheme from context or defaults.
	tc := r.tableColors(ctx)
	useColor := tc.enabled

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

	// Border characters
	const (
		vertBorder = "│"
		horzBorder = "─"
		topLeft    = "┌"
		topRight   = "┐"
		topTee     = "┬"
		botLeft    = "└"
		botRight   = "┘"
		botTee     = "┴"
		midLeft    = "├"
		midRight   = "┤"
		midCross   = "┼"
	)

	// Helper: render a border character in muted color
	bc := func(ch string) string {
		if !useColor {
			return ch
		}
		if tc.truecolor {
			return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", tc.borderR, tc.borderG, tc.borderB, ch)
		}
		return fmt.Sprintf("%s%s\x1b[0m", tc.borderCode, ch)
	}

	// Helper: build a horizontal rule (top, middle, or bottom)
	hRule := func(left, mid, right string) string {
		var parts []string
		for i, w := range widths {
			seg := strings.Repeat(horzBorder, w+2) // +2 for cell padding
			if i == 0 {
				parts = append(parts, bc(left)+bc(seg))
			} else {
				parts = append(parts, bc(mid)+bc(seg))
			}
		}
		parts = append(parts, bc(right))
		return strings.Join(parts, "")
	}

	// --- Top border ---
	fmt.Fprintln(w, hRule(topLeft, topTee, topRight))

	// --- Header row ---
	r.renderStyledRow(w, t.Headers, widths, bc, tc, tc.headerBG, tc.headerBGCode, styles.header, useColor)

	// --- Header separator ---
	fmt.Fprintln(w, hRule(midLeft, midCross, midRight))

	// --- Data rows with alternating backgrounds ---
	for rowIdx, row := range t.Rows {
		var bg [3]int
		var bgCode string
		isOddRow := rowIdx%2 == 1
		if useColor && !isOddRow {
			bg = tc.evenRowBG
		} else if useColor {
			bg = tc.oddRowBG
			bgCode = tc.oddRowBGCode
		}
		r.renderStyledRow(w, row, widths, bc, tc, bg, bgCode, styles.value, isOddRow && useColor)
	}

	// --- Bottom border ---
	fmt.Fprintln(w, hRule(botLeft, botTee, botRight))

	return nil
}

// tableColorScheme holds resolved color values for table styling.
// It supports two modes: ANSI (terminal-relative codes) and truecolor (absolute RGB).
type tableColorScheme struct {
	enabled   bool // false when NO_COLOR is set or colors are disabled
	truecolor bool // true = use RGB values; false = use ANSI codes

	// Truecolor fields (used when ColorProvider is available)
	headerBG  [3]int // header row background
	oddRowBG  [3]int // odd data rows (zebra stripe)
	evenRowBG [3]int // even data rows (transparent / no extra bg)
	borderR   int    // border character RGB
	borderG   int
	borderB   int

	// ANSI fields (used for terminal-relative defaults)
	headerBGCode string // e.g., "\x1b[100m" (bright black bg)
	oddRowBGCode string // e.g., "\x1b[100m" (bright black bg)
	borderCode   string // e.g., "\x1b[90m" (bright black fg)
}

// tableColors resolves table color values from ColorProvider in context,
// with sensible terminal-relative defaults when no provider is available.
func (r *ColoredRenderer) tableColors(ctx context.Context) tableColorScheme {
	// Check NO_COLOR env var
	if os.Getenv("NO_COLOR") != "" {
		return tableColorScheme{enabled: false}
	}

	// If a ColorProvider is available, use truecolor (absolute RGB)
	if provider, ok := colors.FromContext(ctx); ok {
		tc := tableColorScheme{
			enabled:   true,
			truecolor: true,
			headerBG:  [3]int{49, 50, 68},
			oddRowBG:  [3]int{30, 30, 46},
			evenRowBG: [3]int{0, 0, 0},
			borderR:   88, borderG: 91, borderB: 112,
		}

		bg := provider.Background()
		border := provider.Border()
		highlight := provider.Highlight()

		if r, g, b, ok := hexToRGB(bg); ok {
			tc.oddRowBG = [3]int{
				clamp(r + 12),
				clamp(g + 12),
				clamp(b + 12),
			}
			tc.evenRowBG = [3]int{0, 0, 0}
		}
		if r, g, b, ok := hexToRGB(highlight); ok {
			tc.headerBG = [3]int{r, g, b}
		}
		if r, g, b, ok := hexToRGB(border); ok {
			tc.borderR = r
			tc.borderG = g
			tc.borderB = b
		}
		return tc
	}

	// Default: use standard ANSI codes (terminal-relative, adapts to color scheme)
	return tableColorScheme{
		enabled:      true,
		truecolor:    false,
		headerBGCode: "\x1b[100m", // bright black background
		oddRowBGCode: "\x1b[100m", // bright black background
		borderCode:   "\x1b[90m",  // bright black foreground (muted)
	}
}

// renderStyledRow writes a single table row with optional background color.
// bgCode is the ANSI code for the background (used in ANSI mode).
// bg is the RGB triplet (used in truecolor mode).
// applyBG controls whether any background is emitted for this row.
func (r *ColoredRenderer) renderStyledRow(
	w io.Writer,
	cells []string,
	widths []int,
	bc func(string) string,
	tc tableColorScheme,
	bg [3]int,
	bgCode string,
	cellStyle lipgloss.Style,
	applyBG bool,
) {
	var buf strings.Builder

	for i, cell := range cells {
		if i >= len(widths) {
			break
		}
		padded := padToWidth(cell, widths[i])

		// Border + space + cell content + space
		buf.WriteString(bc("│"))

		if applyBG {
			if tc.truecolor {
				buf.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", bg[0], bg[1], bg[2]))
			} else {
				buf.WriteString(bgCode)
			}
		}

		buf.WriteString(" ")
		// For header rows the cellStyle is styles.header (bold+colored),
		// for data rows it's styles.value (foreground).
		if i < len(cells) {
			// Apply lipgloss style only to the text, not padding
			styled := cellStyle.Render(padded)
			buf.WriteString(styled)
		}
		buf.WriteString(" ")

		if applyBG {
			buf.WriteString("\x1b[0m")
		}
	}

	// Closing border
	buf.WriteString(bc("│"))

	fmt.Fprintln(w, buf.String())
}

func (r *ColoredRenderer) renderListWithStyles(w io.Writer, list ListData, styles styles) error {
	for _, item := range list.Items {
		fmt.Fprintf(w, "  %s %s\n", styles.info.Render(r.icons.Bullet), item)
	}
	return nil
}

// RenderMessage outputs a styled message
func (r *ColoredRenderer) RenderMessage(w io.Writer, msg Message) error {
	// Use background context for backward compatibility
	return r.RenderMessageWithContext(context.Background(), w, msg)
}

// RenderMessageWithContext outputs a styled message with context for theming
func (r *ColoredRenderer) RenderMessageWithContext(ctx context.Context, w io.Writer, msg Message) error {
	styles := r.getStyles(ctx)

	var icon string
	var style lipgloss.Style

	switch msg.Level {
	case LevelSuccess:
		icon = r.icons.Success
		style = styles.success
	case LevelWarning:
		icon = r.icons.Warning
		style = styles.warning
	case LevelError:
		icon = r.icons.Error
		style = styles.errStyle
	case LevelProgress:
		icon = r.icons.Progress
		style = styles.info
	case LevelDebug:
		icon = ""
		style = styles.muted
	case LevelInfo:
		fallthrough
	default:
		icon = r.icons.Info
		style = styles.info
	}

	if icon != "" {
		fmt.Fprintln(w, style.Render(icon+" "+msg.Content))
	} else {
		fmt.Fprintln(w, style.Render(msg.Content))
	}
	return nil
}

// Note: Registration is handled centrally by factory.go init()
