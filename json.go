package cloudrunproxy

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse は WriteJSONError が返す固定形式の error body。
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON は v を JSON で書き出す共通ヘルパ。
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteJSONError は固定文言 msg を `{"error": msg}` で返す。
//
// msg には "cloud run upstream error" のような固定文言だけを渡すこと。
// upstream エラーの詳細 (err.Error() 等) を msg に流し込むのは値漏れに
// なるので禁止 — 詳細は caller 側で log にだけ出す。consumer repo は
// この規約を regression test で固定する。
func WriteJSONError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}
