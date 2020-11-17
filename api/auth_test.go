package api

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_handleAuth(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleAuth()
		})
	}
}

func Test_authenticateHandler(t *testing.T) {
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
			authenticateHandler(tt.args.c)
		})
	}
}

func Test_callbackHandler(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := callbackHandler(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("callbackHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("callbackHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
