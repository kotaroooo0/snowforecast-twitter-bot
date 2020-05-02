package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseStringToSnowResorts(t *testing.T) {
	cases := []struct {
		input  string
		output []SnowResort
	}{
		{
			input: "[[\"\",[[\"HiddenValley2\",\"Hidden Valley Ski\"],[\"Snow-Creek\",\"Snow Creek\"]]]]",
			output: []SnowResort{
				SnowResort{"HiddenValley2", "Hidden Valley Ski"},
				SnowResort{"Snow-Creek", "Snow Creek"},
			},
		},
	}

	for _, tt := range cases {
		if diff := cmp.Diff(parseStringToSnowResorts(tt.input), tt.output); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}
