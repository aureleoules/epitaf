package api

import (
	"reflect"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name string
		want *jwt.GinJWTMiddleware
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthMiddleware(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
