package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/computerscienceiscool/grid-examples/ex7-makerspace-stewardship/service"
)

func main() {
	runtimeRoot := flag.String("runtime-root", ".makerspace-stewardship", "directory for append-only local evidence")
	listenAddress := flag.String("listen", "127.0.0.1:7037", "loopback address for this participant agent")
	allowEmptyRecognition := flag.Bool("allow-empty-recognition", false, "start with no recognized public keys")
	identityPassphraseStdin := flag.Bool("identity-passphrase-stdin", false, "read this participant agent's identity passphrase from standard input")
	flag.Parse()
	policy, err := service.LoadRecognitionPolicy(filepath.Join(*runtimeRoot, "recognition.json"), *allowEmptyRecognition)
	if err != nil {
		log.Fatal(err)
	}
	var app *service.App
	if *identityPassphraseStdin {
		passphrase, readErr := io.ReadAll(io.LimitReader(os.Stdin, 4097))
		if readErr != nil {
			log.Fatal(readErr)
		}
		identity, loadErr := service.LoadParticipantIdentity(*runtimeRoot, bytes.TrimSpace(passphrase))
		if loadErr != nil {
			log.Fatal(loadErr)
		}
		app, err = service.NewPersistentParticipantApp(*runtimeRoot, policy, identity)
	} else {
		app, err = service.NewPersistentRecordApp(*runtimeRoot, policy)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Makerspace Stewardship listening on http://%s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, app.Handler()))
}
