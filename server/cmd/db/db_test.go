/*
This integration tests package do a app overview of all the db interactions in one big test
*/
package db_test

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"runtime"
	"teniditter-server/cmd/db"
	"teniditter-server/cmd/global/utils"
	utils_enc "teniditter-server/cmd/global/utils/encryption"
	utils_env "teniditter-server/cmd/global/utils/env"
	ps "teniditter-server/cmd/planetscale"
	"teniditter-server/cmd/services/nitter"
	"testing"
	"time"
)

/* HELPERS */

func generateRandomUsername() string {
	random, err0 := rand.Int(rand.Reader, big.NewInt(1e18))
	wd, err1 := os.Getwd()
	hd, err2 := os.UserHomeDir()
	ex, err3 := os.Executable()
	hn, err4 := os.Hostname()
	rc, err5 := utils.GenerateRandomChars(32)

	for _, err := range []error{err0, err1, err2, err3, err4, err5} {
		if err != nil {
			panic(err)
		}
	}

	username_args := []any{runtime.Version(), runtime.GOOS, runtime.GOARCH, time.Now().UnixNano(), random.String(), wd, hd, ex, hn, runtime.NumGoroutine(), rc}
	utils.ShuffleSlice(username_args)

	return utils_enc.GenerateHashFromArgs(username_args...)
}

func generateRandomPassword() string {
	psw, err := utils.GenerateRandomChars(64)
	if err != nil {
		panic("cannot generate psw")
	}
	return psw
}

func LoadAccount() {
	user, err := db.GetAccount(username)
	if err != nil {
		panic("failed to fetch user by username")
	}
	currUser = user
}

/* INTEGRATION TESTS */

var (
	username = generateRandomUsername()
	password = generateRandomPassword()
	currUser *db.AccountModel
)

func init() {
	os.Setenv("TEST", "1")
	utils_env.LoadEnv()

	db := ps.DBManager.NewDB()
	if db == nil {
		panic("cannot connect to db")
	}
}

// func TestMain(m *testing.M) {
// 	exitVal := m.Run()
// 	os.Exit(exitVal)
// 	defer ps.DBManager.Disconnect()
// }

// DO NOT RUN TEST INDIVIDUALLY, this test handle all the others tests
func TestApp(t *testing.T) {
	defer ps.DBManager.Disconnect()

	testCreateAccount(t)       // [ORDER_IS_IMPORTANT]: create test user
	defer testDeleteAccount(t) // if other tests failed the test user will still be deleted

	testGetAccount(t) //  [ORDER_IS_IMPORTANT]: load user into "currUser" variable for other tests
	testGetAllAccounts(t)

	testPasswordMatch(t)
	testUpdatePassword(t)
	LoadAccount() //  [ORDER_IS_IMPORTANT]: reload user into "currUser" variable for other tests after password change
	testPasswordMatch(t)

	testRecoveryCodes(t)

	testNitter(t)
	testTeddit(t)

	testLists(t)
}

func testCreateAccount(t *testing.T) {
	if utils.IsEmptyString(username) || utils.IsEmptyString(password) {
		t.Fatal("no username/psw")
		return
	}

	if _, err := db.CreateAccount(username, password); err != nil {
		t.Fatal("failed to create account", err)
	}
}

func testGetAccount(t *testing.T) {
	user, err := db.GetAccount(username)
	if err != nil {
		t.Fatal("failed to fetch user by username")
	}
	currUser = user

	user2, err := db.GetAccount(user.AccountId)
	if err != nil {
		t.Fatal("failed to fetch user by id")
	}

	if utils_enc.GenerateHashFromArgs(user) != utils_enc.GenerateHashFromArgs(user2) {
		t.Fatal("users not the same")
	}
}

func testGetAllAccounts(t *testing.T) {
	if accounts, err := db.GetAllAccounts(false); err != nil || len(accounts) <= 0 {
		t.Fatal("GetAllAccounts(false) failed", err)
	}
	if _, err := db.GetAllAccounts(true); err != nil {
		t.Fatal("GetAllAccounts(true) failed", err)
	}
}

