package api

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_handleClasses(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleClasses()
		})
	}
}

func Test_getClassesHandler(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getClassesHandler(tt.args.c)
		})
	}
}
