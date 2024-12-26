package models

import (
	"fmt"
	"strings"
)

type ShitCommitMetadata struct {
	tree         string
	parent       []string
	author       string
	committer    string
	gpgsignature string
	message      string
}

func CreateShitCommitMetadata(data string) (*ShitCommitMetadata, error) {
	var ismessage bool
	var issignature bool
	var msgbuf strings.Builder
	m := &ShitCommitMetadata{}
	for _, line := range strings.Split(data, "\n") {
		if issignature {
			if len(line) > 0 && line[0] == ' ' {
				line = strings.TrimLeft(line, " ")
				m.gpgsignature += string(line) + "\n"
				continue
			} else if len(line) == 0 {
				m.gpgsignature += "\n"
			} else {
				issignature = false
			}
		}
		if !ismessage {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				ismessage = true
				continue
			}
			split := strings.SplitN(line, " ", 2)
			var data string
			if len(split) == 2 {
				data = split[1]
			}
			switch split[0] {
			case "tree":
				m.tree = string(data)
			case "parent":
				if m.parent == nil {
					m.parent = make([]string, 0)
				}
				m.parent = append(m.parent, string(data))
			case "author":
				m.author = string(data)
			case "committer":
				m.committer = string(data)
			case "gpgsig":
				m.gpgsignature += string(data) + "\n"
				issignature = true
			}
		} else {
			msgbuf.WriteString(line)
		}
	}
	m.message = msgbuf.String()
	return m, nil
}

func (s *ShitCommitMetadata) Serialize() ([]byte, error) {
	var buf strings.Builder
	tree := fmt.Sprintf("tree %s\n", s.tree)
	buf.WriteString(tree)

	var parentstr strings.Builder
	for _, parent := range s.parent {
		parentstr.WriteString(fmt.Sprintf("parent %s\n", parent))
	}
	buf.WriteString(parentstr.String())

	author := fmt.Sprintf("author %s\n", s.author)
	buf.WriteString(author)

	committer := fmt.Sprintf("committer %s\n", s.committer)
	buf.WriteString(committer)

	var gpgsig strings.Builder
	gpgsig.WriteString("gpgsig ")
	// s.gpgsignature = strings.TrimSuffix(s.gpgsignature, "\n")
	signature := s.gpgsignature
	lines := strings.Split(signature, "\n")
	n := len(lines)
	gpgsig.WriteString(fmt.Sprintf("%s\n", lines[0]))
	gpgsig.WriteString(strings.Join(lines[1:n-2], "\n "))
	gpgsig.WriteString(fmt.Sprintf("%s", lines[n-2]))
	buf.WriteString(gpgsig.String())

	buf.WriteString(fmt.Sprintf("\n\n%s", s.message))
	return []byte(buf.String()), nil
}

func (s *ShitCommitMetadata) GetTree() string {
	return s.tree
}

func (s *ShitCommitMetadata) GetParent() []string {
	return s.parent
}

func (s *ShitCommitMetadata) GetAuthor() string {
	return s.author
}

func (s *ShitCommitMetadata) GetCommitter() string {
	return s.committer
}

func (s *ShitCommitMetadata) GetGPGSignature() string {
	return s.gpgsignature
}

func (s *ShitCommitMetadata) GetMessage() string {
	return s.message
}
