package composeaddresstranslator

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type ComposeAddressTranslator struct {
	Map map[string]string
}

// New creates a ComposeAddressTranslator with the given translation map form.
// The keys and values for the translation map are the host with port.
// e.g.
// New(map[string]string{
//     "127.0.0.1:9042": "192.168.1.1:1234",
//     "127.0.0.2:9042": "192.168.1.2:5678",
//     "127.0.0.3:9042": "192.168.1.3:9012",
// })
func New(m map[string]string) ComposeAddressTranslator {
	return ComposeAddressTranslator{
		Map: m,
	}
}

// NewFromJSONString creates a ComposeAddressTranslator with the given translation map in JSON string form.
// e.g.
// NewFromJSONString(`{
//         "127.0.0.1:9042": "192.168.1.1:1234",
//         "127.0.0.2:9042": "192.168.1.2:5678",
//         "127.0.0.3:9042": "192.168.1.3:9012"
//     }`)
func NewFromJSONString(s string) (ComposeAddressTranslator, error) {
	cat := ComposeAddressTranslator{
		Map: make(map[string]string),
	}
	err := json.Unmarshal([]byte(s), &cat.Map)
	return cat, err
}

// Translate implements the AddressTranslator interface for the gocql driver.
// In the case that the translation map contains
func (cat ComposeAddressTranslator) Translate(addr net.IP, port int) (net.IP, int) {
	if host, ok := cat.Map[fmt.Sprintf("%s:%s", addr.String(), strconv.Itoa(port))]; ok {
		a, p, err := net.SplitHostPort(host)
		if err != nil {
			return addr, port
		}
		port, err = strconv.Atoi(p)
		if err != nil {
			return addr, port
		}
		// See if it was passed as IP addresses
		ip := net.ParseIP(a)
		if ip != nil {
			return ip, port
		}
		// We have a hostname, do DNS lookup and use first response (gocql makes same assumption
		// about multiple A records)
		ips, err := net.LookupIP(a)
		if err != nil || len(ips) < 1 {
			return addr, port
		}
		return ips[0], port
	}
	return addr, port
}

func (cat ComposeAddressTranslator) ContactPoints() []string {
	s := make([]string, len(cat.Map))
	i := 0
	// Return the _internal_ IP addresses as gocql seems to expect these
	for h := range cat.Map {
		s[i] = h
		i++
	}
	return s
}
