package config

import (
	"io/fs"
	"os"
	"reflect"
	"testing"
)

const testKey = "TEST_KEY"
const testValue = "123"
const badIntKey = "BAD_KEY"

func TestNewConfig(t *testing.T) {
	got := NewConfig(configDirFS())
	want := Config{
		configMap:  makeEmptyMap(),
		fs:         configDirFS(),
		currentEnv: "",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() got = %v, want %v", got, want)
	}
}

func TestConfig_WithEnvironment(t *testing.T) {
	got := NewConfig(configDirFS()).WithEnvironment("custom")
	want := Config{
		configMap:  makeEmptyMap(),
		fs:         configDirFS(),
		currentEnv: "custom",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() got = %v, want %v", got, want)
	}
}

func TestConfig_Load(t *testing.T) {
	type fields struct {
		configMap  map[string]string
		fs         fs.FS
		currentEnv string
	}
	tests := []struct {
		name    string
		fields  fields
		want    Config
		wantErr bool
	}{
		{
			name: "should get config with default environment",
			fields: fields{
				configMap:  makeEmptyMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			want: Config{
				configMap:  makeDefaultConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			wantErr: false,
		},
		{
			name: "should get config with custom environment",
			fields: fields{
				configMap:  makeEmptyMap(),
				fs:         configDirFS(),
				currentEnv: "custom",
			},
			want: Config{
				configMap:  makeCustomConfigMap(),
				fs:         configDirFS(),
				currentEnv: "custom",
			},
			wantErr: false,
		},
		{
			name: "should get error with invalid fs in config",
			fields: fields{
				configMap:  makeEmptyMap(),
				fs:         nil,
				currentEnv: "",
			},
			want: Config{
				configMap:  makeEmptyMap(),
				fs:         nil,
				currentEnv: "",
			},
			wantErr: true,
		},
		{
			name: "should get error with fs that doesn't contain a config file",
			fields: fields{
				configMap:  makeEmptyMap(),
				fs:         os.DirFS(".."),
				currentEnv: "",
			},
			want: Config{
				configMap:  makeEmptyMap(),
				fs:         os.DirFS(".."),
				currentEnv: "",
			},
			wantErr: true,
		},
		{
			name: "should get error with invalid config file",
			fields: fields{
				configMap:  makeEmptyMap(),
				fs:         configDirFS(),
				currentEnv: "invalid",
			},
			want: Config{
				configMap:  makeEmptyMap(),
				fs:         configDirFS(),
				currentEnv: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				configMap:  tt.fields.configMap,
				fs:         tt.fields.fs,
				currentEnv: tt.fields.currentEnv,
			}
			got, err := c.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetKey(t *testing.T) {
	type fields struct {
		configMap  map[string]string
		fs         fs.FS
		currentEnv string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		setup    func()
		teardown func()
		want     string
		wantErr  bool
	}{
		{
			name: "should get key from config",
			fields: fields{
				configMap:  makeDefaultConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			args:     args{key: testKey},
			setup:    func() {},
			teardown: func() {},
			want:     testValue,
			wantErr:  false,
		},
		{
			name: "should get key from config, using os environment override",
			fields: fields{
				configMap:  makeDefaultConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			args: args{key: testKey},
			setup: func() {
				_ = os.Setenv(testKey, "NEW_VALUE")
			},
			teardown: func() {
				_ = os.Unsetenv(testKey)
			},
			want:    "NEW_VALUE",
			wantErr: false,
		},
		{
			name: "should not find value for key",
			fields: fields{
				configMap:  makeDefaultConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			args:     args{key: "MISSING_KEY"},
			setup:    func() {},
			teardown: func() {},
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				configMap:  tt.fields.configMap,
				fs:         tt.fields.fs,
				currentEnv: tt.fields.currentEnv,
			}
			tt.setup()
			defer tt.teardown()
			got, err := c.GetKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetKeyAsInt(t *testing.T) {
	type fields struct {
		configMap  map[string]string
		fs         fs.FS
		currentEnv string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "should get int key from config",
			fields: fields{
				configMap:  makeDefaultConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			args:    args{key: testKey},
			want:    123,
			wantErr: false,
		},
		{
			name: "should not find int key",
			fields: fields{
				configMap:  makeDefaultConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			args:    args{key: "MISSING_KEY"},
			want:    0,
			wantErr: true,
		},
		{
			name: "should get error with invalid int value for key",
			fields: fields{
				configMap:  makeInvalidIntConfigMap(),
				fs:         configDirFS(),
				currentEnv: "",
			},
			args:    args{key: badIntKey},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				configMap:  tt.fields.configMap,
				fs:         tt.fields.fs,
				currentEnv: tt.fields.currentEnv,
			}
			got, err := c.GetKeyAsInt(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeyAsInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKeyAsInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func configDirFS() fs.FS {
	return os.DirFS("internal/testconfig")
}

func makeEmptyMap() map[string]string {
	return map[string]string{}
}

func makeDefaultConfigMap() map[string]string {
	return map[string]string{testKey: testValue}
}

func makeCustomConfigMap() map[string]string {
	return map[string]string{testKey: "789"}
}

func makeInvalidIntConfigMap() map[string]string {
	return map[string]string{badIntKey: "ABC"}
}
