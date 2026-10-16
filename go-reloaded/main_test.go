package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunWritesProcessedFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "in.txt")
	output := filepath.Join(dir, "out.txt")

	if err := os.WriteFile(input, []byte("1E (hex) files were added"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{input, output}); err != nil {
		t.Fatalf("run a échoué : %v", err)
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if want := "30 files were added"; string(got) != want {
		t.Errorf("obtenu %q, attendu %q", got, want)
	}
}

func TestRunInvalidArguments(t *testing.T) {
	if err := run([]string{"only-one"}); err == nil {
		t.Error("une erreur était attendue avec un seul argument")
	}
	if err := run([]string{"missing.txt", "out.txt"}); err == nil {
		t.Error("une erreur était attendue avec un fichier d'entrée absent")
	}
}
