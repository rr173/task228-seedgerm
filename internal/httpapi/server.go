// Package httpapi HTTP 层：以 /api 前缀暴露试验/种子/图像/环境/阶段/复核/结果/自检 API。
package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

// Server HTTP 服务。
type Server struct {
	svc    *service.Service
	store  *store.Store
	addr   string
	dbPath string
}

// New 构造 HTTP 服务。
func New(svc *service.Service, st *store.Store, addr, dbPath string) *Server {
	return &Server{svc: svc, store: st, addr: addr, dbPath: dbPath}
}

// Handler 返回路由 mux。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/trials", s.handleTrials)
	mux.HandleFunc("/api/trials/", s.handleTrialByID)
	mux.HandleFunc("/api/seeds/", s.handleSeeds)       // 含 /api/seeds/:id/...
	mux.HandleFunc("/api/stages/", s.handleStages)     // 含 /api/stages/:id/...
	mux.HandleFunc("/api/results/", s.handleResults)   // 含 /api/results/:id/...
	mux.HandleFunc("/api/selfcheck", s.handleSelfCheck)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

// Start 启动服务（长驻）。
func (s *Server) Start() error {
	srv := &http.Server{Addr: s.addr, Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second}
	return srv.ListenAndServe()
}

// writeJSON 写 JSON 响应。
func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError 写错误响应。
func writeError(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}

// parseID 从路径段取 ID。
func parseID(path, prefix string) (int64, bool) {
	rest := strings.TrimPrefix(path, prefix)
	rest = strings.Trim(rest, "/")
	if rest == "" {
		return 0, false
	}
	// 取第一段
	seg := rest
	if idx := strings.Index(rest, "/"); idx >= 0 {
		seg = rest[:idx]
	}
	id, err := strconv.ParseInt(seg, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

// remaining 返回路径在 id 之后的剩余段（不含前导斜杠）。
func remaining(path, prefix string) string {
	rest := strings.TrimPrefix(path, prefix)
	rest = strings.Trim(rest, "/")
	if idx := strings.Index(rest, "/"); idx >= 0 {
		return strings.Trim(rest[idx+1:], "/")
	}
	return ""
}
