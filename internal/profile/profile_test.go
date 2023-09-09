package profile

import (
	"os"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	type args struct {
		profileKey string
		profile    string
	}
	tests := []struct {
		name     string
		setup    func()
		tearDown func()
		args     args
		want     string
	}{
		{
			name: "should get profile from profile key",
			setup: func() {
				_ = os.Setenv("PROFILE_KEY", "profile_value")
			},
			tearDown: func() {
				_ = os.Unsetenv("PROFILE_KEY")
			},
			args: args{
				profileKey: "PROFILE_KEY",
				profile:    "ignored",
			},
			want: "profile_value.env",
		},
		{
			name:     "should get profile from explicit profile variable",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				profileKey: "MISSING_KEY",
				profile:    "profile",
			},
			want: "profile.env",
		},
		{
			name:     "should get default profile with no profile key or explicit variable",
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
			if got := GetProfile(tt.args.profileKey, tt.args.profile); got != tt.want {
				t.Errorf("GetProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
