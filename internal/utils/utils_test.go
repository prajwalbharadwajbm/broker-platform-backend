package utils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
)

// TestGetEnv tests the GetEnv function with different types
func TestGetEnv(t *testing.T) {
	t.Run("GetEnv_String_WithEnvVar", func(t *testing.T) {
		key := "TEST_STRING"
		expected := "test_value"
		os.Setenv(key, expected)
		defer os.Unsetenv(key)

		result := GetEnv(key, "default")
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("GetEnv_String_WithoutEnvVar", func(t *testing.T) {
		key := "NON_EXISTENT_KEY"
		fallback := "default_value"

		result := GetEnv(key, fallback)
		if result != fallback {
			t.Errorf("Expected %v, got %v", fallback, result)
		}
	})
	// TODO: Add more tests for other types
}

// TestFetchDataFromRequestBody tests the FetchDataFromRequestBody function
func TestFetchDataFromRequestBody(t *testing.T) {
	type TestStruct struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	t.Run("FetchDataFromRequestBody_ValidJSON", func(t *testing.T) {
		expected := TestStruct{Email: "test@test.com", Password: "123"}
		jsonData, _ := json.Marshal(expected)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		result, err := FetchDataFromRequestBody[TestStruct](req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if result.Email != expected.Email || result.Password != expected.Password {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("FetchDataFromRequestBody_InvalidJSON", func(t *testing.T) {
		invalidJSON := bytes.NewBuffer([]byte("{invalid json"))
		req := httptest.NewRequest("POST", "/test", invalidJSON)

		_, err := FetchDataFromRequestBody[TestStruct](req)
		if err == nil {
			t.Fatal("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("FetchDataFromRequestBody_EmptyBody", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer([]byte("")))

		result, err := FetchDataFromRequestBody[TestStruct](req)
		if err != nil {
			t.Fatalf("Expected no error for empty body, got %v", err)
		}

		expected := TestStruct{}
		if result.Email != expected.Email || result.Password != expected.Password {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})
}
