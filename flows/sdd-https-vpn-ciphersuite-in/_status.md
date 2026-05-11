# Status: sdd-https-vpn-ciphersuite-in

## Current Phase

REQUIREMENTS (drafting)

## Last Updated

2026-04-24 by Claude

## Blockers

- Awaiting user clarification on requirements (India does not have national cryptographic algorithms like China's SM or Russia's GOST)

## Progress

- [ ] Requirements drafted ← current
- [ ] Requirements approved
- [ ] Specifications drafted
- [ ] Specifications approved
- [ ] Plan drafted
- [ ] Plan approved
- [ ] Implementation started
- [ ] Implementation complete

## Context Notes

Key decisions and context for resuming:

- Research shows India follows international standards (NIST/ISO)
- India CCA (Controller of Certifying Authorities) uses RSA, ECDSA, SHA-2
- No India-specific national cryptographic algorithms found (unlike CN, RU)
- Need to clarify with user what "India cryptography" means:
  - Option A: CCA-compliant profile (RSA-2048+, ECDSA P-256/P-384, SHA-256/SHA-384)
  - Option B: CERT-In recommended TLS settings
  - Option C: Something else?

## Research Notes

Sources consulted:
- CCA (Controller of Certifying Authorities): cca.gov.in
- CERT-In (Indian Computer Emergency Response Team)
- India PKI Forum: indiapki.org

Findings:
- India uses standard international algorithms (RSA, ECDSA, AES, SHA-2)
- CCA Interoperability Guidelines specify certificate profiles
- SSL/TLS for banking is regulated by RBI (Reserve Bank of India)
- No national cipher algorithms like SM (China) or GOST (Russia)

## Fork History

N/A - New flow
