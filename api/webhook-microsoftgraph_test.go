package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookMicrosoftGraph_ValidationToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/webhook?validationToken=testToken", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook_MicrosoftGraph)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "testToken"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestWebhookMicrosoftGraph_Callback(t *testing.T) {
	event := WebhookEventStruct{
		SubscriptionID:                 "sub123",
		SubscriptionExpirationDateTime: time.Now(),
		ChangeType:                     "updated",
		Resource:                       "resource123",
		ResourceData: struct {
			OdataType string `json:"@odata.type"`
			OdataID   string `json:"@odata.id"`
			OdataEtag string `json:"@odata.etag"`
			ID        string `json:"id"`
		}{
			OdataType: "type",
			OdataID:   "id",
			OdataEtag: "etag",
			ID:        "resourceDataID",
		},
		ClientState: "clientState",
		TenantID:    "tenant123",
	}

	callback := Callback{
		Value: []WebhookEventStruct{event},
	}

	body, err := json.Marshal(callback)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook_MicrosoftGraph)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "received"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestWebhookMicrosoftGraph_BadRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/webhook", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webhook_MicrosoftGraph)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
