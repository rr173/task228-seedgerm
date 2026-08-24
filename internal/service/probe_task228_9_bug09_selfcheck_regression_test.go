package service_test

import (
	"path/filepath"
	"testing"
	"time"

	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug09SelfCheckReportsAllEvidenceTables(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	trial, _ := svc.CreateTrial("P-09", "自检计数", "小麦")
	_, _ = svc.TransitionTrial(trial.ID, model.TrialRunning)
	seed, err := svc.CreateSeed(trial.ID, "S1")
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.IngestImage(seed.ID, "frame", time.Now().UTC(), 100, 100, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RecordEnv(trial.ID, time.Now().UTC(), 22, 60, "thermo-1"); err != nil {
		t.Fatal(err)
	}
	report, err := svc.SelfCheck()
	if err != nil {
		t.Fatal(err)
	}
	if report.Trials != 1 || report.Seeds != 1 || report.Images != 1 || report.EnvSamples != 1 {
		t.Fatalf("self-check omitted persisted evidence: %+v", report)
	}
}
