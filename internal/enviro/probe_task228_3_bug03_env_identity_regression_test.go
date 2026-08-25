package enviro_test

import (
	"path/filepath"
	"testing"
	"time"

	"task228-seedgerm/internal/enviro"
	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug03EnvironmentIdentityUsesStableSampleKey(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	trial, err := st.CreateTrial("P-03", "环境幂等", "水稻", model.TrialRunning)
	if err != nil {
		t.Fatal(err)
	}
	svc := enviro.New(st)
	at := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	if _, err := svc.Record(enviro.SampleInput{TrialID: trial.ID, SampledAt: at, TempC: 25, Humidity: 60, Instrument: "thermo-1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Record(enviro.SampleInput{TrialID: trial.ID, SampledAt: at, TempC: 21, Humidity: 65, Instrument: "thermo-1"}); err != nil {
		t.Fatal(err)
	}
	items, err := st.ListEnv(trial.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].TempC != 21 {
		t.Fatalf("same instrument/time should update one sample, got %+v", items)
	}
}
