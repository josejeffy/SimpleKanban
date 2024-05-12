package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

//go:embed ui
var ui embed.FS

func update(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	err := os.WriteFile("tasks.json", data, fs.FileMode(os.O_CREATE))
	if err != nil {
		fmt.Println(err)
	}
}

func load(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("tasks.json")
	if err != nil {
		return
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	fmt.Fprint(w, string(data))
}

func main() {
	_, err := os.Stat("tasks.json")
	if err != nil {
		os.Create("tasks.json")
	}

	http.Handle("/", http.FileServer(http.FS(ui)))

	http.HandleFunc("/update", update)
	http.HandleFunc("/load", load)

	fmt.Println("Running on http://localhost:8080/ui")
	http.ListenAndServe("127.0.0.1:8080", nil)
}
