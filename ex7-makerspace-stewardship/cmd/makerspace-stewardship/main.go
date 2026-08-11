package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/service"
)

func main() {
	runtimeRoot := flag.String("runtime-root", ".makerspace-stewardship", "directory for append-only local evidence")
	allowEmptyRecognition := flag.Bool("allow-empty-recognition", false, "start with no recognized public keys")
	flag.Parse()
	policy, err := service.LoadRecognitionPolicy(filepath.Join(*runtimeRoot, "recognition.json"), *allowEmptyRecognition)
	if err != nil {
		log.Fatal(err)
	}
	app, err := service.NewPersistentRecordApp(*runtimeRoot, policy)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Makerspace Stewardship listening on http://127.0.0.1:7037")
	log.Fatal(http.ListenAndServe("127.0.0.1:7037", app.Handler()))
}
