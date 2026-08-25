package service_test

import (
	"path/filepath"
	"testing"

	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug05SealedTrialCannotReopen(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	trial, err := svc.CreateTrial("P-05", "封存状态", "大豆")
	if err != nil {
		t.Fatal(err)
	}
	for _, state := range []model.TrialState{model.TrialRunning, model.TrialReviewing, model.TrialCompleted, model.TrialSealed} {
		if _, err := svc.TransitionTrial(trial.ID, state); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.TransitionTrial(trial.ID, model.TrialRunning); err != model.ErrInvalidState {
		t.Fatalf("sealed trial must reject reopening, got %v", err)
	}
}
