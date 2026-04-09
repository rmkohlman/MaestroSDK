package render

import "context"

// TableCellStyle defines the visual style for a table cell.
// This is a simple value type — SDK defines the shape, DVM provides values.
type TableCellStyle struct {
	// FG is the foreground color as a hex string (e.g. "#CBA6F7")
	FG string

	// Bold indicates whether the text should be bold
	Bold bool
}

// TableStyleProvider is the interface that DVM implements to inject
// theme-aware styles into the PrettyRenderer. SDK defines this interface;
// DVM provides the concrete implementation via context injection.
//
// This keeps the dependency arrow correct: DVM → SDK (never SDK → DVM).
type TableStyleProvider interface {
	// HeaderStyle returns the style for table header cells
	HeaderStyle() TableCellStyle

	// CellStyle returns the style for regular table data cells
	CellStyle() TableCellStyle

	// BorderStyle returns the style for table border characters
	BorderStyle() TableCellStyle
}

// tableStyleProviderKey is the context key for storing a TableStyleProvider.
type tableStyleProviderKey struct{}

// WithTableStyleProvider stores a TableStyleProvider in the context.
func WithTableStyleProvider(ctx context.Context, provider TableStyleProvider) context.Context {
	return context.WithValue(ctx, tableStyleProviderKey{}, provider)
}

// TableStyleProviderFromContext retrieves the TableStyleProvider from context.
// Returns the provider and true if found, or nil and false if not present.
func TableStyleProviderFromContext(ctx context.Context) (TableStyleProvider, bool) {
	provider, ok := ctx.Value(tableStyleProviderKey{}).(TableStyleProvider)
	return provider, ok
}
