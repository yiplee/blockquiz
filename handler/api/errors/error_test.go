package errors

import (
	"errors"
	"testing"
)

func TestUnwrap(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
		wantMsg  string
	}{
		{
			name:     "Unauthorized",
			args:     args{err: ErrUnauthorized},
			wantCode: 401,
			wantMsg:  "Unauthorized",
		},
		{
			name:     "simple error",
			args:     args{err: errors.New("failed")},
			wantCode: 0,
			wantMsg:  "failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotMsg := Unwrap(tt.args.err)
			if gotCode != tt.wantCode {
				t.Errorf("Unwrap() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("Unwrap() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}
