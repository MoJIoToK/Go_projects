package mongo

import (
	"GoNews/internal/model"
	"fmt"
	"math/rand"
	"testing"
)

var (
	path  string = "mongodb://localhost:27017/"
	posts        = []model.Post{
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
)

func TestNew(t *testing.T) {

	_, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorage_AddPost(t *testing.T) {

	st, err := New(path)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		news []model.Post
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "OK",
			s:       st,
			args:    args{news: posts},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.AddPost(tt.args.news)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddPost() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetPosts(t *testing.T) {

	st, err := New(path)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		n int
	}
	tests := []struct {
		name    string
		s       *Storage
		args    args
		wantErr bool
	}{
		{
			name:    "OK",
			s:       st,
			args:    args{n: 2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetPosts(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.args.n {
				t.Errorf("Storage.Posts() = %v, want %v", len(got), tt.args.n)
			}
		})
	}
}
