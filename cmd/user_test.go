package main

import (
	"GravitumTask/internal/user"
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func prepareData(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	db.Create(&user.User{
		Model: gorm.Model{
			ID: 1,
		},
		Password: string(password),
		Email:    "user@example.com",
		Name:     "gravitum",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().Where("email = ?", "user@example.com").Delete(&user.User{})
}

func TestUpdateSuccess(t *testing.T) {
	db := initDb()
	prepareData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()
	defer removeData(db)
	data, _ := json.Marshal(user.UserUpdateRequest{
		Password: "password",
		Email:    "user@example.com",
		Name:     "gravitum1",
	})
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/users/1", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status Ok, got %v", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return
	}
	var userResponse user.User
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		t.Fatal(err)
		return
	}
	if userResponse.Name != "gravitum1" {
		t.Fatalf("expected user name 'gravitum_changed', got %v", userResponse.Name)
		return
	}
}

func TestUpdateFail(t *testing.T) {
	db := initDb()
	prepareData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()
	defer removeData(db)
	data, _ := json.Marshal(user.UserUpdateRequest{
		Name:     "Gravitum",
		Email:    "user@example.com",
		Password: "password1",
	})
	req, err := http.NewRequest("PUT", ts.URL+"/users/1", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status BadRequest, got %v", resp.StatusCode)
		return
	}
}

func TestGetSuccess(t *testing.T) {
	db := initDb()
	prepareData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()
	defer removeData(db)
	resp, err := http.Get(ts.URL + "/users/1")
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status Ok, got %v", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return
	}
	var userResponse user.User
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		t.Fatal(err)
		return
	}
	if userResponse.Name != "gravitum" || userResponse.Email != "user@example.com" {
		t.Fatalf("expected user name 'gravitum' and email 'user@example.com', got %v, %v", userResponse.Name, userResponse.Email)
		return
	}
}

func TestPostSuccess(t *testing.T) {
	db := initDb()
	ts := httptest.NewServer(App())
	defer ts.Close()
	defer removeData(db)
	data, _ := json.Marshal(user.UserCreateRequest{
		Name:     "gravitum",
		Email:    "user@example.com",
		Password: "password",
	})
	resp, err := http.Post(ts.URL+"/users", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected code StatusCreated, got %v", resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return
	}
	var userResponse user.User
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		t.Fatal(err)
		return
	}
	if userResponse.Name != "gravitum" || userResponse.Email != "user@example.com" {
		t.Fatalf("expected user name 'gravitum', got %v, %v", userResponse.Name, userResponse.Email)
		return
	}
	var user user.User
	result := db.Clauses(clause.Returning{}).First(&user, "name = ?", userResponse.Name)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if userResponse.Name != user.Name {
		t.Fatalf("expected user name 'gravitum', got %v", userResponse.Name)
		return
	}
}