func testDeleteAccount(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	if ok := db.DeleteAccount(currUser); !ok {
		t.Fatal("failed to delete account")
	}

	account, err := db.GetAccount(currUser.AccountId)
	if userExist := err == nil && account != nil; userExist {
		t.Fatal("user hasn't been deleted")
	}
}

func testPasswordMatch(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	if match := currUser.PasswordMatch(password); !match {
		t.Fatal("user password doesn't match the db one")
	}
}

func testUpdatePassword(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	newPassword := generateRandomPassword()
	if err := currUser.UpdatePassword(newPassword); err != nil {
		t.Fatal("UpdatePassword() failed")
	}
	password = newPassword
}

func testRecoveryCodes(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	// test 1
	oldCodes, err := currUser.DecryptRecoveryCodes()
	if err != nil {
		t.Fatal("couldn't decrypt codes")
	}

	if len(*oldCodes) != db.RECOVERY_CODES_AMOUNT {
		t.Fatal("invalid number of recovery codes", len(*oldCodes), "!=", db.RECOVERY_CODES_AMOUNT)
	}
	for _, code := range *oldCodes {
		if len(code) != db.RECOVERY_CODE_LENGTH {
			t.Fatal("invalid length of recovery code", len(code), "!=", db.RECOVERY_CODE_LENGTH)
		}
	}

	// test 2
	random, err := rand.Int(rand.Reader, big.NewInt(int64(db.RECOVERY_CODES_AMOUNT)))
	if err != nil {
		t.Fatal(err)
	}
	i, err := utils.BigIntToInt(random, 8)
	if err != nil {
		t.Fatal(err)
	}

	codeRemoved := (*oldCodes)[i]
	if err := currUser.UseRecoveryCode(codeRemoved); err != nil {
		t.Fatalf("UseRecoveryCode(%s) failed, %s", (*oldCodes)[i], err.Error())
	}

	// reload datas
	LoadAccount() // reload user into "currUser" variable for other tests after recovery codes changes

	// test 3
	codeAdded, err := utils.GenerateRandomChars(uint(db.RECOVERY_CODE_LENGTH))
	if err != nil {
		t.Fatal(err)
	}
	if err := currUser.AddRecoveryCode(codeAdded); err != nil {
		t.Fatalf("AddRecoveryCode(%s) failed, %s", codeAdded, err.Error())
	}

	// reload datas
	LoadAccount() // reload user into "currUser" variable for other tests after recovery codes changes

	newCodes, err := currUser.DecryptRecoveryCodes()
	if err != nil {
		t.Fatal("couldn't decrypt codes", err)
	}

	// check datas reloaded
	if utils_enc.GenerateHashFromArgs(oldCodes) == utils_enc.GenerateHashFromArgs(newCodes) {
		t.Fatal("no updates on codes")
	}

	// check test 2
	if exist := currUser.HasRecoveryCode(codeRemoved); exist {
		t.Fatal(codeRemoved, "should have been removed from the recovery codes")
	}

	// check test 3
	if exist := currUser.HasRecoveryCode(codeAdded); !exist {
		t.Fatal(codeAdded, "should have been added to the recovery codes")
	}

	// test 4
	regeneratedCodes, err := currUser.RegenerateRecoveryCode()
	if err != nil {
		t.Fatal("couldn't regenerate codes", err)
	}

	// reload datas
	LoadAccount() // reload user into "currUser" variable for other tests after recovery codes changes
	newCodesAfterRegeneration, err := currUser.DecryptRecoveryCodes()
	if err != nil {
		t.Fatal("couldn't decrypt codes", err)
	}

	// check test 4
	if utils_enc.GenerateHashFromArgs(newCodes) == utils_enc.GenerateHashFromArgs(regeneratedCodes) {
		t.Fatal("codes not regenerated")
	}
	if utils_enc.GenerateHashFromArgs(newCodesAfterRegeneration) != utils_enc.GenerateHashFromArgs(regeneratedCodes) {
		t.Fatal("codes not regenerated")
	}
}

