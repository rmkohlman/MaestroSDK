package render

// Tests for Issue #207: Lipgloss Table Rendering (TDD Phase 2 — RED)
//
// These tests drive the creation of:
//   - PrettyRenderer: a new renderer registered under RendererPretty ("pretty")
//   - TableStyleProvider: interface in SDK that DVM implements for style injection
//   - "colored" as a deprecated alias that resolves to the same renderer as "table"
//
// All tests in this file are expected to FAIL until the implementation is added.

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/rmkohlman/MaestroSDK/colors"
)

// ---------------------------------------------------------------------------
// 1. PrettyRenderer existence and registration
// ---------------------------------------------------------------------------

func TestPrettyRenderer_IsRegistered(t *testing.T) {
	r := Get(RendererPretty)
	if r == nil {
		t.Fatal("RendererPretty ('pretty') is not registered in the global registry")
	}
}

func TestPrettyRenderer_ImplementsRendererInterface(t *testing.T) {
	// Compile-time check: PrettyRenderer must implement Renderer
	var _ Renderer = &PrettyRenderer{}
}

func TestPrettyRenderer_Name(t *testing.T) {
	r := NewPrettyRenderer()
	if r == nil {
		t.Fatal("NewPrettyRenderer() returned nil")
	}
	if r.Name() != RendererPretty {
		t.Errorf("PrettyRenderer.Name() = %q, want %q", r.Name(), RendererPretty)
	}
}

func TestPrettyRenderer_SupportsColor(t *testing.T) {
	r := NewPrettyRenderer()
	if !r.SupportsColor() {
		t.Error("PrettyRenderer.SupportsColor() should return true")
	}
}

// ---------------------------------------------------------------------------
// 2. RendererName constant for "pretty"
// ---------------------------------------------------------------------------

func TestRendererPretty_ConstantExists(t *testing.T) {
	// RendererPretty constant must exist with value "pretty"
	if string(RendererPretty) != "pretty" {
		t.Errorf("RendererPretty = %q, want %q", RendererPretty, "pretty")
	}
}

// ---------------------------------------------------------------------------
// 3. PrettyRenderer.RenderTable produces bordered output (box-drawing characters)
// ---------------------------------------------------------------------------

func TestPrettyRenderer_RenderTable_ProducesBorderedOutput(t *testing.T) {
	r := NewPrettyRenderer()
	var buf bytes.Buffer

	data := TableData{
		Headers: []string{"NAME", "STATUS", "APP"},
		Rows: [][]string{
			{"main", "running", "portal"},
			{"dev", "stopped", "api"},
		},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()

	// Must contain box-drawing characters (bordered table)
	boxChars := []string{"│", "─", "┌", "┐", "└", "┘", "├", "┤", "┬", "┴", "┼"}
	foundAny := false
	for _, ch := range boxChars {
		if strings.Contains(output, ch) {
			foundAny = true
			break
		}
	}
	if !foundAny {
		t.Errorf("PrettyRenderer output contains no box-drawing characters.\nOutput:\n%s", output)
	}
}

func TestPrettyRenderer_RenderTable_ContainsHeaders(t *testing.T) {
	r := NewPrettyRenderer()
	var buf bytes.Buffer

	data := TableData{
		Headers: []string{"NAME", "STATUS"},
		Rows: [][]string{
			{"workspace-1", "active"},
		},
	}

	err := r.Render(&buf, data, Options{Type: TypeTable})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "NAME") {
		t.Errorf("Expected output to contain header 'NAME', got:\n%s", output)
	}
	if !strings.Contains(output, "workspace-1") {
		t.Errorf("Expected output to contain row data 'workspace-1', got:\n%s", output)
	}
}

