package redis

import (
	"fmt"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
	"reflect"
	"testing"
)

func TestMapToStruct(t *testing.T) {
	sr := &domain.SnowResort{}
	type args struct {
		mapVal map[string]string
		val    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "etst",
			args: args{
				mapVal: map[string]string{
					"id":        "a",
					"name":      "A",
					"searchKey": "a",
					"elevation": "Elevation",
					"region":    "W",
				},
				val: sr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(sr)
		})
	}
}

func Test_structToMap(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := structToMap(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("structToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
