package dotenvconfig

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

const (
	testKey         = "TEST_KEY"
	testKeyValue    = "123"
	testKeyValueInt = 123
)

// Load

func TestLoad_FailsWithNilFs(t *testing.T) {

	if err := Load(nil, "", Options{}); err == nil {
		t.Errorf("Load() error = %v, wantErr %v", err, true)
	}
}

func TestLoad_FailsWithMissingConfig(t *testing.T) {

	if err := Load(getFS(), "missing.env", Options{}); err == nil {
		t.Errorf("Load() error = %v, wantErr %v", err, true)
	}
}

func TestLoad_FailsWithInvalidConfig(t *testing.T) {

	if err := Load(getFS(), "invalid.env", Options{}); err == nil {
		t.Errorf("Load() error = %v, wantErr %v", err, true)
	}
}

func TestLoad_SetsNewKey(t *testing.T) {

	if err := Load(getFS(), "default.env", Options{}); err != nil {
		t.Errorf("Load() error = %v, wantErr %v", err, true)
	}

	if v := os.Getenv(testKey); v != testKeyValue {
		t.Errorf("Load() invalid value for key %s, expected %s, got %s", testKey, testKeyValue, v)
	}
}

func TestLoad_KeepsExistingKey(t *testing.T) {

	tmpVal := "789"
	_ = os.Setenv(testKey, tmpVal)
	defer func() {
		err := os.Unsetenv(testKey)
		if err != nil {
			panic(err)
		}
	}()

	if err := Load(getFS(), "default.env", Options{}); err != nil {
		t.Errorf("Load() error = %v, wantErr %v", err, true)
	}

	if v := os.Getenv(testKey); v != tmpVal {
		t.Errorf("Load() invalid value for key %s, expected %s, got %s", testKey, tmpVal, v)
	}
}

// GetKey

func TestGetKey(t *testing.T) {
	_ = os.Setenv(testKey, testKeyValue)
	defer func() {
		err := os.Unsetenv(testKey)
		if err != nil {
			panic(err)
		}
	}()

	v, err := GetKey(testKey)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if v != testKeyValue {
		t.Errorf("GetKey(): want %s, got %s", testKeyValue, v)
	}
}

func TestGetKey_Err(t *testing.T) {

	_, err := GetKey(testKey)
	if !errors.Is(err, ErrMissingKey) {
		t.Errorf("GetKey(): want %s, got %s", ErrMissingKey, err)
	}
}

// GetKeyMust

func TestGetKeyMust(t *testing.T) {
	_ = os.Setenv(testKey, testKeyValue)
	defer func() {
		err := os.Unsetenv(testKey)
		if err != nil {
			panic(err)
		}
	}()

	if v := GetKeyMust(testKey); v != testKeyValue {
		t.Errorf("GetKeyMust(): want %s, got %s", testKeyValue, v)
	}
}

func TestGetKeyMust_Panic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic for non-existent key")
		}
	}()

	GetKeyMust(testKey)
}

// GetKeyAsInt

func TestGetKeyAsInt(t *testing.T) {
	_ = os.Setenv(testKey, testKeyValue)
	defer func() {
		err := os.Unsetenv(testKey)
		if err != nil {
			panic(err)
		}
	}()

	v, err := GetKeyAsInt(testKey)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if v != testKeyValueInt {
		t.Errorf("GetKey(): want %d, got %d", testKeyValueInt, v)
	}
}

func TestGetKeyAsInt_Missing(t *testing.T) {

	_, err := GetKeyAsInt(testKey)
	if !errors.Is(err, ErrMissingKey) {
		t.Errorf("GetKeyAsInt(): want %s, got %s", ErrMissingKey, err)
	}
}

func TestGetKeyAsInt_ConversionErr(t *testing.T) {
	_ = os.Setenv(testKey, "ABC")
	defer func() {
		err := os.Unsetenv(testKey)
		if err != nil {
			panic(err)
		}
	}()

	_, err := GetKeyAsInt(testKey)
	if !errors.Is(err, ErrConversion) {
		t.Errorf("GetKeyAsInt(): want %s, got %s", ErrConversion, err)
	}
}

// GetKeyAsIntMust

func TestGetKeyAsIntMust(t *testing.T) {
	_ = os.Setenv(testKey, testKeyValue)
	defer func() {
		err := os.Unsetenv(testKey)
		if err != nil {
			panic(err)
		}
	}()

	if v := GetKeyAsIntMust(testKey); v != testKeyValueInt {
		t.Errorf("GetKeyMust(): want %d, got %d", testKeyValueInt, v)
	}
}

func TestGetKeyAsIntMust_Panic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected a panic for non-existent key")
		}
	}()

	GetKeyAsIntMust(testKey)
}

func getFS() fs.FS {
	return os.DirFS("testconfig/")
}
