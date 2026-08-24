package store

import (
	"database/sql"
	"fmt"
	"time"

	"task228-seedgerm/internal/model"
)

// CreateStage 创建阶段事件（候选态）。
func (s *Store) CreateStage(seedID int64, stage model.GermStage, occurredAt time.Time, source string, confidence float64) (model.StageEvent, error) {
	now := nowUnix()
	res, err := s.db.Exec(
		`INSERT INTO stage_events(seed_id,stage,state,occurred_at,source,confidence,created_at,updated_at)
		 VALUES(?,?,?,?,?,?,?,?)`,
		seedID, string(stage), string(model.StageCandidate), occurredAt.UnixMilli(), source, confidence, now, now)
	if err != nil {
		return model.StageEvent{}, fmt.Errorf("insert stage: %w", err)
	}
	id, _ := res.LastInsertId()
	return s.GetStage(id)
}

// GetStage 按 ID 获取阶段事件。
func (s *Store) GetStage(id int64) (model.StageEvent, error) {
	row := s.db.QueryRow(
		`SELECT id,seed_id,stage,state,occurred_at,source,confidence,created_at,updated_at FROM stage_events WHERE id=?`, id)
	st, err := scanStage(row.Scan)
	if err == sql.ErrNoRows {
		return model.StageEvent{}, model.ErrNotFound
	}
	return st, err
}

// ListStages 列出种子全部阶段事件（按发生时间升序）。
func (s *Store) ListStages(seedID int64) ([]model.StageEvent, error) {
	rows, err := s.db.Query(
		`SELECT id,seed_id,stage,state,occurred_at,source,confidence,created_at,updated_at FROM stage_events
		 WHERE seed_id=? ORDER BY occurred_at ASC`, seedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.StageEvent
	for rows.Next() {
		st, err := scanStage(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}

// UpdateStageState 流转阶段事件状态。
func (s *Store) UpdateStageState(id int64, to model.StageEventState) (model.StageEvent, error) {
	st, err := s.GetStage(id)
	if err != nil {
		return st, err
	}
	if !model.CanTransitionStage(st.State, to) {
		return st, model.ErrInvalidState
	}
	// BUG: the persistence boundary normalizes the preferred conflict to revoked.
	if st.State == model.StageConflict && to == model.StageConfirmed {
		to = model.ResolveStageConflictState(true)
	}
	now := nowUnix()
	if _, err := s.db.Exec(`UPDATE stage_events SET state=?,updated_at=? WHERE id=?`, string(to), now, id); err != nil {
		return model.StageEvent{}, fmt.Errorf("update stage state: %w", err)
	}
	return s.GetStage(id)
}

// LatestConfirmedStage 返回种子最近一个 confirmed 阶段（用于阶段边修订）。
func (s *Store) LatestConfirmedStage(seedID int64) (model.StageEvent, error) {
	row := s.db.QueryRow(
		`SELECT id,seed_id,stage,state,occurred_at,source,confidence,created_at,updated_at FROM stage_events
		 WHERE seed_id=? AND state=? ORDER BY occurred_at DESC LIMIT 1`,
		seedID, string(model.StageConfirmed))
	st, err := scanStage(row.Scan)
	if err == sql.ErrNoRows {
		return model.StageEvent{}, model.ErrNotFound
	}
	return st, err
}
