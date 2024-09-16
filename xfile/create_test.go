package fileutil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCreate_Successful(t *testing.T) {
	tmpDir, cleanUpFn := createTmpDir(t)
	defer cleanUpFn()

	wantSize := 10

	gotName, err := Create(tmpDir, wantSize)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	f, err := os.Stat(gotName)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if int64(wantSize) != f.Size() {
		t.Fatalf("Create file size does not match want size, want: %d got: %d", wantSize, f.Size())
	}
}

func TestCreate_PathShouldBeAbs(t *testing.T) {
	tmpDir, cleanUpFn := createTmpDir(t)
	defer cleanUpFn()

	_, err := Create(filepath.Base(tmpDir), 10)

	if err == nil {
		t.Fatal("Expected error got nil")
	}
}

func TestCreatePartialMatch_Success(t *testing.T) {
	tmpDir, cleanUpFn := createTmpDir(t)
	defer cleanUpFn()

	wantSize := 10

	files, err := CreatePartialMatch(tmpDir, 10)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("Got unexpected amount of files want: 2 got: %d", len(files))
	}

	// Compare the first bytes
	f1, err := os.Open(files[0])
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	defer f1.Close()
	f1FirstBytes := make([]byte, wantSize-numOfDiffBytes)
	if _, err := io.ReadFull(f1, f1FirstBytes); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	f2, err := os.Open(files[1])
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	defer f2.Close()
	f2FirstBytes := make([]byte, wantSize-numOfDiffBytes)
	if _, err := io.ReadFull(f2, f2FirstBytes); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if !reflect.DeepEqual(f1FirstBytes, f2FirstBytes) {
		t.Fatalf("First bytes are not equal")
	}

	// Compare the last bytes are different.
	f1LastBytes := make([]byte, numOfDiffBytes)
	_, err = f1.ReadAt(f1LastBytes, int64(wantSize-numOfDiffBytes))
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	f2LastBytes := make([]byte, numOfDiffBytes)
	_, err = f2.ReadAt(f2LastBytes, int64(wantSize-numOfDiffBytes))
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if reflect.DeepEqual(f1LastBytes, f2LastBytes) {
		t.Fatalf("Last bytes should be different")
	}
}

func TestCreatePartialMatch_PathShouldBeAbs(t *testing.T) {
	tmpDir, cleanUpFn := createTmpDir(t)
	defer cleanUpFn()

	_, err := CreatePartialMatch(filepath.Base(tmpDir), 10)
	if err == nil {
		t.Fatalf("Expected error got nil")
	}
}

func createTmpDir(t *testing.T) (string, func()) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to create tmp dir for testdata: %v", err)
	}

	return tmpDir, func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Logf("Failed to clean up test %s: %v", t.Name(), err)
		}
	}
}
