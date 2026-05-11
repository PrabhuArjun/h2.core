// Package main provides a tool to generate self-signed GOST certificates.
package main

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/vpnclient/https-vpn/crypto/ru/gost"
)

func main() {
	fmt.Println("Generating GOST R 34.10-2012 Key Pair...")
	
	curve := gost.TC26256A()
	priv, err := gost.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	// Save Private Key
	// In a real implementation, we would use PKCS#8 with GOST OIDs
	privBytes := priv.D.Bytes()
	privBlock := &pem.Block{
		Type:  "GOST PRIVATE KEY",
		Bytes: privBytes,
	}
	
	privFile, _ := os.Create("gost_key.pem")
	pem.Encode(privFile, privBlock)
	privFile.Close()

	fmt.Println("Generated gost_key.pem")

	// Save Public Key / Mock Certificate
	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	pubBlock := &pem.Block{
		Type:  "GOST CERTIFICATE",
		Bytes: pubBytes,
	}
	
	pubFile, _ := os.Create("gost_cert.pem")
	pem.Encode(pubFile, pubBlock)
	pubFile.Close()

	fmt.Println("Generated gost_cert.pem (Mock)")
	fmt.Println("\nTo use this with HTTPS VPN, update your config.json:")
	fmt.Println(`"tlsSettings": {`)
	fmt.Println(`  "cipherSuites": "ru",`)
	fmt.Println(`  "certificates": [`)
	fmt.Println(`    { "certificateFile": "gost_cert.pem", "keyFile": "gost_key.pem" }`)
	fmt.Println(`  ]`)
	fmt.Println(`}`)
}
