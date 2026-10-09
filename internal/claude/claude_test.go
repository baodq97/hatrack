package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	root := t.TempDir()
	return &Store{
		ConfigDir:  filepath.Join(root, ".claude"),
		GlobalJSON: filepath.Join(root, ".claude.json"),
		Dir:        filepath.Join(root, ".hatrack", "claude"),
	}
}

// signIn fakes what `claude auth login` leaves in a config dir.
func signIn(t *testing.T, credsPath, jsonPath, email, token string) {
	t.Helper()
	write(t, credsPath, `{"claudeAiOauth":{"refreshToken":"`+token+`"}}`)
	m := map[string]any{"numStartups": 7, "projects": map[string]any{"/x": map[string]any{"a": "<b>"}}}
	if b, err := os.ReadFile(jsonPath); err == nil {
		m = nil
		must(t, json.Unmarshal(b, &m))
	}
	m["oauthAccount"] = map[string]any{"accountUuid": "uuid-" + email, "organizationUuid": "org", "emailAddress": email}
	b, _ := json.Marshal(m)
	write(t, jsonPath, string(b))
}

func activeToken(t *testing.T, s *Store) string {
	t.Helper()
	var c struct{ ClaudeAiOauth struct{ RefreshToken string } }
	b, err := os.ReadFile(filepath.Join(s.ConfigDir, credsFile))
	must(t, err)
	must(t, json.Unmarshal(b, &c))
	return c.ClaudeAiOauth.RefreshToken
}

func TestSwitchKeepsRefreshedTokensAndOtherState(t *testing.T) {
	s := newStore(t)
	signIn(t, filepath.Join(s.ConfigDir, credsFile), s.GlobalJSON, "a@x.com", "a1")

	// add b without touching the active login
	name, err := s.Add("", func(dir string) error {
		signIn(t, filepath.Join(dir, credsFile), filepath.Join(dir, ".claude.json"), "b@x.com", "b1")
		return nil
	})
	must(t, err)
	if name != "b@x.com" || activeToken(t, s) != "a1" {
		t.Fatalf("add: name=%q active=%q", name, activeToken(t, s))
	}

	// a was never saved: switching must park it, not lose it
	must(t, s.Use("b@x.com"))
	if activeToken(t, s) != "b1" {
		t.Fatalf("use b: active=%q", activeToken(t, s))
	}

	// claude refreshes b's token while b is active
	write(t, filepath.Join(s.ConfigDir, credsFile), `{"claudeAiOauth":{"refreshToken":"b2"}}`)
	must(t, s.Use("b@x.com")) // re-selecting the active account keeps the fresh token
	if activeToken(t, s) != "b2" {
		t.Fatalf("reuse b: active=%q", activeToken(t, s))
	}
	must(t, s.Use("a@x.com"))
	must(t, s.Use("b@x.com"))
	if activeToken(t, s) != "b2" {
		t.Fatalf("round trip lost refreshed token: %q", activeToken(t, s))
	}

	var m map[string]any
	b, _ := os.ReadFile(s.GlobalJSON)
	must(t, json.Unmarshal(b, &m))
	if m["numStartups"] != float64(7) || m["projects"].(map[string]any)["/x"].(map[string]any)["a"] != "<b>" {
		t.Fatalf("other .claude.json keys changed: %s", b)
	}
	if m["oauthAccount"].(map[string]any)["emailAddress"] != "b@x.com" {
		t.Fatalf("oauthAccount not swapped: %s", b)
	}

	ps, err := s.List()
	must(t, err)
	if len(ps) != 2 || ps[0].Name != "a@x.com" || ps[0].Active || !ps[1].Active {
		t.Fatalf("list: %+v", ps)
	}
}

func TestSaveRefusesToOverwriteAnotherAccount(t *testing.T) {
	s := newStore(t)
	signIn(t, filepath.Join(s.ConfigDir, credsFile), s.GlobalJSON, "a@x.com", "a1")
	_, err := s.Save("work")
	must(t, err)
	signIn(t, filepath.Join(s.ConfigDir, credsFile), s.GlobalJSON, "b@x.com", "b1")
	if _, err := s.Save("work"); err == nil {
		t.Fatal("saved b over a's profile")
	}
	signIn(t, filepath.Join(s.ConfigDir, credsFile), s.GlobalJSON, "a@x.com", "a2")
	if _, err := s.Save("work2"); err == nil {
		t.Fatal("saved a second copy of a")
	}
}

func TestRejectsPathNames(t *testing.T) {
	s := newStore(t)
	for _, n := range []string{"../x", "a/b", `a\b`, ".hidden", ""} {
		if err := s.Use(n); err == nil {
			t.Errorf("Use(%q) accepted", n)
		}
		if err := s.Remove(n); err == nil {
			t.Errorf("Remove(%q) accepted", n)
		}
	}
}

func TestNotSignedIn(t *testing.T) {
	s := newStore(t)
	if _, err := s.Save(""); err != ErrNotSignedIn {
		t.Fatalf("got %v", err)
	}
}

func write(t *testing.T, path, data string) {
	t.Helper()
	must(t, os.MkdirAll(filepath.Dir(path), 0o700))
	must(t, os.WriteFile(path, []byte(data), 0o600))
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
