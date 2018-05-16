package auth_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/mappcpd/web-services/internal/auth"
	//"github.com/mappcpd/web-services/internal/auth"
	"github.com/mappcpd/web-services/testdata"
)

var db = testdata.NewDataStore()
var helper = testdata.NewHelper()

func TestMain(m *testing.M) {
	err := db.SetupMySQL()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.TearDownMySQL()

	m.Run()
}

func TestPingDatabase(t *testing.T) {
	err := db.Store.MySQL.Session.Ping()
	if err != nil {
		t.Fatal("Could not ping database")
	}
}

func TestAuthMemberClearPass(t *testing.T) {
	id, name, err := auth.AuthMember(db.Store, "michael@mesa.net.au", "password")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Michael Donnici", name)
}

func TestAuthMemberMD5Pass(t *testing.T) {
	id, name, err := auth.AuthMember(db.Store, "michael@mesa.net.au", "5f4dcc3b5aa765d61d8327deb882cf99")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Michael Donnici", name)
}

func TestAuthMemberFail(t *testing.T) {
	id, _, err := auth.AuthMember(db.Store, "michael@mesa.net.au", "wrongPassword")
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, sql.ErrNoRows, err)
	helper.Result(t, 0, id)
}

func TestAuthAdminClearPass(t *testing.T) {
	id, name, err := auth.AdminAuth(db.Store, "demo-admin", "demo-admin")
	if err == sql.ErrNoRows {
		t.Log("Expected result to fail login")
	}
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Demo Admin", name)
}

func TestAuthAdminMD5Pass(t *testing.T) {
	id, name, err := auth.AdminAuth(db.Store, "demo-admin", "41d0510a9067999b72f38ba0ce9f6195")
	if err != nil {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, 1, id)
	helper.Result(t, "Demo Admin", name)
}

func TestAuthAdminFail(t *testing.T) {
	id, _, err := auth.AdminAuth(db.Store, "demo-admin", "wrongPassword")
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Database error: %s", err)
	}
	helper.Result(t, sql.ErrNoRows, err)
	helper.Result(t, 0, id)
}
