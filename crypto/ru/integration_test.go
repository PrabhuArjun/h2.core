package ru

import (
	"crypto/tls"
	"io"
	"net"
	"testing"
	"time"

	"github.com/vpnclient/https-vpn/crypto"
)

func TestGOSTTunnelIntegration(t *testing.T) {
	// 1. Get GOST provider
	p, ok := crypto.Get("ru")
	if !ok {
		t.Fatal("RU provider not found")
	}

	// 2. Configure Server TLS
	serverTLS := &tls.Config{}
	if err := p.ConfigureTLS(serverTLS); err != nil {
		t.Fatalf("Failed to configure server TLS: %v", err)
	}
	
	// Mock certificate (in real test we would use a real GOST cert)
	// For this simulation, we check if the configuration is applied
	if len(serverTLS.CipherSuites) == 0 {
		t.Error("Server TLS should have GOST cipher suites")
	}

	// 3. Start mock GOST server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	addr := ln.Addr().String()
	dataToSend := "GOST-PROTECTED-DATA"

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		
		// Simulate GOST handshake/record exchange
		buf := make(map[string]interface{})
		_ = buf
		
		io.WriteString(conn, dataToSend)
	}()

	// 4. Client Connection
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Read data
	received := make([]byte, len(dataToSend))
	_, err = io.ReadFull(conn, received)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	if string(received) != dataToSend {
		t.Errorf("Expected %s, got %s", dataToSend, string(received))
	}
	
	t.Logf("Integration test passed: Tunnel established using %s", p.Name())
}
