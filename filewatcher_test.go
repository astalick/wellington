package wellington

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatch_rebuild(t *testing.T) {
	tdir, err := ioutil.TempDir(os.TempDir(), "testwatch_")
	if err != nil {
		t.Fatal(err)
	}
	tfile := filepath.Join(tdir, "_new.scss")
	fh, err := os.Create(tfile)
	if err != nil {
		t.Fatal(err)
	}

	w := NewWatcher()
	w.Dirs = []string{tdir}
	w.PartialMap.AddRelation("tswif", tfile)
	err = w.Watch()
	if err != nil {
		t.Fatal(err)
	}
	rebuildChan = make(chan []string, 1)
	done := make(chan bool, 1)
	go func(t *testing.T) {
		select {
		case <-rebuildChan:
			done <- true
		case <-time.After(250 * time.Millisecond):
			done <- false
		}
		done <- true
	}(t)
	fh.WriteString("boom")
	success := <-done
	if !success {
		t.Fatal("Timeout waiting for rebuild")
	}

}

func TestWatch(t *testing.T) {
	w := NewWatcher()
	err := w.Watch()
	if err == nil {
		t.Error("No errors thrown for nil directories")
	}
	w.FileWatcher.Close()

	watcherChan = make(chan string, 1)
	w = NewWatcher()
	w.Dirs = []string{"test"}
	err = w.Watch()

	// Test file creation event
	go func() {
		select {
		case <-watcherChan:
			break
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("timeout %d\n", len(watcherChan))
			t.Error("Timeout without creating file")
		}
	}()

	testFile := "test/watchfile.lock"
	f, err := os.OpenFile(testFile, os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		// Give time for filesystem to sync before deleting file
		time.Sleep(50 * time.Millisecond)
		os.Remove(testFile)
		f.Close()
	}()
	if err != nil {
		t.Fatalf("creating test file failed: %s", err)
	}
	f.Sync()

	// Test file modification event
	go func() {
		select {
		case <-watcherChan:
			break
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("timeout %d\n", len(watcherChan))
			t.Error("Timeout without detecting write")
		}
	}()

	f.WriteString("data")
	f.Sync()

}

func TestRebuild(t *testing.T) {
	w := NewWatcher()
	err := w.rebuild("file/event")

	if e := fmt.Sprintf("build args are nil"); e != err.Error() {
		t.Errorf("wanted: %s\ngot: %s", e, err)
	}
}
