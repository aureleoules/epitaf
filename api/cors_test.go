package api

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_cors(t *testing.T) {
	tests := []struct {
		name string
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cors() = %v, want %v", got, tt.want)
			}
		})
	}
}
