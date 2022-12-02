package teddit_test

import (
	"os"
	"sync"
	"teniditter-server/cmd/services/teddit"
	"testing"
)

func init() {
	os.Setenv("TEST", "1")
}

func TestGetHomePosts(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(5)
	for _, homeType := range []string{"hot", "new", "top", "rising", "controversial"} {
		go func(homeType string) {
			defer wg.Done()
			if res, err := teddit.GetHomePosts(homeType, "", true); err != nil || len(*res) <= 0 {
				t.Errorf("GetHomePosts() + type %s failed", homeType)
			}
		}(homeType)
	}
	wg.Wait()
}
func TestGetPostInfo(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(6)
	for _, sort := range []string{"best", "top", "new", "controversial", "old", "qa"} {
		go func(sort string) {
			defer wg.Done()
			if res, err := teddit.GetPostInfo("golang", "za4kvu", sort); err != nil || res.Metadata.Author != "Jamo008" {
				t.Errorf("GetPostInfo() + sort %s failed", sort)
			}
		}(sort)
	}
	wg.Wait()
}

func TestGetSubredditPosts(t *testing.T) {
	if res, err := teddit.GetSubredditPosts("golang"); err != nil || len(*res) <= 0 {
		t.Fatal("GetSubredditPosts() failed")
	}
}
func TestGetSubredditMetadatas(t *testing.T) {
	if _, err := teddit.GetSubredditMetadatas("golang"); err != nil {
		t.Fatal("GetSubredditMetadatas() failed")
	}
}
func TestGetUserInfos(t *testing.T) {
	if res, err := teddit.GetUserInfos("Jamo008"); err != nil && len(*res) <= 0 {
		t.Fatal("GetUserInfos() failed")
	}
}
