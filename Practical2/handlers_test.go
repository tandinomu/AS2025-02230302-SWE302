// handlers_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// A helper function to reset the state before each test
func resetState() {
	users = make(map[int]User)
	nextID = 1
}

func TestCreateUserHandler(t *testing.T) {
	resetState()

	// 1. Define the user we want to create
	userPayload := `{"name": "John Doe"}`
	req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(userPayload))
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// 3. Create a router and serve the request
	router := chi.NewRouter()
	router.Post("/users", createUserHandler)
	router.ServeHTTP(rr, req)

	// 4. Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// 5. Check the response body
	var createdUser User
	if err := json.NewDecoder(rr.Body).Decode(&createdUser); err != nil {
		t.Fatal(err)
	}

	if createdUser.Name != "John Doe" {
		t.Errorf("handler returned unexpected body: got name %v want %v", createdUser.Name, "John Doe")
	}

	if createdUser.ID != 1 {
		t.Errorf("handler returned unexpected body: got id %v want %v", createdUser.ID, 1)
	}
}

func TestCreateUserHandler_InvalidJSON(t *testing.T) {
	resetState()

	// Test with invalid JSON
	invalidPayload := `{"name": }`
	req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(invalidPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Post("/users", createUserHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetAllUsersHandler(t *testing.T) {
	resetState()

	// Add some test users
	users[1] = User{ID: 1, Name: "Alice"}
	users[2] = User{ID: 2, Name: "Bob"}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Get("/users", getAllUsersHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var userList []User
	if err := json.NewDecoder(rr.Body).Decode(&userList); err != nil {
		t.Fatal(err)
	}

	if len(userList) != 2 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(userList), 2)
	}
}

func TestGetAllUsersHandler_EmptyList(t *testing.T) {
	resetState()

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Get("/users", getAllUsersHandler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var userList []User
	if err := json.NewDecoder(rr.Body).Decode(&userList); err != nil {
		t.Fatal(err)
	}

	if len(userList) != 0 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(userList), 0)
	}
}

func TestGetUserHandler(t *testing.T) {
	resetState()

	// First, create a user to fetch
	users[1] = User{ID: 1, Name: "Jane Doe"}

	// Test case 1: User found
	t.Run("User Found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/users/{id}", getUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var foundUser User
		json.NewDecoder(rr.Body).Decode(&foundUser)
		if foundUser.Name != "Jane Doe" {
			t.Errorf("handler returned unexpected body: got %v want %v", foundUser.Name, "Jane Doe")
		}
	})

	// Test case 2: User not found
	t.Run("User Not Found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/99", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/users/{id}", getUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 3: Invalid ID
	t.Run("Invalid ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/users/{id}", getUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestUpdateUserHandler(t *testing.T) {
	resetState()

	// Create a user to update
	users[1] = User{ID: 1, Name: "Original Name"}

	// Test case 1: Successful update
	t.Run("Successful Update", func(t *testing.T) {
		updatePayload := `{"name": "Updated Name"}`
		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(updatePayload))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Put("/users/{id}", updateUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedUser User
		if err := json.NewDecoder(rr.Body).Decode(&updatedUser); err != nil {
			t.Fatal(err)
		}

		if updatedUser.Name != "Updated Name" {
			t.Errorf("handler returned unexpected name: got %v want %v", updatedUser.Name, "Updated Name")
		}

		if updatedUser.ID != 1 {
			t.Errorf("handler returned unexpected ID: got %v want %v", updatedUser.ID, 1)
		}

		// Verify the user was actually updated in the map
		if users[1].Name != "Updated Name" {
			t.Error("user was not actually updated in the map")
		}
	})

	// Test case 2: User not found
	t.Run("User Not Found", func(t *testing.T) {
		updatePayload := `{"name": "Updated Name"}`
		req, err := http.NewRequest("PUT", "/users/99", bytes.NewBufferString(updatePayload))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Put("/users/{id}", updateUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 3: Invalid ID
	t.Run("Invalid ID", func(t *testing.T) {
		updatePayload := `{"name": "Updated Name"}`
		req, err := http.NewRequest("PUT", "/users/invalid", bytes.NewBufferString(updatePayload))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Put("/users/{id}", updateUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test case 4: Invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		invalidPayload := `{"name": }`
		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBufferString(invalidPayload))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Put("/users/{id}", updateUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestDeleteUserHandler(t *testing.T) {
	resetState()

	// Test case 1: Successful deletion
	t.Run("Successful Deletion", func(t *testing.T) {
		users[1] = User{ID: 1, Name: "To Be Deleted"}

		req, err := http.NewRequest("DELETE", "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Delete("/users/{id}", deleteUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		// Verify the user was actually deleted
		if _, ok := users[1]; ok {
			t.Error("user was not deleted from the map")
		}
	})

	// Test case 2: User not found
	t.Run("User Not Found", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/users/99", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Delete("/users/{id}", deleteUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test case 3: Invalid ID
	t.Run("Invalid ID", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/users/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Delete("/users/{id}", deleteUserHandler)
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
