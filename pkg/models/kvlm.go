package models

import (
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// KVLM stands for KV List with Message
type KVLM struct {
	*orderedmap.OrderedMap[string, []string]
}

func (kvlm *KVLM) Append(key string, value string) {
	if existing, ok := kvlm.Get(key); ok {
		kvlm.Set(key, append(existing, value))
	} else {
		kvlm.Set(key, []string{value})
	}
}

func KVLMSerialize(kvlm *KVLM) (string, error) {
	return "", nil
}

func KVLMDeserialize(data string) (*KVLM, error) {
	kvlm := &KVLM{orderedmap.New[string, []string]()}
	lines := strings.Split(data, "\n")

	var key string
	var value strings.Builder

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// If the line is empty, it's the message body.
		if line == "" {
			// TODO: This is a bit hacky and needs to be fixed.
			contains := strings.HasSuffix(value.String(), "--\n")
			if !contains {
				value.WriteString("\n")
				continue
			}
			message := strings.Join(lines[i+1:], "\n")
			kvlm.Set("", []string{message})
			break
		}

		// Detect continuation lines.
		if line[0] == ' ' {
			value.WriteString(line[1:])
			value.WriteString("\n")
			continue
		}

		// Store the previous key-value pair.
		if key != "" {
			kvlm.Append(key, strings.TrimRight(value.String(), "\n"))
			// key = ""
		}

		// Parse the new key-value pair.
		parts := strings.SplitN(line, " ", 2)
		key = parts[0]
		value.Reset()
		if len(parts) > 1 {
			value.WriteString(parts[1])
		}
	}

	// Store the last key-value pair.
	if key != "" {
		kvlm.Append(key, strings.TrimRight(value.String(), "\n"))
	}

	return kvlm, nil
}

func CreateKVLM(orderedmap *orderedmap.OrderedMap[string, []string]) *KVLM {
	return &KVLM{orderedmap}
}
