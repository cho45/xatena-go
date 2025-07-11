package syntax

import (
	"regexp"
	"testing"
)

func TestNewLineScanner(t *testing.T) {
	input := "line1\nline2\nline3"
	scanner := NewLineScanner(input)

	if scanner.pos != 0 {
		t.Errorf("NewLineScanner pos = %d, want 0", scanner.pos)
	}

	expectedLines := []string{"line1", "line2", "line3"}
	if len(scanner.lines) != len(expectedLines) {
		t.Errorf("NewLineScanner lines count = %d, want %d", len(scanner.lines), len(expectedLines))
	}

	for i, expected := range expectedLines {
		if scanner.lines[i] != expected {
			t.Errorf("NewLineScanner lines[%d] = %q, want %q", i, scanner.lines[i], expected)
		}
	}
}

func TestLineScanner_Peek(t *testing.T) {
	input := "line1\nline2\nline3"
	scanner := NewLineScanner(input)

	// First peek
	line := scanner.Peek()
	if line != "line1" {
		t.Errorf("Peek() = %q, want %q", line, "line1")
	}

	// Position should not change
	if scanner.pos != 0 {
		t.Errorf("Peek() changed position to %d, want 0", scanner.pos)
	}

	// Second peek should return same line
	line2 := scanner.Peek()
	if line2 != "line1" {
		t.Errorf("Second Peek() = %q, want %q", line2, "line1")
	}
}

func TestLineScanner_Next(t *testing.T) {
	input := "line1\nline2\nline3"
	scanner := NewLineScanner(input)

	expected := []string{"line1", "line2", "line3"}
	for i, expectedLine := range expected {
		line := scanner.Next()
		if line != expectedLine {
			t.Errorf("Next() [%d] = %q, want %q", i, line, expectedLine)
		}
		if scanner.pos != i+1 {
			t.Errorf("Next() [%d] position = %d, want %d", i, scanner.pos, i+1)
		}
	}

	// After consuming all lines
	if !scanner.EOF() {
		t.Error("EOF() should return true after consuming all lines")
	}

	// Next on EOF should return empty string
	line := scanner.Next()
	if line != "" {
		t.Errorf("Next() on EOF = %q, want empty string", line)
	}
}

func TestLineScanner_EOF(t *testing.T) {
	input := "line1\nline2"
	scanner := NewLineScanner(input)

	// Not EOF initially
	if scanner.EOF() {
		t.Error("EOF() should return false initially")
	}

	// Consume first line
	scanner.Next()
	if scanner.EOF() {
		t.Error("EOF() should return false after consuming first line")
	}

	// Consume second line
	scanner.Next()
	if !scanner.EOF() {
		t.Error("EOF() should return true after consuming all lines")
	}
}

func TestLineScanner_Reset(t *testing.T) {
	input := "line1\nline2\nline3"
	scanner := NewLineScanner(input)

	// Consume some lines
	scanner.Next()
	scanner.Next()

	// Reset to beginning
	scanner.Reset(0)
	if scanner.pos != 0 {
		t.Errorf("Reset(0) position = %d, want 0", scanner.pos)
	}

	// Should be able to read from beginning
	line := scanner.Peek()
	if line != "line1" {
		t.Errorf("Peek() after Reset(0) = %q, want %q", line, "line1")
	}

	// Reset to middle
	scanner.Reset(1)
	if scanner.pos != 1 {
		t.Errorf("Reset(1) position = %d, want 1", scanner.pos)
	}

	line = scanner.Peek()
	if line != "line2" {
		t.Errorf("Peek() after Reset(1) = %q, want %q", line, "line2")
	}
}

func TestLineScanner_Pos(t *testing.T) {
	input := "line1\nline2\nline3"
	scanner := NewLineScanner(input)

	if scanner.Pos() != 0 {
		t.Errorf("Initial Pos() = %d, want 0", scanner.Pos())
	}

	scanner.Next()
	if scanner.Pos() != 1 {
		t.Errorf("Pos() after Next() = %d, want 1", scanner.Pos())
	}

	scanner.Next()
	if scanner.Pos() != 2 {
		t.Errorf("Pos() after second Next() = %d, want 2", scanner.Pos())
	}
}

