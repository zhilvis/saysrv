package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
)

type speech struct {
	voice  string
	format string
}

func newSpeech(voice string) speech {
	if voice == "" {
		voice = "cellos"
	}
	return speech{voice, "m4af"}
}

func tmpFile() string {
	f, e := ioutil.TempFile(os.TempDir(), "saysrv")
	if e != nil {
		log.Panic(e)
	}
	f.Close()
	os.Remove(f.Name())
	return f.Name() + ".m4a"
}

func (s speech) speak(text string) []byte {
	f := tmpFile()
	c := exec.Command("say", "-v", s.voice, text, "--file-format="+s.format, "-o", f)
	c.Start()
	c.Process.Wait()

	b, e := ioutil.ReadFile(f)
	if e != nil {
		log.Panic(e)
	}
	os.Remove(f)
	return b
}

func getQ(vals url.Values, key string) string {
	if len(vals[key]) > 0 {
		return vals[key][0]
	}
	return ""
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		if path[1] == "speak" {
      log.Println("REQUEST - ", r)
			text := path[3]
			voice := path[2]
			w.Header().Add("Content-Type", "audio/mp4")
			w.Write(newSpeech(voice).speak(text))
		}
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	log.Fatal(http.ListenAndServe(":60222", nil))
}
