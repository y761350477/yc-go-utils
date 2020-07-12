package convert

import (
	"testing"
)

func TestStringToInt(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{str: `1`},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToInt(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringToInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
