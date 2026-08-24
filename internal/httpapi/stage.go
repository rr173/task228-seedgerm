package httpapi

import (
	"encoding/json"
	"net/http"

	"task228-seedgerm/internal/model"
)

// handleStages 处理 /api/stages/:id 及其子资源（confirm / resolve）。
func (s *Server) handleStages(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.URL.Path, "/api/stages/")
	if !ok {
		writeError(w, http.StatusBadRequest, model.ErrBadInput)
		return
	}
	rem := remaining(r.URL.Path, "/api/stages/")
	switch {
	case rem == "":
		s.stageDetail(w, r, id)
	case rem == "confirm":
		s.stageConfirm(w, r, id)
	case rem == "resolve":
		s.stageResolve(w, r, id)
	default:
		writeError(w, http.StatusNotFound, model.ErrNotFound)
	}
}

// stageDetail GET 阶段事件详情。
func (s *Server) stageDetail(w http.ResponseWriter, r *http.Request, id int64) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, model.ErrBadInput)
		return
	}
	ev, err := s.store.GetStage(id)
	if err != nil {
		s.mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ev)
}

// stageConfirm POST 确认阶段事件。
func (s *Server) stageConfirm(w http.ResponseWriter, r *http.Request, id int64) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, model.ErrBadInput)
		return
	}
	ev, err := s.svc.ConfirmStage(id)
	if err != nil {
		s.mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ev)
}

// stageResolve POST 处理冲突：revoke 冲突事件并 confirm 偏好事件。
func (s *Server) stageResolve(w http.ResponseWriter, r *http.Request, conflictID int64) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, model.ErrBadInput)
		return
	}
	var body struct {
		PreferID int64 `json:"prefer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, model.ErrBadInput)
		return
	}
	seedID, err := s.stageSeedID(conflictID)
	if err != nil {
		s.mapErr(w, err)
		return
	}
	if _, err := s.svc.Review.RevokeStage(conflictID); err != nil {
		s.mapErr(w, err)
		return
	}
	if body.PreferID != conflictID {
		if _, err := s.svc.ConfirmStage(body.PreferID); err != nil {
			s.mapErr(w, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]int64{"seed_id": seedID})
}

// stageSeedID 返回阶段事件所属种子。
func (s *Server) stageSeedID(stageID int64) (int64, error) {
	ev, err := s.store.GetStage(stageID)
	if err != nil {
		return 0, err
	}
	return ev.SeedID, nil
}