func testNitter(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	// test 1
	nittosName := generateRandomUsername()
	nittos, err := db.GetNittos(nittosName)
	if err != nil {
		t.Fatalf("GetNittos(%s) failed, %s", nittosName, err)
	}
	if nittos.Username != nittosName {
		t.Fatalf("expected: %s, got: %s", nittosName, nittos.Username)
	}

	// test 2
	result, err := db.SearchNittos(nittosName)
	if err != nil {
		t.Fatalf("SearchNittos(%s) failed, %s", nittosName, err)
	}

	// check test 1
	if len(result) < 1 {
		t.Fatalf("SearchNittos(%s): len of 0, expected non-empty", nittosName)
	}

	// test 3
	if success := currUser.SubTo(nittos); !success {
		t.Fatal("SubTo() failed")
	}

	// check test 3
	subs, err := currUser.GetNitterSubs()
	if err != nil {
		t.Fatal("GetNitterSubs() failed", err)
	}
	if len(subs) != 1 || subs[0] != nittosName {
		t.Fatal("GetNitterSubs() returned something not expected:", subs, "does not have", nittosName, "as an only child")
	}

	accounts, err := db.GetAllAccounts(true)
	if err != nil {
		t.Fatal("GetAllAccounts(true) failed", err)
	}
	if len(accounts) <= 0 {
		t.Fatal("GetAllAccounts(true) gave an empty result")
	}

	// test 4
	if success := currUser.UnsubFrom(nittos); !success {
		t.Fatal("UnsubFrom() failed")
	}

	// check test 4
	subs, err = currUser.GetNitterSubs()
	if err != nil {
		t.Fatal("GetNitterSubs() failed", err)
	}
	if len(subs) != 0 {
		t.Fatal("GetNitterSubs() returned a non-empty sub list:", subs)
	}
}

func testTeddit(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	// test 1
	subname := generateRandomUsername()
	subteddit, err := db.GetSubteddit(subname)
	if err != nil {
		t.Fatalf("GetSubteddit(%s) failed, %s", subname, err)
	}
	if subteddit.Subname != subname {
		t.Fatalf("expected: %s, got: %s", subname, subteddit.Subname)
	}

	// test 2
	result, err := db.SearchSubteddit(subname)
	if err != nil {
		t.Fatalf("SearchSubteddit(%s) failed, %s", subname, err)
	}

	// check test 1
	if len(result) < 1 {
		t.Fatalf("SearchSubteddit(%s): len of 0, expected non-empty", subname)
	}

	// test 3
	if success := currUser.SubTo(subteddit); !success {
		t.Fatal("SubTo() failed")
	}

	// check test 3
	subs, err := currUser.GetTedditSubs()
	if err != nil {
		t.Fatal("GetTedditSubs() failed", err)
	}
	if len(subs) != 1 || subs[0] != subname {
		t.Fatal("GetTedditSubs() returned something not expected:", subs, "does not have", subname, "as an only child")
	}

	accounts, err := db.GetAllAccounts(true)
	if err != nil {
		t.Fatal("GetAllAccounts(true) failed", err)
	}
	if len(accounts) <= 0 {
		t.Fatal("GetAllAccounts(true) gave an empty result")
	}

	// test 4
	if success := currUser.UnsubFrom(subteddit); !success {
		t.Fatal("UnsubFrom() failed")
	}

	// check test 4
	subs, err = currUser.GetTedditSubs()
	if err != nil {
		t.Fatal("GetTedditSubs() failed", err)
	}
	if len(subs) != 0 {
		t.Fatal("GetTedditSubs() returned a non-empty sub list:", subs)
	}
}

