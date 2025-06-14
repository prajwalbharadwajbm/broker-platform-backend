package validator

import (
	"errors"
	"strings"
)

func IsValidEmail(email string) (bool, error) {
	if !strings.Contains(email, "@") {
		return false, errors.New("BPB002")
	}

	return true, nil

}

func IsValidPassword(username, password string) (bool, error) {
	if !(len(password) >= 8) {
		return false, errors.New("BPB003")
	}
	if password == username {
		return false, errors.New("BPB004")
	}
	return true, nil
}
