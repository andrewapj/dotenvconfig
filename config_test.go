package dotenvconfig

import (
	"io/fs"
	"os"
	"reflect"
	"strconv"
	"testing"
)

const testKey = "TEST_KEY"

func TestLoad(t *testing.T) {
	type args struct {
		fSys fs.FS
		opts Options
	}

	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedKey   string
		expectedValue string
		setup         func()
		tearDown      func()
	}{
		{
			name: "should get error with nil fs",
			args: args{
				fSys: nil,
				opts: Options{},
			},
			wantErr:  true,
			setup:    func() {},
			tearDown: func() {},
		},
		{
			name: "should get config with profile set by environment variable",
			args: args{
				fSys: getFS(),
				opts: Options{ProfileKey: "key"},
			},
			wantErr:       false,
			expectedKey:   testKey,
			expectedValue: "000",
			setup:         func() { _ = os.Setenv("key", "envvar") },
			tearDown: func() {
				_ = os.Unsetenv("key")
				_ = os.Unsetenv(testKey)
			},
		},
		{
			name: "should get default profile, when environment variable is missing",
			args: args{
				fSys: getFS(),
				opts: Options{ProfileKey: "missingKey"},
			},
			wantErr:       false,
			expectedKey:   testKey,
			expectedValue: "123",
			setup:         func() {},
			tearDown:      func() { _ = os.Unsetenv(testKey) },
		},
		{
			name: "should get config when profile set explicitly",
			args: args{
				fSys: getFS(),
				opts: Options{Profile: "custom"},
			},
			wantErr:       false,
			expectedKey:   testKey,
			expectedValue: "789",
			setup:         func() {},
			tearDown:      func() { _ = os.Unsetenv(testKey) },
		},
		{
			name: "should get config with default profile",
			args: args{
				fSys: getFS(),
				opts: Options{},
			},
			wantErr:       false,
			expectedKey:   "TEST_KEY2",
			expectedValue: "456",
			setup:         func() {},
			tearDown: func() {
				_ = os.Unsetenv(testKey)
				_ = os.Unsetenv("TEST_KEY2")
			},
		},
		{
			name: "should get config from environment with Load not overriding",
			args: args{
				fSys: getFS(),
				opts: Options{Profile: "custom"}},
			wantErr:       false,
			expectedKey:   testKey,
			expectedValue: "preserved",
			setup: func() {
				_ = os.Setenv(testKey, "preserved")
			},
			tearDown: func() { _ = os.Unsetenv(testKey) },
		},
		{
			name: "should get error with profile that points to a missing .env file",
			args: args{
				fSys: getFS(),
				opts: Options{Profile: "missing"},
			},
			wantErr:  true,
			setup:    func() {},
			tearDown: func() {},
		},
		{
			name: "should get error with invalid config",
			args: args{
				fSys: getFS(),
				opts: Options{Profile: "invalid"},
			},
			wantErr:  true,
			setup:    func() {},
			tearDown: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.tearDown()
			if err := Load(tt.args.fSys, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			v := os.Getenv(tt.expectedKey)
			if !reflect.DeepEqual(tt.expectedValue, v) {
				t.Errorf("Load() got = %v, want %v", v, tt.expectedValue)
			}
		})
	}
}

// GetKey

func TestGetKey_ExistingKey(t *testing.T) {
	value := "TEST_VALUE"
	_ = os.Setenv(testKey, value)
	defer os.Unsetenv(testKey)

	if result := GetKey(testKey); result != value {
		t.Errorf("Expected %s, got %s", value, result)
	}
}

func TestGetKey_NonExistentKey_NoPanic(t *testing.T) {
	nonExistentKey := "NON_EXISTENT_KEY"

	if result := GetKey(nonExistentKey); result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}
}

// GetKeyMust

func TestGetKeyMust(t *testing.T) {
	value := "TEST_VALUE"
	_ = os.Setenv(testKey, value)
	defer os.Unsetenv(testKey)

	if result := GetKeyMust(testKey); result != value {
		t.Errorf("Expected %s, got %s", value, result)
	}
}

func TestGetKeyMust_WithPanic(t *testing.T) {
	nonExistentKey := "NON_EXISTENT_KEY"

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic for non-existent key with panicOnError set to true")
		}
	}()

	GetKeyMust(nonExistentKey)
}

// GetKeyAsInt

func TestGetKeyAsInt_ValidInt(t *testing.T) {
	key := testKey
	value := "12345"
	expectedIntValue, _ := strconv.Atoi(value)
	_ = os.Setenv(key, value)
	defer os.Unsetenv(key)

	if result := GetKeyAsInt(key); result != expectedIntValue {
		t.Errorf("Expected %d, got %d", expectedIntValue, result)
	}
}

func TestGetKeyAsInt_InvalidInt(t *testing.T) {
	key := testKey
	value := "abc"
	expectedIntValue := 0
	_ = os.Setenv(key, value)
	defer os.Unsetenv(key)

	if result := GetKeyAsInt(key); result != expectedIntValue {
		t.Errorf("Expected %d, got %d", expectedIntValue, result)
	}
}

func TestGetKeyAsInt_NonExistentKey(t *testing.T) {
	nonExistentKey := "NON_EXISTENT_INT_KEY"

	if result := GetKeyAsInt(nonExistentKey); result != 0 {
		t.Errorf("Expected %d, got %d", 0, result)
	}
}

// GetKeyAsIntMust

func TestGetKeyAsIntMust(t *testing.T) {
	key := testKey
	value := "12345"
	expectedIntValue, _ := strconv.Atoi(value)
	_ = os.Setenv(key, value)
	defer os.Unsetenv(key)

	if result := GetKeyAsIntMust(key); result != expectedIntValue {
		t.Errorf("Expected %d, got %d", expectedIntValue, result)
	}
}

func TestGetKeyAsIntMust_InvalidInt(t *testing.T) {
	key := testKey
	value := "abc"
	_ = os.Setenv(key, value)
	defer os.Unsetenv(key)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic for non-existent key with panicOnError set to true")
		}
	}()

	GetKeyAsIntMust(key)
}

func TestGetKeyAsIntMust_MissingInt(t *testing.T) {
	nonExistentKey := "NON_EXISTENT_INT_KEY"

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic for non-existent key with panicOnError set to true")
		}
	}()

	GetKeyAsIntMust(nonExistentKey)
}

func getFS() fs.FS {
	return os.DirFS("testconfig/")
}
