package editor

import (
	"bufio"
	"context"
	"strings"
	"unicode/utf8"

	"github.com/fakovacic/editor/internal/app"
)

func (s *editor) Change(_ context.Context, msg *app.ChangeMsg) error {
	s.file.Lock()
	defer s.file.Unlock()

	switch msg.Action {
	case OpInsert:
		s.file.Contents = s.Insert(s.file.Contents, msg)
	case OpRemove:
		s.file.Contents = s.Remove(s.file.Contents, msg)
	}

	return nil
}

func (s *editor) Insert(content string, msg *app.ChangeMsg) string {
	contentScanner := bufio.NewScanner(strings.NewReader(content))
	contentScanner.Split(bufio.ScanLines)

	var (
		newValue strings.Builder
		row      int
	)

	totalRows := linesStringCount(s.file.Contents)

	lineSuffix := "\n"

	for contentScanner.Scan() {
		// last line should not have new line
		if row == totalRows {
			lineSuffix = ""
		}

		line := contentScanner.Text()

		switch {
		case row == msg.Start.Row:
			textBefore := line[0:msg.Start.Column]
			textAfter := line[msg.Start.Column:]

			// insert content between text
			line = textBefore + strings.Join(msg.Lines, "\n") + textAfter

			newValue.WriteString(line + lineSuffix)

			if len(msg.Lines) > 1 {
				row += (len(msg.Lines) - 1)

				continue
			}

			row++
		default:
			newValue.WriteString(line + lineSuffix)

			row++
		}
	}

	// write additional lines
	if msg.End.Row > totalRows {
		newValue.WriteString(strings.Join(msg.Lines, "\n"))
	}

	return newValue.String()
}

func (s *editor) Remove(content string, msg *app.ChangeMsg) string {
	contentScanner := bufio.NewScanner(strings.NewReader(content))
	contentScanner.Split(bufio.ScanLines)

	var (
		newValue strings.Builder
		row      int
		linesRow int
	)

	totalRows := linesStringCount(s.file.Contents)

	lineSuffix := "\n"

	for contentScanner.Scan() {
		// last line should not have new line
		if row == totalRows {
			lineSuffix = ""
		}

		line := contentScanner.Text()

		switch {
		case row < msg.Start.Row:
			newValue.WriteString(line + lineSuffix)

			row++
		case row > msg.End.Row:
			newValue.WriteString(line + lineSuffix)

			row++
		case msg.End.Row == msg.Start.Row:
			contentCount := utf8.RuneCountInString(msg.Lines[linesRow])

			textBefore := line[0:msg.Start.Column]
			textAfter := line[(msg.Start.Column + contentCount):]

			line = textBefore + textAfter

			newValue.WriteString(line + lineSuffix)

			row++
			linesRow++
		case row == msg.Start.Row:
			newValue.WriteString(line[0:msg.Start.Column])

			row++
			linesRow++
		case row == msg.End.Row:
			pointer := msg.End.Column
			if msg.End.Column != 0 {
				pointer = msg.End.Column - 1
			}

			line = line[pointer:]

			newValue.WriteString(line + lineSuffix)

			row++
			linesRow++
		default:
			row++
			linesRow++
		}
	}

	return newValue.String()
}

func linesStringCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}

	return n
}
