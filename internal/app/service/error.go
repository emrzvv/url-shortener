package service

import "fmt"

type InvalidURLError struct {
	value string
}

func (e *InvalidURLError) Error() string {
	return fmt.Sprintf("invalid URL: %s", e.value)
}

type InvalidIDError struct {
	value string
}

func (e *InvalidIDError) Error() string {
	return fmt.Sprintf("invalid ID: %s", e.value)
}
