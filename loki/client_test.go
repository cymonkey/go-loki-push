package loki

import (
	"context"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		cfg Config
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.cfg, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
