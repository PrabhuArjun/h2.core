# Fork changes

This repository is a **fork** of the official [VPNclient/h2.core](https://github.com/VPNclient/h2.core).

- Copyright © 2010–2025 NativeMind (copyright holder).
- Original license: NativeMind NONC — see [`LICENSE`](./LICENSE).
- All credit for the HTTPS VPN engine belongs to the original authors.

Changes are documented here per license §2.3 (Attribution / record changes in README).

## Branch: refactor/separation-cleanup

Goal: finish the h2.core ↔ libh2 separation started upstream on 2026-05-16
(`flows/sdd-h2-libh2-separation/`), which left wrapper code duplicated across
both repositories.

**Removed (duplicates that belong in the `libh2` repository):**

- `cgo/` — cgo wrapper. Byte-identical to `libh2/wrappers/cgo/` (`client.go` 5866 B, `h2core.go` 5932 B).
- `mobile/` — gomobile wrapper. Byte-identical to `libh2/mobile/` (`client.go`, `doc.go`, `socks.go`).

Rationale: `libh2/wrappers/cgo/go.mod` already imports
`github.com/vpnclient/https-vpn`, i.e. `libh2` is intended to wrap `h2.core`.
The copies under `h2.core/cgo/` and `h2.core/mobile/` are leftovers from the
incomplete migration and are not imported anywhere inside the `https-vpn`
module (verified via grep for `https-vpn/cgo` and `https-vpn/mobile`).

Result: `h2.core` keeps the pure Go core + CLI (`core/`, `crypto/`,
`transport/`, `cmd/https-vpn/`); all platform wrappers live in `libh2`.

No core logic changed. No functional code modified — only removal of
duplicated wrapper directories.

## Branch: feature/gostvpn-mvp

Goal: некоммерческий MVP «ГОСТ VPN» на движке h2.core (сессия
`gostvpn-noncommercial-mvp`, 2026-05-20). NONC разрешает некоммерческое
использование без отдельной лицензии — бесплатная раздача легальна сейчас.

**Added:**

- `docs/gostvpn-mvp-setup.md` — runbook сборки/конфигурации/smoke + результаты
  локального smoke-теста + список известных ограничений ГОСТ-крипто.

**No core logic changed.** Только добавлена документация. Найденные при smoke
ограничения движка (загрузка GOST-cert PEM, привязка ГОСТ cipher к серту) —
зафиксированы в docs как апстрим-задачи, код движка не правился в этой ветке.

Smoke-результат: транспорт HTTP/2 CONNECT поверх TLS 1.3 работает (ALPN h2,
curl-туннель 200 OK). ГОСТ end-to-end требует доработки crypto/ru/tls
(GOST-cert PEM + cipher binding) — задокументировано.
