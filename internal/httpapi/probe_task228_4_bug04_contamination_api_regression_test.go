package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"task228-seedgerm/internal/httpapi"
	"task228-seedgerm/internal/service"
	"task228-seedgerm/internal/store"
)

func TestTask228Bug04ContaminationScoreReachesSeedState(t *testing.T) {
	st, err := store.Open(t.TempDir() + "/probe.db")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	server := httptest.NewServer(httpapi.New(service.New(st), st, ":0", "").Handler())
	defer server.Close()
	post := func(path, body string) map[string]interface{} {
		resp, err := http.Post(server.URL+path, "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			t.Fatalf("POST %s returned %d", path, resp.StatusCode)
		}
		var value map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
			t.Fatal(err)
		}
		return value
	}
	trial := post("/api/trials", `{"code":"P-04","name":"污染识别","species":"玉米"}`)
	trialID := int64(trial["id"].(float64))
	trialPath := strconv.FormatInt(trialID, 10)
	post("/api/trials/"+trialPath+"/transition", `{"to":"running"}`)
	seed := post("/api/trials/"+trialPath+"/seeds", `{"seed_no":"S1"}`)
	seedID := int64(seed["id"].(float64))
	seedPath := strconv.FormatInt(seedID, 10)
	event := post("/api/seeds/"+seedPath+"/detect", `{"contam_score":0.9}`)
	eventPath := strconv.FormatInt(int64(event["id"].(float64)), 10)
	post("/api/stages/"+eventPath+"/confirm", `{}`)
	resp, err := http.Get(server.URL + "/api/seeds/" + seedPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got struct {
		State string `json:"state"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.State != "contaminated" {
		t.Fatalf("contamination score did not mark seed: %q", got.State)
	}
}
