package validator

import (
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	t.Run("IsValidEmail_ValidEmail", func(t *testing.T) {
		validEmails := []string{
			// Examples from https://en.wikipedia.org/wiki/Email_address#Examples
			"simple@example.com",
			"FirstName.LastName@EasierReading.org",
			"long.email-address-with-hyphens@and.subdomains.example.com",
			"user%example.com@example.org",
		}

		for _, email := range validEmails {
			valid, err := IsValidEmail(email)
			if err != nil {
				t.Errorf("Expected no error for valid email '%s', got %v", email, err)
			}
			if !valid {
				t.Errorf("Expected email '%s' to be valid", email)
			}
		}
	})

	t.Run("IsValidEmail_InvalidEmail", func(t *testing.T) {
		invalidEmails := []string{
			// Examples from https://en.wikipedia.org/wiki/Email_address#Examples
			"abc.example.com",
			"just\"not\"right@example.com",
		}

		for _, email := range invalidEmails {
			valid, err := IsValidEmail(email)
			if err == nil {
				t.Errorf("Expected error for invalid email '%s', got nil", email)
			}
			if valid {
				t.Errorf("Expected email '%s' to be invalid", email)
			}

			// Check error code
			if err != nil && err.Error() != "BPB002" {
				t.Errorf("Expected error code 'BPB002' for email '%s', got '%s'", email, err.Error())
			}
		}
	})

	t.Run("IsValidEmail_EmptyString", func(t *testing.T) {
		valid, err := IsValidEmail("")
		if err == nil {
			t.Error("Expected error for empty email")
		}
		if valid {
			t.Error("Expected empty email to be invalid")
		}
		if err.Error() != "BPB002" {
			t.Errorf("Expected error code 'BPB002', got '%s'", err.Error())
		}
	})
}

func TestIsValidPassword(t *testing.T) {
	t.Run("IsValidPassword_ValidPasswords", func(t *testing.T) {
		testCases := []struct {
			username string
			password string
		}{
			{"user123", "password123"},
			{"admin", "strongP@ssw0rd"},
			{"testuser", "12345678"},
			{"john", "verylongpassword"},
			{"alice", "P@ssw0rd!"},
		}

		for _, tc := range testCases {
			valid, err := IsValidPassword(tc.username, tc.password)
			if err != nil {
				t.Errorf("Expected no error for username '%s' and password '%s', got %v",
					tc.username, tc.password, err)
			}
			if !valid {
				t.Errorf("Expected password '%s' to be valid for username '%s'",
					tc.password, tc.username)
			}
		}
	})

	t.Run("IsValidPassword_TooShort", func(t *testing.T) {
		shortPasswords := []string{
			"",
			"1",
			"12",
			"123",
			"1234",
			"12345",
			"123456",
			"1234567",
		}

		username := "testuser"
		for _, password := range shortPasswords {
			valid, err := IsValidPassword(username, password)
			if err == nil {
				t.Errorf("Expected error for short password '%s'", password)
			}
			if valid {
				t.Errorf("Expected password '%s' to be invalid (too short)", password)
			}
			if err != nil && err.Error() != "BPB003" {
				t.Errorf("Expected error code 'BPB003' for password '%s', got '%s'",
					password, err.Error())
			}
		}
	})

	t.Run("IsValidPassword_SameAsUsername", func(t *testing.T) {
		testCases := []struct {
			username string
			password string
		}{
			{"username", "username"},
			{"admin123", "admin123"},
			{"testuser", "testuser"},
			{"12345678", "12345678"},
		}

		for _, tc := range testCases {
			valid, err := IsValidPassword(tc.username, tc.password)
			if err == nil {
				t.Errorf("Expected error when password equals username ('%s')", tc.username)
			}
			if valid {
				t.Errorf("Expected password to be invalid when it equals username ('%s')", tc.username)
			}
			if err != nil && err.Error() != "BPB004" {
				t.Errorf("Expected error code 'BPB004' for username/password '%s', got '%s'",
					tc.username, err.Error())
			}
		}
	})

	t.Run("IsValidPassword_MinLength", func(t *testing.T) {
		// Test exactly 8 characters (minimum length)
		username := "testuser"
		password := "12345678"

		valid, err := IsValidPassword(username, password)
		if err != nil {
			t.Errorf("Expected no error for 8-character password, got %v", err)
		}
		if !valid {
			t.Error("Expected 8-character password to be valid")
		}
	})

	t.Run("IsValidPassword_CaseSensitive", func(t *testing.T) {
		// Test case sensitivity
		testCases := []struct {
			username string
			password string
			expected bool
		}{
			{"User", "user12345", true},     // Different case should be valid
			{"USER", "user12345", true},     // Different case should be valid
			{"user", "USER12345", true},     // Different case should be valid
			{"TestUser", "TestUser", false}, // Same case should be invalid
		}

		for _, tc := range testCases {
			valid, err := IsValidPassword(tc.username, tc.password)
			if tc.expected {
				if err != nil {
					t.Errorf("Expected no error for username '%s' and password '%s', got %v",
						tc.username, tc.password, err)
				}
				if !valid {
					t.Errorf("Expected password '%s' to be valid for username '%s'",
						tc.password, tc.username)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for username '%s' and password '%s'",
						tc.username, tc.password)
				}
				if valid {
					t.Errorf("Expected password '%s' to be invalid for username '%s'",
						tc.password, tc.username)
				}
			}
		}
	})

	t.Run("IsValidPassword_EmptyUsername", func(t *testing.T) {
		// Test with empty username
		valid, err := IsValidPassword("", "validpassword123")
		if err != nil {
			t.Errorf("Expected no error for empty username, got %v", err)
		}
		if !valid {
			t.Error("Expected password to be valid with empty username")
		}

		// Test when password equals empty username (should be valid as password is 8+ chars)
		valid, err = IsValidPassword("", "")
		if err == nil {
			t.Error("Expected error for empty password")
		}
		if valid {
			t.Error("Expected empty password to be invalid")
		}
	})
}
