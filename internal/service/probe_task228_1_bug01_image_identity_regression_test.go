package service_test

import (
	"path/filepath"
	"testing"
	"time"

	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug01ImageIdentityIgnoresNote(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := service.New(st)
	trial, err := svc.CreateTrial("P-01", "图像幂等", "小麦")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.TransitionTrial(trial.ID, model.TrialRunning); err != nil {
		t.Fatal(err)
	}
	seed, err := svc.CreateSeed(trial.ID, "S1")
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 8, 24, 8, 0, 0, 0, time.UTC)
	if _, _, err := svc.IngestImage(seed.ID, "same-frame", at, 640, 480, "first annotation"); err != nil {
		t.Fatal(err)
	}
	_, created, err := svc.IngestImage(seed.ID, "same-frame", at.Add(time.Hour), 640, 480, "revised annotation")
	if err != nil {
		t.Fatal(err)
	}
	if created {
		t.Fatal("same image hash must remain idempotent when the note changes")
	}
	images, err := st.ListImages(seed.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(images) != 1 {
		t.Fatalf("expected one evidence frame, got %d", len(images))
	}
}
