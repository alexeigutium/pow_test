package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidHash(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		difficulty int
		hash       []byte
		result     bool
	}{
		{
			name:       "valid",
			difficulty: 12,
			hash:       []byte{0, 0b00001000, 0, 0, 0},
			result:     true,
		},
		{
			name:       "valid, many zeros",
			difficulty: 13,
			hash:       []byte{0, 0, 0, 0b00001000, 0, 0, 0},
			result:     true,
		},
		{
			name:       "invalid",
			difficulty: 13,
			hash:       []byte{0, 0b00001000, 0, 0, 0},
			result:     false,
		},
		{
			name:       "too short hash",
			difficulty: 18,
			hash:       []byte{0, 0},
			result:     false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual := IsValidHash(testCase.difficulty, testCase.hash)

			assert.Equal(t, testCase.result, actual)
		})
	}
}

func Test_IsCorrectHash(t *testing.T) {
	correctChallenge := append([]byte("hello, world!"), []byte{3, 0, 0, 0, 0, 0, 0, 0}...)
	correctSolution := []byte{18, 51, 240, 86, 164, 199, 238, 238, 228, 40, 190, 39, 54, 162, 150, 156, 59, 160, 224, 184}
	for _, testCase := range []struct {
		name       string
		difficulty int
		hash       []byte
		data       []byte
		result     bool
	}{
		{
			name:       "correct",
			difficulty: 3,
			hash:       correctSolution,
			data:       correctChallenge,
			result:     true,
		},
		{
			name:       "invalid solution",
			difficulty: 13,
			hash:       correctSolution,
			data:       correctChallenge,
			result:     false,
		},
		{
			name:       "incorrect hash",
			difficulty: 3,
			hash:       correctSolution,
			data:       []byte("hello, world!"),
			result:     false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			actual := IsCorrectHash(testCase.difficulty, testCase.hash, testCase.data)

			assert.Equal(t, testCase.result, actual)
		})
	}
}
