package cli

import (
	"reflect"
	"testing"
)

func TestCombineStructs(t *testing.T) {
	type args struct {
		structs []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Combine structs",
			args: args{
				structs: []interface{}{
					&struct{ I int }{I: 10},
					struct{ S string }{S: "string"},
				},
			},
			want: &struct {
				I int
				S string
			}{I: 10, S: "string"},
			wantErr: false,
		},
		{
			name: "Combine structs, invalid parameters",
			args: args{
				structs: []interface{}{
					&struct{ I int }{I: 10},
					[]string{},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineStructs(tt.args.structs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CombineStructs() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CombineStructs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCopyStruct(t *testing.T) {
	type args struct {
		src  interface{}
		dest interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			args: args{
				src: &struct{ I int }{I: 10},
				dest: &struct {
					I int
					S string
				}{S: "string"},
			},
			want: &struct {
				I int
				S string
			}{I: 10, S: "string"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyStruct(tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("CopyStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dest, tt.want) {
				t.Errorf("dest = %#v, want %#v", tt.args.dest, tt.want)
			}
		})
	}
}
