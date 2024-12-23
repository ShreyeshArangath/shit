package models

import "fmt"

type ShitException struct {
	Message string
}

func (e *ShitException) Error() string {
	return fmt.Sprintf("ShitException: %s", e.Message)
}
