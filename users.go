package storj

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type UserService struct {
	client *Client
}

type User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	Pubkey          string `json:"pubkey,omitempty"`
	ReferralPartner string `json:"referralPartner,omitempty"`
}

func hashPassword(password string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(password))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (u *UserService) Create(user User) error {
	if user.Password == "" {
		return fmt.Errorf("Password field empty")
	}

	p, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = p

	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	rel, _ := url.Parse("/users")
	url := u.client.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	return nil
}
