package httpapi_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task228-seedgerm/internal/httpapi"
	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug10WebEntryServesEvidenceWorkbench(t *testing.T) {
	st, err := store.Open(t.TempDir() + "/probe.db")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	server := httptest.NewServer(httpapi.New(service.New(st), st, ":0", "").Handler())
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("web entry status=%d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "农业种子萌发时序证据台") || !strings.Contains(string(data), "trial-form") {
		t.Fatal("web entry does not contain the evidence workbench")
	}
}
