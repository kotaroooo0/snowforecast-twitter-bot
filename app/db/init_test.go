package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kotaroooo0/snowforecast-twitter-bot/repository"
)

func TestParseStringToSnowResorts(t *testing.T) {
	cases := []struct {
		input  string
		output []repository.SnowResort
	}{
		{
			input: "[[\"\",[[\"HiddenValley2\",\"Hidden Valley Ski\"],[\"Snow-Creek\",\"Snow Creek\"]]]]",
			output: []repository.SnowResort{
				repository.SnowResort{"HiddenValley2", "Hidden Valley Ski"},
				repository.SnowResort{"Snow-Creek", "Snow Creek"},
			},
		},
	}

	for _, tt := range cases {
		if diff := cmp.Diff(parseStringToSnowResorts(tt.input), tt.output); diff != "" {
			t.Errorf("Diff: (-got +want)\n%s", diff)
		}
	}
}
