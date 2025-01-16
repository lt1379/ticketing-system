package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {
	now, _, _ := time.Now().UTC().Date()
	tests := []struct {
		name string
		want int
	}{
		{
			name: "TestGetCurrentTime - 1",
			want: now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := GetCurrentTime().Date(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrentTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
