package service_test

import (
	"path/filepath"
	"testing"
	"time"

	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug02RejectsOutOfOrderImage(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	trial, _ := svc.CreateTrial("P-02", "时间顺序", "小麦")
	_, _ = svc.TransitionTrial(trial.ID, model.TrialRunning)
	seed, err := svc.CreateSeed(trial.ID, "S1")
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 8, 24, 9, 0, 0, 0, time.UTC)
	if _, _, err := svc.IngestImage(seed.ID, "later", at, 640, 480, ""); err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.IngestImage(seed.ID, "earlier", at.Add(-time.Minute), 640, 480, ""); err != model.ErrTimeReversed {
		t.Fatalf("older frame must be rejected, got %v", err)
	}
}
