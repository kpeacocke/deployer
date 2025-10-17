package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTarGz(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a test tar.gz file
	tarPath := filepath.Join(tmpDir, "test.tar.gz")
	f, err := os.Create(tarPath)
	if err != nil {
		t.Fatalf("failed to create tar file: %v", err)
	}

	gzWriter := gzip.NewWriter(f)
	tarWriter := tar.NewWriter(gzWriter)

	// Add a test file to the archive
	testContent := []byte("Hello, World!")
	header := &tar.Header{
		Name: "testfile.txt",
		Mode: 0o644,
		Size: int64(len(testContent)),
	}
	if err = tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err = tarWriter.Write(testContent); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}

	tarWriter.Close()
	gzWriter.Close()
	f.Close()

	// Extract the archive
	extractDir := filepath.Join(tmpDir, "extracted")
	if err = ExtractTarGz(tarPath, extractDir); err != nil {
		t.Fatalf("failed to extract tar.gz: %v", err)
	}

	// Verify extracted file
	extractedFile := filepath.Join(extractDir, "testfile.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("extracted content mismatch: got %q, want %q", string(content), string(testContent))
	}
}

func TestExtractZip(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a test zip file
	zipPath := filepath.Join(tmpDir, "test.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create zip file: %v", err)
	}

	zipWriter := zip.NewWriter(f)

	// Add a test file to the archive
	testContent := []byte("Hello, Zip!")
	writer, err := zipWriter.Create("testfile.txt")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	if _, err = writer.Write(testContent); err != nil {
		t.Fatalf("failed to write zip content: %v", err)
	}

	zipWriter.Close()
	f.Close()

	// Extract the archive
	extractDir := filepath.Join(tmpDir, "extracted")
	if err = ExtractZip(zipPath, extractDir); err != nil {
		t.Fatalf("failed to extract zip: %v", err)
	}

	// Verify extracted file
	extractedFile := filepath.Join(extractDir, "testfile.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("extracted content mismatch: got %q, want %q", string(content), string(testContent))
	}
}

func TestParseChecksums(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a test checksums file
	checksumsPath := filepath.Join(tmpDir, "checksums.txt")
	testData := `abc123def456  file1.tar.gz
789ghi012jkl  file2.zip
mno345pqr678  file3.bin
`
	if err := os.WriteFile(checksumsPath, []byte(testData), 0o644); err != nil {
		t.Fatalf("failed to write checksums file: %v", err)
	}

	// Parse checksums
	checksums, err := ParseChecksums(checksumsPath)
	if err != nil {
		t.Fatalf("failed to parse checksums: %v", err)
	}

	// Verify parsed data
	expected := map[string]string{
		"file1.tar.gz": "abc123def456",
		"file2.zip":    "789ghi012jkl",
		"file3.bin":    "mno345pqr678",
	}

	for filename, expectedHash := range expected {
		gotHash, ok := checksums[filename]
		if !ok {
			t.Errorf("missing checksum for %s", filename)
			continue
		}
		if gotHash != expectedHash {
			t.Errorf("checksum mismatch for %s: got %s, want %s", filename, gotHash, expectedHash)
		}
	}
}

func TestVerifyFileSHA256(t *testing.T) {
	// Create a temporary directory for test
	tmpDir := t.TempDir()

	// Create a test file with known content
	testContent := []byte("test content for checksum")
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, testContent, 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Calculate expected checksum
	h := sha256.New()
	h.Write(testContent)
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Test with correct checksum
	if err := VerifyFileSHA256(testFile, expectedHash); err != nil {
		t.Errorf("verification failed with correct checksum: %v", err)
	}

	// Test with incorrect checksum
	wrongHash := "0000000000000000000000000000000000000000000000000000000000000000"
	if err := VerifyFileSHA256(testFile, wrongHash); err == nil {
		t.Error("verification should fail with incorrect checksum")
	}
}
