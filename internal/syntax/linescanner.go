package syntax

import (
	"regexp"
	"strings"
)

type LineScanner struct {
	lines   []string
	pos     int
	matched []string
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

func (s *LineScanner) Scan(re *regexp.Regexp) bool {
	if s.EOF() {
		return false
	}
	line := s.lines[s.pos]
	m := re.FindStringSubmatch(line)
	if m != nil {
		s.matched = m
		s.pos++
		return true
	}
	s.matched = nil
	return false
}

func (s *LineScanner) Matched() []string {
	return s.matched
}

func (s *LineScanner) ScanUntil(re *regexp.Regexp) []string {
	var result []string
	for !s.EOF() {
		if s.Scan(re) {
			if len(s.matched) > 0 {
				result = append(result, s.matched[0])
			}
			break
		}
		result = append(result, s.Next())
	}
	return result
}