func TestLineScanner_SetLines(t *testing.T) {
	scanner := NewLineScanner("original")

	newLines := []string{"new1", "new2", "new3"}
	scanner.SetLines(newLines)

	if scanner.pos != 0 {
		t.Errorf("SetLines() position = %d, want 0", scanner.pos)
	}

	if len(scanner.lines) != len(newLines) {
		t.Errorf("SetLines() lines count = %d, want %d", len(scanner.lines), len(newLines))
	}

	for i, expected := range newLines {
		if scanner.lines[i] != expected {
			t.Errorf("SetLines() lines[%d] = %q, want %q", i, scanner.lines[i], expected)
		}
	}
}

func TestLineScanner_Scan(t *testing.T) {
	input := "* item1\n+ item2\nregular text"
	scanner := NewLineScanner(input)

	// Pattern for list items
	listPattern := regexp.MustCompile(`^([*+])\s+(.+)$`)

	// Should match first line
	if !scanner.Scan(listPattern) {
		t.Error("Scan() should match first line")
	}

	matched := scanner.Matched()
	if len(matched) != 3 {
		t.Errorf("Matched() length = %d, want 3", len(matched))
	}
	if matched[0] != "* item1" {
		t.Errorf("Matched()[0] = %q, want %q", matched[0], "* item1")
	}
	if matched[1] != "*" {
		t.Errorf("Matched()[1] = %q, want %q", matched[1], "*")
	}
	if matched[2] != "item1" {
		t.Errorf("Matched()[2] = %q, want %q", matched[2], "item1")
	}

	// Should match second line
	if !scanner.Scan(listPattern) {
		t.Error("Scan() should match second line")
	}

	matched = scanner.Matched()
	if matched[1] != "+" {
		t.Errorf("Matched()[1] = %q, want %q", matched[1], "+")
	}
	if matched[2] != "item2" {
		t.Errorf("Matched()[2] = %q, want %q", matched[2], "item2")
	}

	// Should not match third line
	if scanner.Scan(listPattern) {
		t.Error("Scan() should not match third line")
	}

	if scanner.Matched() != nil {
		t.Error("Matched() should be nil after failed scan")
	}
}

func TestLineScanner_ScanUntil(t *testing.T) {
	input := "line1\nline2\n>>>\nline4\nline5"
	scanner := NewLineScanner(input)

	// Pattern for end marker
	endPattern := regexp.MustCompile(`^>>>$`)

	result := scanner.ScanUntil(endPattern)

	expected := []string{"line1", "line2", ">>>"}
	if len(result) != len(expected) {
		t.Errorf("ScanUntil() length = %d, want %d", len(result), len(expected))
	}

	for i, expectedLine := range expected {
		if result[i] != expectedLine {
			t.Errorf("ScanUntil() result[%d] = %q, want %q", i, result[i], expectedLine)
		}
	}

	// Scanner should be positioned after the matched line
	if scanner.Peek() != "line4" {
		t.Errorf("Peek() after ScanUntil() = %q, want %q", scanner.Peek(), "line4")
	}
}

func TestLineScanner_ScanUntil_NoMatch(t *testing.T) {
	input := "line1\nline2\nline3"
	scanner := NewLineScanner(input)

	// Pattern that won't match
	endPattern := regexp.MustCompile(`^END$`)

	result := scanner.ScanUntil(endPattern)

	expected := []string{"line1", "line2", "line3"}
	if len(result) != len(expected) {
		t.Errorf("ScanUntil() no match length = %d, want %d", len(result), len(expected))
	}

	for i, expectedLine := range expected {
		if result[i] != expectedLine {
			t.Errorf("ScanUntil() no match result[%d] = %q, want %q", i, result[i], expectedLine)
		}
	}

	// Scanner should be at EOF
	if !scanner.EOF() {
		t.Error("Scanner should be at EOF after ScanUntil with no match")
	}
}

func TestLineScanner_EmptyInput(t *testing.T) {
	scanner := NewLineScanner("")

	// Empty input creates a slice with one empty string
	if scanner.EOF() {
		t.Error("EOF() should return false for empty input initially")
	}

	if scanner.Peek() != "" {
		t.Errorf("Peek() on empty input = %q, want empty string", scanner.Peek())
	}

	if scanner.Next() != "" {
		t.Errorf("Next() on empty input = %q, want empty string", scanner.Next())
	}
}

func TestLineScanner_SingleLine(t *testing.T) {
	scanner := NewLineScanner("single line")

	if scanner.EOF() {
		t.Error("EOF() should return false for single line input")
	}

	line := scanner.Next()
	if line != "single line" {
		t.Errorf("Next() = %q, want %q", line, "single line")
	}

	if !scanner.EOF() {
		t.Error("EOF() should return true after consuming single line")
	}
}