func testLists(t *testing.T) {
	if currUser == nil {
		t.Fatal("no user registered")
	}

	listname := generateRandomUsername()

	// test 1
	if success := db.CreateNitterList(currUser, listname); !success {
		t.Fatal("CreateNitterList() failed")
	}

	// test 2
	lists, err := currUser.GetNitterLists()
	if err != nil {
		t.Fatal("GetNitterLists() failed", err)
	}
	// check test 1 and 2
	if len(lists) != 1 || lists[0].ListName != listname {
		t.Fatal("User lists is invalid, expected:", []db.NitterListModel{{ListName: listname}}, "got:", lists)
	}

	// test 3
	list, err := db.GetListById(lists[0].ListID)
	if err != nil {
		t.Fatalf("GetListById(%d) failed", lists[0].ListID)
	}
	// check test 3
	if utils_enc.GenerateHashFromArgs(*list) != utils_enc.GenerateHashFromArgs(lists[0]) {
		t.Fatalf("Invalid fetch from GetListById(%d) or currUser.GetNitterLists()", lists[0].ListID)
	}

	// test 4
	neetId, err := utils.GenerateRandomChars(19)
	if err != nil {
		t.Fatal(err)
	}

	neetMock := nitter.NeetComment{NeetBasicComment: nitter.NeetBasicComment{
		Id:        neetId,
		Content:   "test",
		Creator:   nitter.NittosPreview{Username: "test", Description: "test"},
		CreatedAt: int(time.Now().Unix()),
		Stats: nitter.NeetCommentStats{
			LikesCounts:  50,
			RTCounts:     3,
			ReplyCounts:  15,
			QuotesCounts: 1,
		},
	}}

	if success := db.InsertNewNeet(neetMock); !success {
		t.Fatal("InsertNewNeet() failed")
	}
	defer func() {
		// test 12
		if success := db.DeleteNeet(neetMock.Id); !success {
			t.Fatalf("DeleteNeet(%s) failed", neetMock.Id)
		}

		// check test 12
		if db.IsNeetAlreadyExist(neetMock.Id) {
			t.Fatal("DeleteNeet() failed to delete")
		}
	}()

	// test 5 + check test 4
	if !db.IsNeetAlreadyExist(neetMock.Id) {
		t.Fatal("InsertNewNeet() failed to add")
	}

	// test 6
	neetDB, err := db.GetNeetById(neetMock.Id)
	if err != nil {
		t.Fatalf("GetNeetById(%s) failed", neetMock.Id)
	}

	// check test 6
	var neet nitter.NeetComment
	if err := json.Unmarshal([]byte(neetDB.Neet), &neet); err != nil {
		t.Fatal(err)
	}
	if utils_enc.GenerateHashFromArgs(neet) != utils_enc.GenerateHashFromArgs(neetMock) {
		t.Fatalf("InsertNewNeet() inserted datas wrongly")
	}

	// test 7
	if success := list.AddNeet(neetMock.Id); !success {
		t.Fatalf("list.AddNeet(%s) failed", neetMock.Id)
	}

	// test 8
	listNeets, err := db.GetListContentById(list.ListID)
	if err != nil {
		t.Fatalf("GetListContentById(%d) failed", list.ListID)
	}

	// check test 7
	if len(listNeets) != 1 || listNeets[0].Id != neetMock.Id {
		t.Fatal("list.AddNeet() added wrongly", listNeets, neetMock.Id)
	}

	// test 9
	if success := list.RemoveNeet(neetMock.Id); !success {
		t.Fatalf("list.RemoveNeet(%s) failed", neetMock.Id)
	}

	// test 10
	listNeets, err = db.GetListContentById(list.ListID)
	if err != nil {
		t.Fatalf("GetListContentById(%d) failed", list.ListID)
	}

	// check test 9
	if len(listNeets) != 0 {
		t.Fatal("list.RemoveNeet() removed wrongly", listNeets)
	}

	// test 11
	if success := list.Delete(); !success {
		t.Fatal("list.Delete() failed")
	}

	// check test 11
	if list, err := db.GetListById(list.ListID); err == nil || list != nil {
		t.Fatal("list.Delete() failed to delete")
	}
}
