# go-cloudrun-proxy

ippoan org の「CF Worker → Cloud Run proxy → GCP API」2 段構成 service
([release-wave-gcp](https://github.com/ippoan/release-wave-gcp) /
[secrets-inventory-gcp](https://github.com/ippoan/secrets-inventory-gcp)) が
手コピーで共有していた bootstrap skeleton の集約 Go module。

設計の親 issue: [#1](https://github.com/ippoan/go-cloudrun-proxy/issues/1)
(Refs ippoan/claude-md#76)

## Install

```sh
go get github.com/ippoan/go-cloudrun-proxy
```

## API

```go
import cloudrunproxy "github.com/ippoan/go-cloudrun-proxy"
```

| symbol | 用途 |
|---|---|
| `MustEnv(key string) string` | boot 時の必須 env 検証。空なら `log.Fatalf` で落とす |
| `RequireAPIKey(headerName, key string, next http.Handler) http.Handler` | shared secret 認証 middleware。constant-time 比較、key 未設定なら fail-closed で全 401 |
| `WriteJSON(w http.ResponseWriter, status int, v any)` | JSON response helper |
| `WriteJSONError(w http.ResponseWriter, status int, msg string)` | `{"error": msg}` を返す。**msg は固定文言のみ** — upstream error 詳細は log にだけ出す (値漏れ防止) |
| `StatusFromGRPC(err error) int` | gRPC error code → HTTP status mapping。非 gRPC error は従来互換の 502 |
| `HandleHealth(serviceName string) http.HandlerFunc` | `{"ok":true,"service":"<name>"}` health handler |

## 使用例

```go
package main

import (
	"log"
	"net/http"

	cloudrunproxy "github.com/ippoan/go-cloudrun-proxy"
)

func main() {
	apiKey := cloudrunproxy.MustEnv("MY_PROXY_API_KEY")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", cloudrunproxy.HandleHealth("my-proxy"))
	mux.Handle("/do-something", cloudrunproxy.RequireAPIKey(
		"X-My-Proxy-API-Key", apiKey,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result, err := callUpstream(r.Context())
			if err != nil {
				// 詳細は log のみ。body は固定文言 + gRPC code 由来の status。
				log.Printf("upstream error: %v", err)
				cloudrunproxy.WriteJSONError(w,
					cloudrunproxy.StatusFromGRPC(err), "upstream error")
				return
			}
			cloudrunproxy.WriteJSON(w, http.StatusOK, result)
		}),
	))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## 規約 (consumer 側で緩めないこと)

- **upstream エラー詳細を response body に echo しない**。body は固定文言、
  詳細は log にだけ出す。consumer repo は regression test で固定する
- shared secret 認証は constant-time 比較 + key 未設定 fail-closed
- health check の path は `/health` を使う (`/healthz` は Google フロントが
  外部 request をインターセプトして汎用 404 を返す)

## 開発

```sh
go vet ./...
go test ./... -race
go build ./...
```

CI は [ippoan/ci-workflows](https://github.com/ippoan/ci-workflows) の
`go-ci.yml` caller。`coverage_100.toml` により登録 file の coverage 100% が
gate される。
