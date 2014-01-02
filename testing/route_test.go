package testing

import (
	"github.com/dbettin/cork"
	"net/http"
	"testing"
)

func TestSettingsCheckFailed(t *testing.T) {
	if _, err := Get("somepath"); err == nil {
		t.Error("Settings check should have failed since it wasn't set")
	}
}

func TestSettingsCheckFailIfRoutesNotGiven(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	Settings = cork.Pop()
	if _, err := Get("somepath"); err == nil {
		t.Error("Settings check should have errored since routes were not set properly")
	}
}

func TestSettingsCheckPass(t *testing.T) {
	Settings = cork.Pop()
	Settings.Routes = func(r cork.Routes) {}
	if _, err := Get("somepath"); err != nil {
		t.Error("Settings check should not have errored as it was set properly")
	}
}

func Routes(r cork.Routes) {
	r.Get("/foo", func(res http.ResponseWriter, req *http.Request) {})
}

func TestGetRoute(t *testing.T) {
	Settings = cork.Pop()
	Settings.Routes = Routes

	r, err := Get("/foo")

	if err != nil || r == nil {
		t.Error("Get test should have found route: '/foo'")
	}
}

func TestGetRouteFailed(t *testing.T) {
	Settings = cork.Pop()
	Settings.Routes = Routes

	r, err := Get("/bar")

	if err == nil && r != nil {
		t.Error("Get test shouldn't have found route: '/foo'")
	}
}
