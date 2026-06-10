// Package cloudrunproxy は ippoan org の「CF Worker → Cloud Run proxy → GCP API」
// 2 段構成 service (release-wave-gcp / secrets-inventory-gcp 等) が手コピーで
// 共有していた bootstrap skeleton の集約 lib。
//
// 規約 (consumer 側で緩めないこと):
//   - upstream エラー詳細を response body に echo しない。body は固定文言、
//     詳細は log にだけ出す (WriteJSONError の doc 参照)
//   - shared secret 認証は constant-time 比較 + key 未設定 fail-closed
//     (RequireAPIKey 参照)
//
// Refs https://github.com/ippoan/go-cloudrun-proxy/issues/1
package cloudrunproxy

import (
	"log"
	"os"
)

// fatalf は MustEnv の失敗経路。テストから差し替えて process exit なしで
// 検証できるよう variable にしている (本番は log.Fatalf = exit(1))。
var fatalf = log.Fatalf

// MustEnv は環境変数 key を読み、空なら log.Fatalf で process を落とす。
// boot 時の必須 env 検証用 (optional な env には使わないこと)。
func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fatalf("env %s is required", key)
	}
	return v
}
