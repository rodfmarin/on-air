package lifxutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("test-token")
	if c.Token != "test-token" {
		t.Errorf("expected token 'test-token', got %v", c.Token)
	}
}

func TestLightStruct(t *testing.T) {
	light := Light{ID: "id", Label: "label", Power: "on", Color: "red", Brightness: 0.5}
	if light.ID != "id" || light.Label != "label" || light.Power != "on" || light.Color != "red" || light.Brightness != 0.5 {
		t.Errorf("Light struct fields not set correctly: %+v", light)
	}
}

func TestListLights_Success(t *testing.T) {
	lights := []Light{{ID: "id1", Label: "label1"}, {ID: "id2", Label: "label2"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/lights/all" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lights); err != nil {
			t.Errorf("json encode error: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	got, err := c.ListLights()
	if err != nil {
		t.Fatalf("ListLights failed: %v", err)
	}
	if !reflect.DeepEqual(got, lights) {
		t.Errorf("ListLights: got %+v, want %+v", got, lights)
	}
}

func TestListLights_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("error")); err != nil {
			t.Errorf("write error: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	_, err := c.ListLights()
	if err == nil {
		t.Error("expected error for ListLights, got nil")
	}
}

func TestTogglePower_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	if err := c.TogglePower("id:test"); err != nil {
		t.Errorf("TogglePower failed: %v", err)
	}
}

func TestTogglePower_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("error")); err != nil {
			t.Errorf("write error: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	if err := c.TogglePower("id:test"); err == nil {
		t.Error("expected error for TogglePower, got nil")
	}
}

func TestSetState_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var state map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			t.Errorf("json decode error: %v", err)
		}
		if state["power"] != "on" || state["color"] != "red" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	state := map[string]interface{}{"power": "on", "color": "red"}
	if err := c.SetState("id:test", state); err != nil {
		t.Errorf("SetState failed: %v", err)
	}
}

func TestSetState_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("error")); err != nil {
			t.Errorf("write error: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	state := map[string]interface{}{"power": "on", "color": "red"}
	if err := c.SetState("id:test", state); err == nil {
		t.Error("expected error for SetState, got nil")
	}
}

func TestSetBusy_FallbackColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var state map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			t.Errorf("json decode error: %v", err)
		}
		if state["color"] != "red saturation:0.5" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	light := Light{ID: "test"}
	if err := c.SetBusy(light, ""); err != nil {
		t.Errorf("SetBusy fallback color failed: %v", err)
	}
}

func TestSetFree_FallbackColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var state map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
			t.Errorf("json decode error: %v", err)
		}
		if state["color"] != "kelvin:2671" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Client{Token: "test-token", BaseURL: server.URL + "/v1/"}

	light := Light{ID: "test"}
	if err := c.SetFree(light, ""); err != nil {
		t.Errorf("SetFree fallback color failed: %v", err)
	}
}
