package stage_test

import (
	"path/filepath"
	"testing"
	"time"

	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/stage"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug07ConflictResolutionRevokesOnlyRejectedStage(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	trial, _ := st.CreateTrial("P-07", "冲突复核", "玉米", model.TrialRunning)
	seed, err := st.CreateSeed(trial.ID, "S1", model.SeedPending)
	if err != nil {
		t.Fatal(err)
	}
	detector := stage.New(st)
	first, err := st.CreateStage(seed.ID, model.StageRadicle, time.Now().UTC(), "auto", .9)
	if err != nil {
		t.Fatal(err)
	}
	second, err := st.CreateStage(seed.ID, model.StageColeoptile, time.Now().UTC().Add(time.Minute), "manual", .8)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpdateStageState(first.ID, model.StageConflict); err != nil {
		t.Fatal(err)
	}
	if err := detector.ResolveConflict(first.ID, second.ID); err != nil {
		t.Fatal(err)
	}
	gotFirst, _ := st.GetStage(first.ID)
	gotSecond, _ := st.GetStage(second.ID)
	if gotFirst.State != model.StageRevoked || gotSecond.State != model.StageConfirmed {
		t.Fatalf("unexpected resolution: first=%s second=%s", gotFirst.State, gotSecond.State)
	}
}
