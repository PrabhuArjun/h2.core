# ГОСТ VPN — MVP setup (некоммерческий)

> Развёртывание бесплатного некоммерческого VPN на движке **h2.core**
> (HTTPS VPN + ГОСТ-криптография). Сессия `gostvpn-noncommercial-mvp`,
> ветка `feature/gostvpn-mvp`, 2026-05-20.
>
> **Лицензия.** h2.core под NativeMind NONC. Некоммерческое использование
> (бесплатная раздача) разрешено без отдельной лицензии. **Никаких платных
> фич** — коммерция заблокирована до письменной лицензии NativeMind. Этот
> MVP строго некоммерческий. См. `FORK_CHANGES.md` (§2.3 атрибуция).

---

## Что это

h2.core — «HTTPS VPN»: туннель — это стандартный **HTTP/2 CONNECT поверх
TLS 1.3**, статистически неотличимый от браузерного HTTPS. DPI/ТСПУ не
выделяет его как VPN — это и есть защита от блокировок (в отличие от
VLESS-Reality, который блокируется в РФ-регионах с 17.02.2026, см.
`shantinet/docs/rkn-baseline-2026-05-19.md`).

Плюс — подключаемая национальная криптография, включая **ГОСТ R 34.10/12/15**
(`crypto/ru/`).

---

## Сборка

Требования: Go 1.25+ (тестировалось на 1.26.3 darwin/arm64).

### Локально (для разработки и smoke)

```bash
cd ~/Projects/h2.core
CGO_ENABLED=0 go build -tags "bundle_custom crypto_ru" -o /tmp/gostvpn ./cmd/https-vpn
```

Tags:
- `bundle_custom` — включить кастомный crypto-бандл (не дефолтный all)
- `crypto_ru` — включить ГОСТ-провайдер (`crypto/ru`)

### Кросс-компиляция для сервера (linux)

```bash
make build-unix
# → dist/https_vpn (darwin host), dist/https_vpn_linux_amd64, _linux_arm64
# build/unix.sh использует CGO_ENABLED=0, статичные бинари
```

Для ГОСТ-бандла на linux:
```bash
TAGS="bundle_custom crypto_ru" make build-unix
```

---

## Конфигурация

`https-vpn init -crypto ru` создаёт `config.json`:

```json
{
  "inbounds": [{
    "port": 443,
    "protocol": "https-vpn",
    "streamSettings": {
      "network": "h2",
      "security": "tls",
      "tlsSettings": {
        "serverName": "<your-domain>",
        "certificates": [
          { "certificateFile": "cert.pem", "keyFile": "key.pem" }
        ],
        "cipherSuites": "ru",
        "minVersion": "TLS1.3"
      }
    }
  }],
  "outbounds": [{ "protocol": "freedom" }]
}
```

> **NB:** поле `cryptoProvider` устарело — движок выводит warning и просит
> использовать `cipherSuites`. На момент 2026-05-20 (ветка
> `refactor/separation-cleanup` + `feature/gostvpn-mvp`) оба поля читаются,
> но канон — `cipherSuites`.

Запуск:
```bash
https-vpn run -c config.json
```

---

## Smoke-тест (локальный, проверено 2026-05-20)

```bash
# 1. Сгенерировать сертификат
#    ГОСТ self-signed (см. ограничение ниже):
(cd ~/Projects/h2.core && go run ./crypto/ru/tools/gencert.go)  # → gost_cert.pem, gost_key.pem
#    ИЛИ обычный (для проверки транспорта):
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 7 -nodes -subj "/CN=localhost"

# 2. config.json с portом 8443 (443 требует root), TLS1.3, cipherSuites:ru

# 3. Запустить сервер
/tmp/gostvpn run -c config.json &

# 4. Проверить TLS + ALPN
echo | openssl s_client -connect localhost:8443 -tls1_3 -alpn h2 2>&1 | grep -iE "protocol|cipher|ALPN"
#   → Protocol: TLSv1.3, ALPN protocol: h2  ✅

# 5. Проверить HTTP/2 CONNECT туннель
curl -sS -x https://localhost:8443 --proxy-insecure -o /dev/null \
  -w "via-proxy HTTP %{http_code}, %{time_total}s\n" https://example.com
#   → via-proxy HTTP 200  ✅ — туннель проксирует реальный трафик
```

