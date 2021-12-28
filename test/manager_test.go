package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sessionManager"
	_ "github.com/sessionManager/store"
)

func TestSessionStart(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(`GET`, "/", nil)
	cookieName := "gosessionid"
	sm, err := sessionManager.NewManager("store", cookieName, 360)
	if err != nil {
		t.Errorf("NewManager err")
	}

	session := sm.SessionStart(w, r)

	// session store内のデータ保存機能テスト
	err = session.Set("name", "test")
	if err != nil {
		t.Errorf("set session store err")
	}
	name := session.Get("name")
	if name != "test" {
		t.Errorf("get session store err")
	}
	// set-cookieの中身をテスト
	setCookie := w.Header().Get("Set-Cookie")
	ary := strings.Split(setCookie, "; ")
	for idx, cookie := range ary {
		cary := strings.Split(cookie, "=")
		if idx == 0 && cary[0] != cookieName {
			t.Errorf("cookie name err")
		}

	}
}

func TestSessionDelete(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(`GET`, "/", nil)
	cookieName := "gosessionid"
	sm, err := sessionManager.NewManager("store", cookieName, 360)
	if err != nil {
		t.Errorf("NewManager err")
	}

	session := sm.SessionStart(w, r)

	// 同じクライアントから再度リクエストした時のsessionid
	r1, _ := http.NewRequest(`GET`, "/", nil)
	r1.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	newSession := sm.SessionStart(w, r1)
	if session.SessionID() != newSession.SessionID() {
		t.Errorf("no same sessionid")
	}

	// 異なるクライアントからリクエストした時のsessionid
	sm.SessionDestroy(w, r1)
	r2, _ := http.NewRequest(`GET`, "/", nil)
	r2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()
	newSession = sm.SessionStart(w, r2)
	if session.SessionID() == newSession.SessionID() {
		t.Errorf("same sessionid")
	}
}