// Package ru provides Russian national cryptography (GOST) implementation.
package ru

import (
	"crypto/tls"

	"github.com/vpnclient/https-vpn/crypto"
	rutls "github.com/vpnclient/https-vpn/crypto/ru/tls"
)

// Provider implements crypto.Provider for Russian cryptography (GOST).
type Provider struct{}

// Name returns the provider identifier.
func (p *Provider) Name() string {
	return "ru"
}

// ConfigureTLS applies GOST cryptography settings to tls.Config.
func (p *Provider) ConfigureTLS(cfg *tls.Config) error {
	// Note: Standard Go TLS doesn't support GOST cipher suites natively.
	// We use the custom TLS implementation in crypto/ru/tls.
	
	// Set minimum TLS version to TLS 1.2 (commonly used with GOST)
	if cfg.MinVersion < rutls.VersionTLS12 {
		cfg.MinVersion = rutls.VersionTLS12
	}
	
	// Set supported cipher suites
	cfg.CipherSuites = p.SupportedCipherSuites()

	return nil
}

// SupportedCipherSuites returns the list of supported GOST cipher suite IDs.
func (p *Provider) SupportedCipherSuites() []uint16 {
	return rutls.SupportedCipherSuites()
}

func init() {
	crypto.Register(&Provider{})
}
