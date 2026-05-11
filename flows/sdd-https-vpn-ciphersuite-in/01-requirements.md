# Requirements: India Cryptography Provider

> Version: 0.1 (DRAFT)
> Status: PENDING CLARIFICATION
> Last Updated: 2026-04-24

## Problem Statement

Добавить поддержку индийской криптографии в систему HTTPS VPN.

**Важное замечание:** В отличие от Китая (SM2/SM3/SM4) или России (ГОСТ), Индия **не имеет собственных национальных криптографических алгоритмов**. Индийские регуляторы (CCA, CERT-In, RBI) требуют использования международных стандартов.

Система уже поддерживает:
- "us" - стандартная криптография (RSA/ECDSA/AES/SHA)
- "ru" - российская криптография (ГОСТ Р 34.10/34.11/34.12)
- "cn" - китайская криптография (SM2/SM3/SM4/SM9)
- "uk" - британские рекомендации NCSC (TLS 1.3, AES-256-GCM)
- "fr" - французские рекомендации ANSSI

## Research Findings

### India CCA (Controller of Certifying Authorities)

Источник: [cca.gov.in](https://cca.gov.in)

CCA регулирует цифровые подписи в Индии по IT Act 2000. Поддерживаемые алгоритмы:
- **Подписи**: RSA (≥2048 бит), ECDSA (P-256, P-384)
- **Хэш**: SHA-256, SHA-384, SHA-512
- **Кривые**: NIST P-256, P-384

### India CERT-In

Источник: [cert-in.org.in](https://cert-in.org.in)

CERT-In выпускает рекомендации по кибербезопасности, но не специфицирует уникальные криптографические алгоритмы.

### Reserve Bank of India (RBI)

Для банковских систем RBI требует SSL/TLS с международными стандартами.

## Open Questions (для уточнения)

1. **Что именно подразумевается под "криптографией India"?**
   - [ ] Вариант A: CCA-совместимый профиль (RSA-2048+, ECDSA P-256/P-384, SHA-256+)
   - [ ] Вариант B: CERT-In рекомендации по TLS
   - [ ] Вариант C: Профиль для банковского сектора (RBI)
   - [ ] Вариант D: Другое (пожалуйста, уточните)

2. **Предпочтительные кривые**:
   - [ ] P-256 (128-bit security)
   - [ ] P-384 (192-bit security)
   - [ ] Обе

3. **Минимальная версия TLS**:
   - [ ] TLS 1.2 (для совместимости)
   - [ ] TLS 1.3 only (для безопасности)

4. **Нужны ли какие-то специфические cipher suites?**

## Possible Implementation Approach

Если требуется CCA-совместимый профиль, реализация будет похожа на `uk` provider:

```go
// crypto/in/provider.go
package in

type Provider struct{}

func (p *Provider) Name() string { return "in" }

func (p *Provider) SupportedCipherSuites() []uint16 {
    return []uint16{
        tls.TLS_AES_256_GCM_SHA384,       // TLS 1.3
        tls.TLS_AES_128_GCM_SHA256,       // TLS 1.3
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, // TLS 1.2
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,   // TLS 1.2
    }
}

func (p *Provider) ConfigureTLS(cfg *tls.Config) error {
    cfg.MinVersion = tls.VersionTLS12  // or TLS13
    cfg.CurvePreferences = []tls.CurveID{
        tls.CurveP384,
        tls.CurveP256,
    }
    // ...
}
```

## Constraints

- **Technical**: Должен использовать интерфейс `crypto.Provider`
- **Standards**: Следует индийским регуляторным требованиям
- **Pattern**: Следует существующей структуре провайдеров

## References

- [CCA Interoperability Guidelines](https://cca.gov.in/sites/files/pdf/guidelines/CCA-IOG.pdf)
- [India PKI Forum](https://www.indiapki.org)
- [CERT-In](https://cert-in.org.in)
- [Digital India - CCA](https://www.digitalindia.gov.in/di_ecosystem/controller-of-certifying-authorities-cca/)

---

## Approval

- [ ] Reviewed by:
- [ ] Approved on:
- [ ] Notes:
