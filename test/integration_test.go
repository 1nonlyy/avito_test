package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func TestIntegrationFlow(t *testing.T) {
	moderatorToken := registerAndLogin(t, "mod@avito.ru", "1234", "moderator")
	pvzID := createPVZ(t, moderatorToken, "Казань")

	employeeToken := registerAndLogin(t, "emp@avito.ru", "1234", "employee")
	createReception(t, employeeToken, pvzID)

	for i := 0; i < 50; i++ {
		addProduct(t, employeeToken, pvzID, "электроника")
	}

	closeReception(t, employeeToken, pvzID)
}

func registerAndLogin(t *testing.T, email, password, role string) string {
	req := map[string]string{"role": role}
	body, _ := json.Marshal(req)

	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("dummyLogin failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("dummyLogin returned %d: %s", resp.StatusCode, string(b))
	}

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode dummyLogin response: %v", err)
	}

	token := res["token"]
	if token == "" {
		t.Fatal("dummyLogin returned empty token")
	}
	return token
}

func createPVZ(t *testing.T, token, city string) string {
	req := map[string]string{"city": city}
	body, _ := json.Marshal(req)
	resp := authenticatedRequest(t, "POST", "/pvz", token, body)
	defer resp.Body.Close()

	var res map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	return res["id"].(string)
}

func createReception(t *testing.T, token, pvzID string) {
	req := map[string]string{"pvzId": pvzID}
	body, _ := json.Marshal(req)
	resp := authenticatedRequest(t, "POST", "/receptions", token, body)
	defer resp.Body.Close()
}

func addProduct(t *testing.T, token, pvzID, prodType string) {
	req := map[string]string{"pvzId": pvzID, "type": prodType}
	body, _ := json.Marshal(req)
	resp := authenticatedRequest(t, "POST", "/products", token, body)
	defer resp.Body.Close()
}

func closeReception(t *testing.T, token, pvzID string) {
	path := fmt.Sprintf("/pvz/%s/close_last_reception", pvzID)
	resp := authenticatedRequest(t, "POST", path, token, []byte("{}"))
	defer resp.Body.Close()
}

func authenticatedRequest(t *testing.T, method, path, token string, body []byte) *http.Response {
	req, _ := http.NewRequest(method, baseURL+path, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request to %s failed: %v", path, err)
	}
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("Request to %s failed with status %d: %s", path, resp.StatusCode, string(b))
	}
	return resp
}
