package fr

import (
	"testing"
	"github.com/vpnclient/https-vpn/crypto"
)

func TestFRProviderRegistration(t *testing.T) {
	p, ok := crypto.Get("fr")
	if !ok {
		t.Fatal("FR provider not registered")
	}
	if p.Name() != "fr" {
		t.Errorf("Expected name 'fr', got '%s'", p.Name())
	}
}
