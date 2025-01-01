package models

import (
	"fmt"
	"strings"
)

type ShitTagMetadata struct {
	object     string
	objecttype string
	tag        string
	tagger     string
	message    string
}

func (s *ShitTagMetadata) Serialize() (string, error) {
	var serialized strings.Builder
	serialized.WriteString(fmt.Sprintf("object %s\n", s.object))
	serialized.WriteString(fmt.Sprintf("type %s\n", s.objecttype))
	serialized.WriteString(fmt.Sprintf("tag %s\n", s.tag))
	serialized.WriteString(fmt.Sprintf("tagger %s\n", s.tagger))
	serialized.WriteString("\n")
	serialized.WriteString(s.message)
	return serialized.String(), nil
}

func CreateShitTagMetadata(data string) (*ShitTagMetadata, error) {
	m := &ShitTagMetadata{}
	var ismessage bool
	var msgbuf strings.Builder
	for _, line := range strings.Split(data, "\n") {
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
			case "object":
				m.object = string(data)
			case "type":
				m.objecttype = string(data)
			case "tag":
				m.tag = string(data)
			case "tagger":
				m.tagger = string(data)
			}
		} else {
			msgbuf.WriteString(line)
		}
	}
	m.message = msgbuf.String()
	return m, nil
}

func CreateShitTagMetadataFromAttr(
	object string,
	objecttype string,
	tag string,
	tagger string,
	message string) *ShitTagMetadata {
	return &ShitTagMetadata{
		object:     object,
		objecttype: objecttype,
		tag:        tag,
		tagger:     tagger,
		message:    message,
	}
}
