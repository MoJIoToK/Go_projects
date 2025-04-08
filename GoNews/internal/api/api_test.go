package api

import (
	"GoNews/internal/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock для тестирования методов пакета API для подключения БД
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetPosts(id int) ([]model.Post, error) {
	args := m.Called(id)
	return args.Get(0).([]model.Post), args.Error(1)
}

func (m *MockStorage) AddPost(post []model.Post) (int, error) {
	args := m.Called(post)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestAPI_PostsHandler(t *testing.T) {

	mockStorage := new(MockStorage)
	api := New(mockStorage)

	expectedPosts := []model.Post{
		{
			Title:   fmt.Sprintf("Test post %d", rand.Int()),
			Content: "Test content",
			Link:    "https://google.com",
			PubTime: int64(rand.Int()),
		},
		{
			Title:   fmt.Sprintf("Test post %d", rand.Int()),
			Content: "Test content",
			Link:    "https://google.com",
			PubTime: int64(rand.Int()),
		},
	}

	mockStorage.On("GetPosts", 1).Return(expectedPosts, nil)

	req, err := http.NewRequest(http.MethodGet, "/news/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/news/{id:[0-9]+}", api.PostsHandler).Methods(http.MethodGet)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responsePosts []model.Post
	err = json.NewDecoder(rr.Body).Decode(&responsePosts)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedPosts, responsePosts)

	mockStorage.AssertExpectations(t)
}
