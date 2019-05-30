package auth_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/cardiacsociety/web-services/internal/auth"
	"github.com/cardiacsociety/web-services/internal/platform/datastore"
	"github.com/cardiacsociety/web-services/testdata"
)

var ds datastore.Datastore

func TestAuth(t *testing.T) {

	var teardown func()
	ds, teardown = setup()
	defer teardown()

	t.Run("auth", func(t *testing.T) {
		t.Run("testPingDatabase", testPingDatabase)
		t.Run("testAuthMemberClearPass", testAuthMemberClearPass)
		t.Run("testAuthMemberMD5Pass", testAuthMemberMD5Pass)
		t.Run("testAuthMemberFail", testAuthMemberFail)
		t.Run("testAuthAdminClearPass", testAuthAdminClearPass)
		t.Run("testAuthAdminMD5Pass", testAuthAdminMD5Pass)
		t.Run("testAuthAdminFail", testAuthAdminFail)
	})
}

func setup() (datastore.Datastore, func()) {
	var db = testdata.NewDataStore()
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalf("SetupMySQL() err = %s", err)
	}
	return db.Store, func() {
		err := db.TearDownMySQL()
		if err != nil {
			log.Fatalf("TearDownMySQL() err = %s", err)
		}
	}
}

func testPingDatabase(t *testing.T) {
	err := ds.MySQL.Session.Ping()
	if err != nil {
		t.Fatalf("Ping() err = %s", err)
	}
}

func testAuthMemberClearPass(t *testing.T) {
	gotId, gotName, err := auth.AuthMember(ds, "michael@mesa.net.au", "password")
	if err != nil {
		t.Fatalf("auth.AuthMember() err = %s", err)
	}
	wantId := 1
	if gotId != wantId {
		t.Errorf("Auth.AuthMember() id = %d, want %d", gotId, wantId)
	}
	wantName := "Michael Donnici"
	if gotName != wantName {
		t.Errorf("Auth.AuthMember() name = %q, want %q", gotName, wantName)
	}
}

func testAuthMemberMD5Pass(t *testing.T) {
	gotId, gotName, err := auth.AuthMember(ds, "michael@mesa.net.au", "5f4dcc3b5aa765d61d8327deb882cf99")
	if err != nil {
		t.Fatalf("auth.AuthMember() err = %s", err)
	}
	wantId := 1
	if gotId != wantId {
		t.Errorf("Auth.AuthMember() id = %d, want %d", gotId, wantId)
	}
	wantName := "Michael Donnici"
	if gotName != wantName {
		t.Errorf("Auth.AuthMember() name = %q, want %q", gotName, wantName)
	}
}

func testAuthMemberFail(t *testing.T) {
	_, _, err := auth.AuthMember(ds, "michael@mesa.net.au", "wrongPassword")
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("auth.AuthMember() err = %s", err)
	}
	if err == nil {
		t.Errorf("auth.AuthMember() err = %v, want %v", err, sql.ErrNoRows)
	}
}

func testAuthAdminClearPass(t *testing.T) {
	gotId, gotName, err := auth.AdminAuth(ds, "demo-admin", "demo-admin")
	if err != nil {
		t.Fatalf("auth.AdminAuth() err = %s", err)
	}
	wantId := 1
	if gotId != wantId {
		t.Errorf("Auth.AdminAuth() id = %d, want %d", gotId, wantId)
	}
	wantName := "Demo Admin"
	if gotName != wantName {
		t.Errorf("Auth.AdminAuth() name = %q, want %q", gotName, wantName)
	}
}

func testAuthAdminMD5Pass(t *testing.T) {
	gotId, gotName, err := auth.AdminAuth(ds, "demo-admin", "41d0510a9067999b72f38ba0ce9f6195")
	if err != nil {
		t.Fatalf("auth.AdminAuth() err = %s", err)
	}
	wantId := 1
	if gotId != wantId {
		t.Errorf("Auth.AdminAuth() id = %d, want %d", gotId, wantId)
	}
	wantName := "Demo Admin"
	if gotName != wantName {
		t.Errorf("Auth.AdminAuth() name = %q, want %q", gotName, wantName)
	}
}

func testAuthAdminFail(t *testing.T) {
	_, _, err := auth.AdminAuth(ds, "demo-admin", "wrongPassword")
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("auth.AdminAuth() err = %s", err)
	}
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("auth.AdminAuth() err = %s", err)
	}
	if err == nil {
		t.Errorf("auth.AdminAuth() err = %v, want %v", err, sql.ErrNoRows)
	}
}
