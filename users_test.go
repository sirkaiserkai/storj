package storj

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

var users = map[string]User{}

func TestUsersCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, "POST")
		assertHeader(t, r, "Content-Type", "application/json")

		m := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			t.Errorf("Failed to parse request body")
		}

		e, ok := m["email"].(string)
		if !ok {
			t.Error("missing \"email\" field")
			w.WriteHeader(400)
			return
		}

		p, ok := m["password"].(string)
		if !ok {
			t.Error("mssing \"password\" field")
			w.WriteHeader(400)
			return
		}

		if _, ok := users[e]; ok {
			fmt.Fprintf(w, `{"error":"Email is already registered"}`)
			w.WriteHeader(400)
			return
		}

		users[e] = User{Email: e, Password: p}

		w.WriteHeader(201)
	})

	err := client.Users.Create(User{Email: "test@email.com", Password: "password"})
	if err != nil {
		t.Errorf("Users.Create returned error: %v", err)
	}

	err = client.Users.Create(User{Email: "test@email.com", Password: "password"})
	if err == nil {
		t.Errorf("Users.Create returned error: should cause duplicate error")
	}
}