func TestPrettyRenderer_RenderTable_WithContextUsingStyleProvider(t *testing.T) {
	r := NewPrettyRenderer()
	var buf bytes.Buffer

	data := TableData{
		Headers: []string{"ECOSYSTEM", "DOMAINS"},
		Rows: [][]string{
			{"prod", "3"},
			{"staging", "2"},
		},
	}

	// Inject a TableStyleProvider via context
	provider := &mockTableStyleProvider{}
	ctx := WithTableStyleProvider(context.Background(), provider)

	err := r.RenderWithContext(ctx, &buf, data, Options{Type: TypeTable})
	if err != nil {
		t.Fatalf("RenderWithContext failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("Expected non-empty output with style provider")
	}
	if !strings.Contains(output, "ECOSYSTEM") {
		t.Errorf("Expected output to contain 'ECOSYSTEM', got:\n%s", output)
	}
}

func TestPrettyRenderer_RenderTable_FallbackWithoutStyleProvider(t *testing.T) {
	r := NewPrettyRenderer()
	var buf bytes.Buffer

	data := TableData{
		Headers: []string{"NAME"},
		Rows:    [][]string{{"test-ws"}},
	}

	// No style provider — should fall back gracefully (no panic, valid output)
	err := r.RenderWithContext(context.Background(), &buf, data, Options{Type: TypeTable})
	if err != nil {
		t.Fatalf("RenderWithContext without provider failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected non-empty output even without style provider")
	}
}

// ---------------------------------------------------------------------------
// 4. TableStyleProvider interface exists and is usable
// ---------------------------------------------------------------------------

// mockTableStyleProvider implements TableStyleProvider for testing.
// The method signatures here drive the interface definition.
type mockTableStyleProvider struct{}

// Compile-time assertion: mockTableStyleProvider must implement TableStyleProvider
var _ TableStyleProvider = &mockTableStyleProvider{}

func (m *mockTableStyleProvider) HeaderStyle() TableCellStyle {
	return TableCellStyle{FG: "#CBA6F7", Bold: true}
}
func (m *mockTableStyleProvider) CellStyle() TableCellStyle   { return TableCellStyle{FG: "#F5E0DC"} }
func (m *mockTableStyleProvider) BorderStyle() TableCellStyle { return TableCellStyle{FG: "#585B70"} }

func TestTableStyleProvider_InterfaceExists(t *testing.T) {
	// This test passes if the above var _ declaration compiles.
	// If TableStyleProvider doesn't exist, the file won't compile.
	var p TableStyleProvider = &mockTableStyleProvider{}
	if p == nil {
		t.Error("TableStyleProvider should not be nil")
	}
}

func TestWithTableStyleProvider_StoresInContext(t *testing.T) {
	provider := &mockTableStyleProvider{}
	ctx := WithTableStyleProvider(context.Background(), provider)
	retrieved, ok := TableStyleProviderFromContext(ctx)
	if !ok {
		t.Fatal("TableStyleProviderFromContext should find the provider stored by WithTableStyleProvider")
	}
	if retrieved == nil {
		t.Error("Retrieved TableStyleProvider should not be nil")
	}
}

func TestTableStyleProviderFromContext_MissingProvider(t *testing.T) {
	_, ok := TableStyleProviderFromContext(context.Background())
	if ok {
		t.Error("TableStyleProviderFromContext should return false when no provider in context")
	}
}

// ---------------------------------------------------------------------------
// 5. "colored" as deprecated alias for the table renderer
// ---------------------------------------------------------------------------

func TestResolveRenderer_ColoredMapsToTableRenderer(t *testing.T) {
	// "colored" is a deprecated alias — resolving it should give the same
	// renderer type as resolving "table"
	coloredR := ResolveRenderer("colored")
	if coloredR == nil {
		t.Fatal("ResolveRenderer('colored') returned nil")
	}

	// The resolved renderer for "colored" must be the table renderer (alias)
	if coloredR.Name() != RendererTable {
		t.Errorf("ResolveRenderer('colored').Name() = %q, want %q (deprecated alias)", coloredR.Name(), RendererTable)
	}
}

func TestResolveRenderer_TableAndColoredReturnSameType(t *testing.T) {
	tableR := ResolveRenderer("table")
	coloredR := ResolveRenderer("colored")

	if tableR == nil {
		t.Fatal("ResolveRenderer('table') returned nil")
	}
	if coloredR == nil {
		t.Fatal("ResolveRenderer('colored') returned nil")
	}

	// Both should resolve to the same renderer name
	if tableR.Name() != coloredR.Name() {
		t.Errorf("'table' resolves to %q but 'colored' resolves to %q — they should be the same",
			tableR.Name(), coloredR.Name())
	}
}

// ---------------------------------------------------------------------------
// 6. OutputWithContext — style provider flows through to renderer
// ---------------------------------------------------------------------------

func TestOutputWithContext_WithTableStyleProvider_ProducesOutput(t *testing.T) {
	// Save and restore writer
	orig := GetWriter()
	defer SetWriter(orig)
	var buf bytes.Buffer
	SetWriter(&buf)

	provider := &mockTableStyleProvider{}
	ctx := WithTableStyleProvider(context.Background(), provider)
	ctx = colors.WithProvider(ctx, &MockColorProvider{})

	data := TableData{
		Headers: []string{"NAME", "STATUS"},
		Rows:    [][]string{{"ws-1", "active"}},
	}

	err := OutputWithContext(ctx, data, Options{Type: TypeTable})
	if err != nil {
		t.Fatalf("OutputWithContext failed: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("Expected non-empty output from OutputWithContext with TableStyleProvider")
	}
}

func TestOutputWithContext_WithoutTableStyleProvider_FallsBackGracefully(t *testing.T) {
	orig := GetWriter()
	defer SetWriter(orig)
	var buf bytes.Buffer
	SetWriter(&buf)

	// No TableStyleProvider in context — must not panic or error
	data := TableData{
		Headers: []string{"NAME"},
		Rows:    [][]string{{"ws-1"}},
	}

	err := OutputWithContext(context.Background(), data, Options{Type: TypeTable})
	if err != nil {
		t.Fatalf("OutputWithContext without TableStyleProvider failed: %v", err)
	}
}
