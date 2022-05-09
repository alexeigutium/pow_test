package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Quotes interface {
	Get() string
}

func NewFileQuotes(filename string) (Quotes, error) {
	rand.Seed(time.Now().Unix())
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't read quotes file: %w", err)
	}

	quotes := strings.Split(string(content), "\n")
	if len(quotes) < 5 {
		return nil, errors.New("too few quotes in the file, minimum is 5")
	}

	return &fileQuotes{
		quotes: quotes,
	}, nil
}

type fileQuotes struct {
	quotes []string
}

func (q fileQuotes) Get() string {
	return q.quotes[rand.Intn(len(q.quotes))]
}
