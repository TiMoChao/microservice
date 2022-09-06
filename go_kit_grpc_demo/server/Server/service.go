package Server

import (
	"context"
	"errors"
	gl "github.com/drhodes/golorem"
	"strings"
)

var (
	ErrRequestTypeNotFound = errors.New("Request type only valid for word, sentence and paragraph")
)

// Define service interface
type Service interface {
	// generate a word with at least min letters and at most max letters.
	Timor(ctx context.Context, requestType string, min, max int) (string, error)
}

// Implement service with empty struct
type TimorService struct {
}

// Implement service functions
func (TimorService) Timor(_ context.Context, requestType string, min, max int) (string, error) {
	var result string
	var err error
	if strings.EqualFold(requestType, "Word") {
		result = gl.Word(min, max)
	} else if strings.EqualFold(requestType, "Sentence") {
		result = gl.Sentence(min, max)
	} else if strings.EqualFold(requestType, "Paragraph") {
		result = gl.Paragraph(min, max)
	} else {
		err = ErrRequestTypeNotFound
	}
	return result, err
}
