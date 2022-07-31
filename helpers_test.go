package moccasin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFnName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "FullPath",
			input:  "github.com/mypackage.(*MyStruct).myFunc",
			output: "myFunc",
		},
		{
			name:   "ShortPath",
			input:  "(*MyStruct).myFunc",
			output: "myFunc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.output, parseFnName(tt.input))
		})
	}
}
