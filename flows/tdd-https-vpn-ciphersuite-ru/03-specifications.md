# Спецификации: tdd-https-vpn-ciphersuite-ru

> Версия: 1.0  
> Статус: DRAFT  
> Последнее обновление: 2026-05-11  
> Требования: [01-requirements.md](01-requirements.md)  
> Тесты: [02-tests.md](02-tests.md)

## Архитектура (на основе тестов)

Для удовлетворения требований по регистрации (Case 1.1) и настройке TLS (Case 2.1, 2.2), система разделена на слои:

### 1. Слой провайдера (`crypto/ru`)
- **Компонент**: `Provider`
- **Интерфейс**: `crypto.Provider`
- **Функция**: Точка входа для регистрации "ru" и конфигурации TLS.

### 2. Слой TLS (`crypto/ru/tls`)
- **Компонент**: `GostTLS`
- **Функция**: Реализация специфичных для ГОСТ протокольных обменов (VKO), не поддерживаемых стандартным `crypto/tls`.

### 3. Слой примитивов (`crypto/ru/gost`)
- **Компоненты**: `Kuznyechik`, `Magma`, `Streebog`, `Gost3410`
- **Функция**: Математическая реализация алгоритмов (удовлетворяет Case 3.1, 3.2).

## Проектирование интерфейсов

### Провайдер ГОСТ
```go
type Provider struct{}

func (p *Provider) Name() string { return "ru" }
func (p *Provider) ConfigureTLS(cfg *tls.Config) error
func (p *Provider) SupportedCipherSuites() []uint16
```

## Модели данных

### Идентификаторы наборов шифров (RFC 9189)
```go
const (
    TLS_GOSTR341112_256_WITH_KUZNYECHIK_CTR_OMAC uint16 = 0xC100
    TLS_GOSTR341112_256_WITH_KUZNYECHIK_MGM_L    uint16 = 0xC101
    TLS_GOSTR341112_256_WITH_MAGMA_CTR_OMAC      uint16 = 0xC103
)
```

## Стратегия тестирования

### Модульные тесты (Unit)
- `crypto/ru/gost/streebog_test.go`: Тестирование хеша на векторах.
- `crypto/ru/gost/kuznyechik_test.go`: Тестирование блочного шифра.

### Интеграционные тесты
- `crypto/ru/provider_test.go`: Проверка корректности заполнения `tls.Config`.
- Эмуляция Handshake с использованием ГОСТ сертификатов для проверки VKO.