### Результаты smoke (2026-05-20)

| Проверка | Результат |
|---|---|
| Сборка `bundle_custom crypto_ru` | ✅ бинарь 9 MB |
| TLS 1.3 handshake | ✅ |
| ALPN `h2` (HTTP/2) | ✅ — выглядит как браузерный HTTPS |
| HTTP/2 CONNECT туннель | ✅ curl → example.com → 200 за 0.94s |

---

## ⚠️ Известные ограничения (на 2026-05-20) — что доделать до prod

Smoke выявил, что **ГОСТ-крипто end-to-end ещё не завершён** в форке:

1. **Загрузка ГОСТ-сертификата не работает.** `gencert.go` создаёт PEM-блок
   `GOST CERTIFICATE`, но загрузчик сертификатов использует стандартный Go
   `crypto/tls`, который принимает только блок `CERTIFICATE`:
   ```
   Failed to start: failed to load certificates: tls: failed to find
   "CERTIFICATE" PEM block ... after skipping: [GOST CERTIFICATE]
   ```
   Сам gencert.go помечает вывод как «Mock» и в комментарии:
   «In a real implementation, we would use PKCS#8 with GOST OIDs».

2. **При обычном (RSA) серте провайдер откатывается на `us`.** Лог запуска
   с `cipherSuites:ru` + RSA-сертом:
   ```
   Crypto provider: ru
   Certificate loaded: provider=us file=...
   CertificateStore: providers=[us], priority=[us]
   ```
   То есть фактический cipher — `TLS_AES_128_GCM_SHA256` (NIST/AES), не ГОСТ.

**Вывод:** транспортный слой (HTTP/2 CONNECT, DPI-устойчивость) — **работает
и готов**. ГОСТ-криптография как сквозной TLS — **требует доработки**:
поддержка GOST-cert PEM (PKCS#8 + GOST OIDs) в загрузчике + привязка
ГОСТ cipher suites к ГОСТ-серту. Это апстрим-задача (вероятно для самого
NativeMind / Ананты) либо отдельная сессия по `crypto/ru/tls`.

---

## Удалённый деплой (следующий шаг, требует тестовый сервер)

Не выполнен в этой итерации (выбран локальный smoke). Когда будет тестовый
сервер (НЕ боевая инфра ShantiNet):

1. `TAGS="bundle_custom crypto_ru" make build-unix` → `dist/https_vpn_linux_amd64`
2. `scp dist/https_vpn_linux_amd64 root@<TEST_IP>:/usr/local/bin/gostvpn`
3. Реальный TLS-серт (Let's Encrypt) для домена → `config.json` порт 443
4. systemd-unit + запуск
5. Клиент: CLI h2.core или обёртка из `libh2` (сессия `vpnclient-libh2-reorg`)

---

## Юридический чеклист (NONC, соблюдено)

- ✅ §2.1 — работаем в публичном форке `PrabhuArjun/h2.core`
- ✅ §2.3 — атрибуция © NativeMind + ссылка на оригинал + лог в `FORK_CHANGES.md`
- ✅ §2.8 ShareAlike — код ГОСТ VPN остаётся open-source под NONC
- ✅ §2.5 — своё имя «ГОСТ VPN», без копирования интерфейса оригинала
- ✅ §2.4 — никакого вредоносного кода
- ✅ Некоммерческий — без платных фич / биллинга
- ⛔ §2.6/§2.7 — коммерция (платные подписки) до письменной лицензии NativeMind НЕ делается
