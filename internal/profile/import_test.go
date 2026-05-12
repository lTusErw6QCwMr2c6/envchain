package profile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/profile"
)

func newImportStore(t *testing.T) profile.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := profile.NewStore(filepath.Join(dir, "profiles"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func TestImportProfile_Dotenv_NewProfile(t *testing.T) {
	st := newImportStore(t)
	input := "FOO=bar\nBAZ=qux\n"
	err := profile.ImportProfile(st, "myapp", strings.NewReader(input), profile.ImportFormatDotenv, profile.ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := st.Load("myapp")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(p.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(p.Vars))
	}
}

func TestImportProfile_Export_StripsPrefix(t *testing.T) {
	st := newImportStore(t)
	input := "export FOO=hello\nexport BAR=world\n"
	err := profile.ImportProfile(st, "myapp", strings.NewReader(input), profile.ImportFormatExport, profile.ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, _ := st.Load("myapp")
	if len(p.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(p.Vars))
	}
}

func TestImportProfile_Overwrite(t *testing.T) {
	st := newImportStore(t)
	initial := "FOO=original\n"
	profile.ImportProfile(st, "myapp", strings.NewReader(initial), profile.ImportFormatDotenv, profile.ImportOptions{})

	update := "FOO=updated\n"
	err := profile.ImportProfile(st, "myapp", strings.NewReader(update), profile.ImportFormatDotenv, profile.ImportOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, _ := st.Load("myapp")
	if p.Vars[0].Value != "updated" {
		t.Errorf("expected 'updated', got %q", p.Vars[0].Value)
	}
}

func TestImportProfile_NoOverwrite_Preserves(t *testing.T) {
	st := newImportStore(t)
	profile.ImportProfile(st, "myapp", strings.NewReader("FOO=original\n"), profile.ImportFormatDotenv, profile.ImportOptions{})
	profile.ImportProfile(st, "myapp", strings.NewReader("FOO=new\n"), profile.ImportFormatDotenv, profile.ImportOptions{Overwrite: false})
	p, _ := st.Load("myapp")
	if p.Vars[0].Value != "original" {
		t.Errorf("expected 'original', got %q", p.Vars[0].Value)
	}
}

func TestImportProfile_InvalidName(t *testing.T) {
	st := newImportStore(t)
	err := profile.ImportProfile(st, "bad name!", strings.NewReader("A=1\n"), profile.ImportFormatDotenv, profile.ImportOptions{})
	if err == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestImportProfile_MalformedLine(t *testing.T) {
	st := newImportStore(t)
	err := profile.ImportProfile(st, "myapp", strings.NewReader("NOEQUALS\n"), profile.ImportFormatDotenv, profile.ImportOptions{})
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestImportProfile_SkipsComments(t *testing.T) {
	st := newImportStore(t)
	input := "# comment\nFOO=bar\n"
	err := profile.ImportProfile(st, "myapp", strings.NewReader(input), profile.ImportFormatDotenv, profile.ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, _ := st.Load("myapp")
	if len(p.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(p.Vars))
	}
}

var _ = os.DevNull // suppress unused import
