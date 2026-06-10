package cloudrunproxy

import "net/http"

type healthResponse struct {
	OK      bool   `json:"ok"`
	Service string `json:"service"`
}

// HandleHealth は `{"ok":true,"service":"<serviceName>"}` を返す health
// handler を生成する。path は `/health` を使うこと — `/healthz` は Google
// フロント (run.app) が外部 request をインターセプトして汎用 404 を返す。
func HandleHealth(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		WriteJSON(w, http.StatusOK, healthResponse{OK: true, Service: serviceName})
	}
}
