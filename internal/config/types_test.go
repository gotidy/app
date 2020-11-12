package config

import (
	"testing"
)

func TestConnection_StringURL(t *testing.T) {
	tests := []struct {
		connection Connection
		want       string
	}{
		{
			connection: Connection{
				Scheme:  "http",
				User:    UserCredential{User: "", Password: nil},
				Address: NetAddress{Host: ""},
			},
			want: "",
		},
		{
			connection: Connection{
				Scheme:  "http",
				User:    UserCredential{User: "", Password: nil},
				Address: NetAddress{Host: "something.com"},
			},
			want: "something.com",
		},
		{
			connection: Connection{
				Scheme:  "http",
				User:    UserCredential{User: "", Password: nil},
				Address: NetAddress{Host: "something.com", Port: 80},
			},
			want: "something.com:80",
		},
		{
			connection: Connection{
				Scheme:  "http",
				User:    UserCredential{User: "john", Password: nil},
				Address: NetAddress{Host: "something.com", Port: 80},
			},
			want: "john@something.com:80",
		},
		{
			connection: Connection{
				Scheme:  "http",
				User:    NewUserCredential("john", ""),
				Address: NetAddress{Host: "something.com", Port: 80},
			},
			want: "john:@something.com:80",
		},
		{
			connection: Connection{
				Scheme:  "http",
				User:    NewUserCredential("john", "qwerty"),
				Address: NetAddress{Host: "something.com", Port: 80},
			},
			want: "john:qwerty@something.com:80",
		},
		{
			connection: Connection{
				Scheme:  "http",
				User:    NewUserCredential("", "qwerty"),
				Address: NetAddress{Host: "something.com", Port: 80},
			},
			want: ":qwerty@something.com:80",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.connection.String(); got != tt.want {
				t.Errorf("Connection.String() = %v, want %v", got, tt.want)
			}
			if got, want := tt.connection.URL(), "http://"+tt.want+"/"; got != want {
				t.Errorf("Connection.String() = %v, want %v", got, want)
			}
		})
	}
}

func TestNetAddress_Network(t *testing.T) {
	type fields struct {
		Host string
		Port int
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{
			fields: fields{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			a := NetAddress{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if got := a.Network(); got != tt.want {
				t.Errorf("NetAddress.Network() = %v, want %v", got, tt.want)
			}
		})
	}
}
