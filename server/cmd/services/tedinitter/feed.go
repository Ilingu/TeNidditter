package teddinitter

import (
	"fmt"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/console"
)

// Generate all users feed
func GenerateFeeds() {
	accounts, err := db.GetAllAccounts(true)
	if err != nil {
		return
	}

	console.Log(fmt.Sprintf("Generating feeds for %d users...", len(accounts)))
	for _, user := range accounts {
		go user.GenerateTedditFeed()
		go user.GenerateNitterFeed()
	}
}
