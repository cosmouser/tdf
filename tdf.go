package tdf

import (
	"bufio"
	"io"
	"strings"
)

// Entry is an entry in a tdf file
type Entry struct {
	Name string
	Info map[string]interface{}
}
type tdfNode struct {
	value *Entry
	next  *tdfNode
}
type tdfStack struct {
	top  *tdfNode
	size int
}

func (s *tdfStack) pop() *tdfNode {
	if s.size == 0 {
		return nil
	}
	tmp := s.top
	s.top = s.top.next
	s.size--
	return tmp
}
func (s *tdfStack) push(n *tdfNode) {
	if s.size == 0 {
		s.top = n
		s.size++
		return
	}
	n.next = s.top
	s.top = n
	s.size++
}

// Scan reads tdf data and returns a slice of TdfMap structs. The keys
// map is used so that attributes are recorded in a consistent casing
func Scan(obj io.Reader, keys map[string]string) []*Entry {
	scanner := bufio.NewScanner(obj)
	var comment, inEntry bool
	entries := []*Entry{}
	stack := &tdfStack{}
	for scanner.Scan() {
		line := scanner.Text()
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
					info := make(map[string]interface{})
					stack.push(&tdfNode{value: &Entry{Name: entryName, Info: info}})
					inEntry = false
					continue
				}
				if line[i] == '{' {
					inEntry = true
				}
				if line[i] == '}' && stack.size > 0 {
					if stack.size == 1 {
						entries = append(entries, stack.pop().value)
						inEntry = false
					} else {
						tmpEntry := stack.pop().value
						stack.top.value.Info[tmpNode.value.Name] = tmpEntry
					}
				}
				if line[i] == '=' && stack.size > 0 && inEntry {
					attribute := strings.Trim(line[:i], " \t")
					value := strings.Split(line[i+1:], ";")
					if key, ok := keys[strings.ToLower(attribute)]; !ok {
						attrLower := strings.ToLower(attribute)
						keys[attrLower] = attribute
						stack.top.value.Info[attribute] = value[0]
						break
					} else {
						stack.top.value.Info[key] = value[0]
						break
					}
				}
			}
		}
	}
	return entries
}
