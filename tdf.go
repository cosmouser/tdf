package tdf

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Node is an entry in a tdf/ota/fbi file
type Node struct {
	Name     string
	Fields   map[string]string
	Children []*Node
}

// /* ignore these comments */
// and these: // until EOL

// Decode parses an input stream and returns a slice of Nodes
func Decode(reader io.Reader) ([]*Node, error) {
	result := []*Node{}
	stack := []*Node{}
	var (
		comment bool
		inEntry bool
	)
	parser := func(line string) error {
		for i := 0; i < len(line); i++ {
			if !comment && i+1 < len(line) {
				if line[i] == '/' && line[i+1] == '/' {
					i = len(line)
					continue
				}
				if line[i] == '/' && line[i+1] == '*' {
					comment = true
					continue
				}
			}
			if comment && i != 0 {
				if line[i] == '/' && line[i-1] == '*' {
					comment = false
					continue
				}
			}
			if !comment {
				if line[i] == '[' {
					var entryName string
					i++
					for line[i] != ']' {
						entryName += string(line[i])
						i++
					}
					stack = append(stack, &Node{
						Name:     entryName,
						Fields:   make(map[string]string),
						Children: make([]*Node, 0),
					})
					inEntry = false
					continue
				}
				if line[i] == '{' {
					inEntry = true
				}
				if line[i] == '}' && len(stack) > 0 {
					if len(stack) == 1 {
						result = append(result, stack[0])
						stack = stack[:len(stack)-1]
						inEntry = false
					} else {
						tmpEntry := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, tmpEntry)
					}
				}
				if line[i] == '=' && len(stack) > 0 && inEntry {
					attribute := strings.Trim(line[:i], " \t")
					value := strings.Split(line[i+1:], ";")
					if len(value) < 1 {
						return errors.New("found = with no ;")
					}
					stack[len(stack)-1].Fields[strings.ToLower(attribute)] = value[0]
				}
			}
		}
		return nil
	}
	lineScanner := bufio.NewScanner(reader)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		line = strings.TrimSpace(line)
		err := parser(line)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
