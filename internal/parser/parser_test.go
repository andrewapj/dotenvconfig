package parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "should parse valid config",
			args:    args{data: getValidConfig()},
			want:    map[string]string{"TEST_KEY": "123", "TEST_KEY2": "456"},
			wantErr: false,
		},
		{
			name:    "should not parse invalid config",
			args:    args{data: getInvalidConfig()},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should parse config with empty lines",
			args:    args{data: getValidConfigWithEmptyLinesAndComments()},
			want:    map[string]string{"TEST_KEY": "123", "TEST_KEY2": "456"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getValidConfig() []byte {
	return []byte("TEST_KEY=123\nTEST_KEY2=456")
}

func getValidConfigWithEmptyLinesAndComments() []byte {
	return []byte(`
	TEST_KEY=123

	# Comment

	TEST_KEY2=456`)
}

func getInvalidConfig() []byte {
	return []byte("TEST_KEY,123")
}
