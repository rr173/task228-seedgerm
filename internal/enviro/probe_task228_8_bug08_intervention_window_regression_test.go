package enviro_test

import (
	"path/filepath"
	"testing"
	"time"

	"task228-seedgerm/internal/enviro"
	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug08InterventionWindowDetectsCoolingRegardlessOfRowOrder(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	trial, err := st.CreateTrial("P-08", "干预窗口", "水稻", model.TrialRunning)
	if err != nil {
		t.Fatal(err)
	}
	svc := enviro.New(st)
	at := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	for _, item := range []struct {
		d time.Duration
		t float64
	}{
		{-time.Hour, 25},
		{time.Hour, 21},
	} {
		if _, err := svc.Record(enviro.SampleInput{TrialID: trial.ID, SampledAt: at.Add(item.d), TempC: item.t, Humidity: 60, Instrument: "thermo-1"}); err != nil {
			t.Fatal(err)
		}
	}
	iv, err := svc.InterveneAround(trial.ID, at, 2*time.Hour, 2*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if !iv.Cooling || iv.TempDrop != 4 {
		t.Fatalf("expected cooling summary, got %+v", iv)
	}
}
