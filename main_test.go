// hard integration test =))
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

type goroutines struct {
	All  int `json:"Goroutines"`
	Exec int `json:"GoroutinesInExecutor"`
}

func getExecGoroutines() (*goroutines, error) {
	resp, err := http.Get(`http://localhost:5678/debug/vars`)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	out := &goroutines{}
	if err = json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

func isUp() (up bool) {
	for i := 0; i < 5; i++ {
		time.Sleep(200 * time.Millisecond)
		_, err := getExecGoroutines()
		if err == nil {
			return true
		}
	}
	return false
}

func TestByFiles(t *testing.T) {
	prefix := `./testdata/`

	cmd := exec.Command(
		`go`,
		`run`,
		`-tags=expvar`,
		`main.go`,
		`-type`,
		`file`,
		`-k`,
		`100`,
	)

	writer, err := cmd.StdinPipe()
	if err != nil {
		t.Errorf(`can't open stdin to process`, err)
	}

	var (
		result = make(chan string)

		up               = make(chan *struct{})
		done             = make(chan *struct{})
		onHot            = make(chan *struct{})
		doneAfterTimeout = make(chan *struct{})
	)

	// run cmd
	go func() {
		close(up)
		out, err := cmd.CombinedOutput()
		if err != nil {
			result <- err.Error()
		} else {
			result <- string(out)
		}
	}()

	<-up
	if !isUp() {
		t.Error(`stats server shutdown, can't test of executed pool`)
	}

	// start write to stdin
	go func() {
		for i := 0; i < 1000; i++ {
			_, err := writer.Write([]byte(prefix + "file_with_go.txt\n"))
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
			if i == 100 {
				close(onHot)
			}
			_, err = writer.Write([]byte(prefix + "just_file.txt\n"))
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
		}
		close(done)
	}()

	runtime.Gosched()
	// for compile time and up stats
	<-onHot
	statsOnRun, err := getExecGoroutines()
	if err != nil {
		t.Error(err)
	} else {
		if statsOnRun.Exec > 100 || statsOnRun.Exec < 50 {
			t.Errorf(`expected exec goroutines for tasks <= 100 && > 50, got %d`, statsOnRun.Exec)
		}

	}
	<-done

	// sleep for a wait and check  keeping pause of stdin
	time.Sleep(500 * time.Millisecond)
	// check pool of goroutines equals == 0
	stats, err := getExecGoroutines()
	if err != nil {
		t.Error(err)
	} else if stats.Exec != 0 {
		t.Errorf(`expected exec goroutines for tasks == 0, got %d`, stats.Exec)
	}

	// start write to stdin
	go func() {
		for i := 0; i < 1000; i++ {
			_, err := writer.Write([]byte(prefix + "file_with_go.txt\n"))
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
		}
		close(doneAfterTimeout)
	}()
	<-doneAfterTimeout
	writer.Close()

	stdout := <-result
	if !strings.Contains(stdout, `Total: 140000`) {
		t.Error(`expected got total: 140000, got `, stdout)
	}
}

var goPage = []byte(`
<html>
<body>
	<h1>Go Go Go Go Go Go Go Go Go Go</h1>
	<div>gogogogogogo Go</div>
</body>
</html>
`)

var page = []byte(`
<html>
<body>
	<h1>Hello!</h1>
</body>
</html>
`)

func startTestServer() (string, func()) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, `hello`) {
			rw.Write(page)
		} else {
			rw.Write(goPage)
		}
	}))

	return srv.Listener.Addr().String(), srv.Close
}

func TestByUrls(t *testing.T) {

	addr, stopTestSrv := startTestServer()
	defer stopTestSrv()

	urlWithGo := []byte(fmt.Sprintf("http://%s/\n", addr))
	urlWithoutGo := []byte(fmt.Sprintf("http://%s/hello\n", addr))

	cmd := exec.Command(
		`go`,
		`run`,
		`-tags=expvar`,
		`main.go`,
		`-type`,
		`url`,
		`-k`,
		`20`,
	)

	writer, err := cmd.StdinPipe()
	if err != nil {
		t.Errorf(`can't open stdin to process`, err)
	}

	var (
		result = make(chan string)

		up               = make(chan *struct{})
		done             = make(chan *struct{})
		doneHalf         = make(chan *struct{})
		doneAfterTimeout = make(chan *struct{})
	)

	// run cmd
	go func() {
		close(up)
		out, err := cmd.CombinedOutput()
		if err != nil {
			result <- err.Error()
		} else {
			result <- string(out)
		}
	}()

	<-up
	if !isUp() {
		t.Error(`stats server shutdown, can't test of executed pool`)
	}

	// start write to stdin
	go func() {
		for i := 0; i < 50; i++ {
			_, err := writer.Write(urlWithGo)
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
			if i == 25 {
				close(doneHalf)
			}
			_, err = writer.Write(urlWithoutGo)
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
		}
		close(done)
	}()

	// for compile time and up stats
	<-doneHalf

	statsOnRun, err := getExecGoroutines()
	if err != nil {
		t.Error(err)
	} else {
		if statsOnRun.Exec != 20 {
			t.Errorf(`expected exec goroutines for tasks == 20, got %d`, statsOnRun.Exec)
		}
	}
	<-done

	// sleep for a wait and check  keeping pause of stdin
	time.Sleep(1 * time.Second)

	// check pool of goroutines equals == 0
	stats, err := getExecGoroutines()
	if err != nil {
		t.Error(err)
	} else if stats.Exec != 0 {
		t.Errorf(`expected exec goroutines for tasks == 0, got %d`, stats.Exec)
	}

	// start write to stdin
	go func() {
		for i := 0; i < 10; i++ {
			_, err := writer.Write(urlWithGo)
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
			_, err = writer.Write([]byte("https://golangxxx.org\n"))
			if err != nil {
				t.Errorf(`can't write to stdin`, err)
			}
		}
		close(doneAfterTimeout)
	}()
	<-doneAfterTimeout
	writer.Close()

	stdout := <-result
	if !strings.Contains(stdout, `Total: 660`) {
		t.Error(`expected got total: 660, got `, stdout)
	}
}
