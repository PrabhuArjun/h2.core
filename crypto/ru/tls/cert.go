// Package tls provides GOST TLS implementation.
package tls

import (
	"crypto/rand"
	"fmt"
	"github.com/vpnclient/https-vpn/crypto/ru/gost"
)

// GenerateGOSTKeyPair generates a GOST R 34.10-2012 key pair.
func GenerateGOSTKeyPair() (*gost.PrivateKey, error) {
	// Using default paramSetA for 256-bit keys
	curve := gost.TC26256A()
	priv, err := gost.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate GOST key: %w", err)
	}
	return priv, nil
}
