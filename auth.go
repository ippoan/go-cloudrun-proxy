package cloudrunproxy

import (
	"crypto/subtle"
	"net/http"
)

// RequireAPIKey は headerName の shared secret を constant-time 比較する
// 認証 middleware。header 名は repo ごとに異なる (X-Release-Wave-API-Key /
// X-Inventory-API-Key 等) ため引数で受ける。
//
// expected が空 (= secret 未注入) の場合は fail-closed で全 request を
// 401 reject する。subtle.ConstantTimeCompare は空 vs 空で一致 (=1) を
// 返すため、空チェックを比較より先に行う。
func RequireAPIKey(headerName, expected string, next http.Handler) http.Handler {
	expectedBytes := []byte(expected)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(headerName)
		if expected == "" ||
			subtle.ConstantTimeCompare([]byte(got), expectedBytes) != 1 {
			WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
