package result_test

import (
	"path/filepath"
	"testing"

	"task228-seedgerm/internal/model"
	"task228-seedgerm/internal/result"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug06ResultVersionsRemainContiguousAndReplace(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "probe.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	trial, err := st.CreateTrial("P-06", "结果版本", "小麦", model.TrialReviewing)
	if err != nil {
		t.Fatal(err)
	}
	svc := result.New(st)
	first, err := svc.Draft(trial.ID, "v1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Publish(first.ID); err != nil {
		t.Fatal(err)
	}
	second, err := svc.Draft(trial.ID, "v2")
	if err != nil {
		t.Fatal(err)
	}
	if second.Version != 2 || second.PrevVersion != 1 {
		t.Fatalf("unexpected second version: %+v", second)
	}
	if _, err := svc.Publish(second.ID); err != nil {
		t.Fatal(err)
	}
	items, err := svc.List(trial.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].State != model.ResultSuperseded || items[1].State != model.ResultPublished {
		t.Fatalf("version replacement failed: %+v", items)
	}
}
