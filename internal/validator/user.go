package validator

import (
	"errors"
	"net/mail"
)

// TODO: Need better validation for email, custom regex may be needed
func IsValidEmail(email string) (bool, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
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
