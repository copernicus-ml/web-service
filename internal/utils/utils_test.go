package utils_test

import (
	"testing"

	"github.com/copernicus-ml/web-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseSymbols(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"valid sequence", "a, b, c", []string{"A", "B", "C"}},
		{"partial empty sequence", "a, ,, , c, dEf", []string{"A", "C", "DEF"}},
		{"empty sequence", ", ,, , , ", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedSymbols := utils.ParseSymbols(tt.input)
			assert.Equal(t, parsedSymbols, tt.expected)
		})
	}
}
