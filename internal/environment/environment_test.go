package environment

import (
	"os"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	type args struct {
		environmentKey string
		environment    string
	}
	tests := []struct {
		name     string
		setup    func()
		tearDown func()
		args     args
		want     string
	}{
		{
			name: "should get environment from environment key",
			setup: func() {
				_ = os.Setenv("ENV_KEY", "env_value")
			},
			tearDown: func() {
				_ = os.Unsetenv("ENV_KEY")
			},
			args: args{
				environmentKey: "ENV_KEY",
				environment:    "ignored",
			},
			want: "env_value.env",
		},
		{
			name:     "should get environment from explicit environment variable",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				environmentKey: "MISSING_KEY",
				environment:    "env",
			},
			want: "env.env",
		},
		{
			name:     "should get default environment with no environment key or explicit variable",
			setup:    func() {},
			tearDown: func() {},
			args:     args{},
			want:     "default.env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.tearDown()
			if got := GetEnvironment(tt.args.environmentKey, tt.args.environment); got != tt.want {
				t.Errorf("GetEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
