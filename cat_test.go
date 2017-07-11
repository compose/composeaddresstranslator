package composeaddresstranslator

import (
	"net"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	cat := ComposeAddressTranslator{map[string]string{
		"127.0.0.1:0": "192.168.1.1:1",
		"127.0.0.1:2": "192.168.1.1:3",
	}}
	for _, tt := range []struct {
		want     net.IP
		wantPort int
		host     net.IP
		port     int
	}{
		{
			want:     net.ParseIP("192.168.1.1"),
			wantPort: 1,
			host:     net.ParseIP("127.0.0.1"),
			port:     0,
		},
		{
			want:     net.ParseIP("192.168.1.1"),
			wantPort: 3,
			host:     net.ParseIP("127.0.0.1"),
			port:     2,
		},
		{
			want:     net.ParseIP("127.0.0.2"),
			wantPort: 1,
			host:     net.ParseIP("127.0.0.2"),
			port:     1,
		},
	} {
		got, gotPort := cat.Translate(tt.host, tt.port)
		if !got.Equal(tt.want) {
			t.Errorf("host got %v, wanted %v", got, tt.want)
		}
		if gotPort != tt.wantPort {
			t.Errorf("port got %v, wanted %v", gotPort, tt.wantPort)
		}
	}
}

func TestContactPoints(t *testing.T) {
	got := ComposeAddressTranslator{map[string]string{
		"127.0.0.1:0": "192.168.1.1:1",
		"127.0.0.1:2": "192.168.1.1:3",
	}}.ContactPoints()
	want := []string{
		"127.0.0.1:0",
		"127.0.0.1:2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, wanted: %v", got, want)
	}
}
