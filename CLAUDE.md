# CLAUDE.md

Claude Code 向けの本リポジトリ作業ルール。

設計の親 issue: [#1](https://github.com/ippoan/go-cloudrun-proxy/issues/1)
(Refs ippoan/claude-md#76 の org 横断重複監査が起点)

## このリポジトリの方針

- ippoan org の「CF Worker → Cloud Run proxy → GCP API」2 段構成 service
  (`ippoan/release-wave-gcp` / `ippoan/secrets-inventory-gcp`) が手コピーで
  共有していた bootstrap skeleton を集約した **共有 Go module**
  (`github.com/ippoan/go-cloudrun-proxy`)
- **薄く保つ**。export するのは skeleton (env / auth / JSON helper /
  gRPC status mapping / health) だけ。GCP SDK client や token source 層は
  scope 外 (= consumer repo 側の責任)
- **upstream エラー詳細を response body に echo しない** (固定文言 + log のみ)
  が consumer 共通の確立規約 — lib 化で緩めない。test で固定する
- shared secret 認証は **constant-time 比較 + key 未設定 fail-closed**
  (`RequireAPIKey`)。header 名・JSON body 形式は repo ごとに差があるため
  引数化する (強制統一しない)

## 公開 API (SemVer 対象)

| symbol | 用途 |
|---|---|
| `MustEnv(key string) string` | boot 時の必須 env 検証 (空なら log.Fatalf) |
| `RequireAPIKey(headerName, key string, next http.Handler) http.Handler` | constant-time な shared secret 認証 middleware (key 未設定 fail-closed) |
| `WriteJSON(w, status, v)` / `WriteJSONError(w, status, msg)` | JSON response helper。`WriteJSONError` の msg は固定文言のみ |
| `StatusFromGRPC(err error) int` | gRPC code → HTTP status mapping (非 gRPC error は 502) |
| `HandleHealth(serviceName string) http.HandlerFunc` | `{"ok":true,"service":...}` health handler |

## Worktree / branch 命名規則

形式: `<issue-number>-<type>-<short-description>`

- `issue-number`: 必須。先に issue を立ててから worktree / branch を作る
- `type`: `feat` | `fix` | `refactor` | `infra`
- `short-description`: 半角小文字英数字とハイフン

Claude Code が自動採番する `claude/...` で実装に入る場合は、対応する issue を
紐付けた上で PR description に `Refs #N` を明記する。

## PR description / commit message のキーワード

- 使用禁止: `Closes #N` / `Fixes #N` / `Resolves #N`
  - PR auto-merge が走った瞬間に issue が自動 close されるため、release 時の
    close 確認 UI と整合しない
- 使用推奨: `Refs #N` / `Related to #N` / `Part of #N`

PR テンプレートは `.github/pull_request_template.md` で `Refs` を強制する。

## ビルド / テスト

PR を出す前に手元で green に:

```sh
go vet ./...
go test ./... -race
go build ./...
```

CI (`.github/workflows/ci.yml`) は ci-workflows の `go-ci.yml` caller。
`coverage_100.toml` を repo root に置いてあるため **登録 file の coverage
100% gate が常時有効** — function を追加したらテストも同 PR で書くこと。
`MustEnv` の fail 経路のような untestable に見える path は `fatalf` package
variable のような注入点を作ってテストする (exclude_funcs に逃がす前に検討)。

## GitHub 自動化 (重要)

- **`main` に直 push しない。** PR を作る。
- PR / commit は `Refs #N` を使う (`Closes/Fixes/Resolves` は禁止 — auto-close 防止)。
- `mcp__github__enable_pr_auto_merge` を reflex で呼ばない (user 明示指示時のみ)。
- PR 作成後は同じ turn で `mcp__github__subscribe_pr_activity` を呼び CI を watch する。

## consumer 移行 (親 issue #1 の checklist)

公開 API を変える時は consumer の利用箇所を必ず確認:

- `ippoan/release-wave-gcp` — `StatusFromGRPC` 採用で blanket 502 改善も合わせる
- `ippoan/secrets-inventory-gcp` — bare `http.Error` → `WriteJSONError` 統一は repo 側判断

---

_共通項を直すときは [`ippoan/claude-md`](https://github.com/ippoan/claude-md) の
`CLAUDE.md.template` を更新すること。_
