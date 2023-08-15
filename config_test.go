package dotenvconfig

import (
	"context"
	"io/fs"
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		fSys fs.FS
		opts Options
	}
	tests := []struct {
		name     string
		setup    func()
		tearDown func()
		args     args
		want     Config
		wantErr  bool
	}{
		{
			name: "should get correct config with env variable",
			setup: func() {
				_ = os.Setenv("ENV_VAR", "envvar")
			},
			tearDown: func() {
				_ = os.Unsetenv("ENV_VAR")
			},
			args: args{
				fSys: getFS(),
				opts: Options{
					environment:    "ignored",
					environmentKey: "ENV_VAR",
				},
			},
			want:    Config{configMap: map[string]string{"TEST_KEY": "000"}},
			wantErr: false,
		},
		{
			name:     "should get correct config with specific environment",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				fSys: getFS(),
				opts: Options{environment: "custom"},
			},
			want:    Config{configMap: map[string]string{"TEST_KEY": "789"}},
			wantErr: false,
		},
		{
			name:     "should get correct config with default environment",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				fSys: getFS(),
				opts: Options{
					jsonLogging:    true,
					loggingEnabled: true,
				},
			},
			want: Config{configMap: map[string]string{
				"TEST_KEY":  "123",
				"TEST_KEY2": "456",
			}},
			wantErr: false,
		},
		{
			name:     "should get error with nil fs",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				fSys: nil,
				opts: Options{},
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name:     "should get error with invalid config",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				fSys: getFS(),
				opts: Options{environment: "invalid"},
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name:     "should get error with missing config",
			setup:    func() {},
			tearDown: func() {},
			args: args{
				fSys: getFS(),
				opts: Options{environment: "missing"},
			},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.tearDown()
			got, err := NewConfig(tt.args.fSys, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetKey(t *testing.T) {
	type fields struct {
		configMap map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name     string
		setup    func()
		tearDown func()
		fields   fields
		args     args
		want     string
	}{
		{
			name: "should get key with env var override",
			setup: func() {
				_ = os.Setenv("TEST_KEY", "override")
			},
			tearDown: func() {
				_ = os.Unsetenv("TEST_KEY")
			},
			fields: fields{configMap: map[string]string{"TEST_KEY": "123"}},
			args:   args{key: "TEST_KEY"},
			want:   "override",
		},
		{
			name:     "should get key without env var override",
			setup:    func() {},
			tearDown: func() {},
			fields:   fields{configMap: map[string]string{"TEST_KEY": "123"}},
			args:     args{key: "TEST_KEY"},
			want:     "123",
		},
		{
			name:     "should get zero value with missing key",
			setup:    func() {},
			tearDown: func() {},
			fields:   fields{configMap: make(map[string]string)},
			args:     args{key: "TEST_KEY"},
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.tearDown()
			c := Config{
				configMap: tt.fields.configMap,
			}
			if got := c.GetKey(tt.args.key); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetKeyAsInt(t *testing.T) {
	type fields struct {
		configMap map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "should get valid int value",
			fields: fields{configMap: map[string]string{"TEST_KEY": "123"}},
			args:   args{key: "TEST_KEY"},
			want:   123,
		},
		{
			name:   "should get zero value with invalid int",
			fields: fields{configMap: map[string]string{"TEST_KEY": "ABC"}},
			args:   args{key: "TEST_KEY"},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				configMap: tt.fields.configMap,
			}
			if got := c.GetKeyAsInt(tt.args.key); got != tt.want {
				t.Errorf("GetKeyAsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToContext(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg Config
	}
	tests := []struct {
		name    string
		args    args
		want    context.Context
		wantErr bool
	}{
		{
			name: "should create ctx with config",
			args: args{
				ctx: context.Background(),
				cfg: Config{configMap: make(map[string]string)},
			},
			want: context.WithValue(context.Background(),
				contextKey, Config{make(map[string]string)}),
			wantErr: false,
		},
		{
			name: "should get error with nil ctx",
			args: args{
				ctx: nil,
				cfg: Config{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToContext(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "should get config from ctx",
			args: args{ctx: context.WithValue(context.Background(),
				contextKey, Config{make(map[string]string)})},
			want:    Config{make(map[string]string)},
			wantErr: false,
		},
		{
			name:    "should get error with nil ctx",
			args:    args{ctx: nil},
			want:    Config{},
			wantErr: true,
		},
		{
			name:    "should get error with ctx that has no config",
			args:    args{ctx: context.Background()},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "should get error with ctx that has non config type",
			args: args{ctx: context.WithValue(context.Background(),
				contextKey, 1)},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getFS() fs.FS {
	return os.DirFS("testconfig/")
}
