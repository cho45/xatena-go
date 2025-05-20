package syntax

import "strings"

type LineScanner struct {
	lines []string
	pos   int
}

func NewLineScanner(input string) *LineScanner {
	lines := strings.Split(input, "\n")
	return &LineScanner{lines: lines, pos: 0}
}

func (s *LineScanner) Next() string {
	if s.pos >= len(s.lines) {
		return ""
	}
	line := s.lines[s.pos]
	s.pos++
	return line
}

func (s *LineScanner) Peek() string {
	if s.pos >= len(s.lines) {
		return ""
	}
	return s.lines[s.pos]
}

func (s *LineScanner) EOF() bool {
	return s.pos >= len(s.lines)
}

func (s *LineScanner) Reset(pos int) {
	s.pos = pos
}

func (s *LineScanner) Pos() int {
	return s.pos
}

func (s *LineScanner) SetLines(lines []string) {
	s.lines = lines
	s.pos = 0
}
