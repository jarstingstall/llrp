package llrp

import (
	"testing"
)

func TestItReturnsNewClientWithHostAndPortSet(t *testing.T) {
	host := "10.0.0.29"
	port := "5084"
	c := NewClient(host, port)
	if c.Host != host {
		t.Fatalf("Client.Host not set properly. Expected %s but got %s", host, c.Host)
	}
	if c.Port != port {
		t.Fatalf("Client.Port not set properly. Expected %s but got %s", port, c.Port)
	}
}

func TestClientCanConnectToServerAndClose(t *testing.T) {
	host := "10.0.0.29"
	port := "5084"
	c := NewClient(host, port)
	err := c.Connect()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
}
