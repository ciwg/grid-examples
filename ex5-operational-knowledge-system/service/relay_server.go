package service

import (
	"io"
	"net/http"
	"strings"
)

type RelayServer struct {
	relay *Relay
}

func NewRelayServer(relay *Relay) *RelayServer {
	return &RelayServer{relay: relay}
}

// Intent: Expose the dedicated remote relay surface as feed-plus-blob routes
// only, so remote transport stays separate from the local embodiment adapter.
// Source: DI-rovik
func (server *RelayServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(relayRoutePrefix+"/meta", server.handleMeta)
	mux.HandleFunc(relayRoutePrefix+"/feed/publish", server.handlePublish)
	mux.HandleFunc(relayRoutePrefix+"/feed/pull", server.handlePull)
	mux.HandleFunc(relayRoutePrefix+"/blobs/", server.handleBlob)
	return mux
}

func (server *RelayServer) handleMeta(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(writer, http.StatusOK, server.relay.Meta())
}

func (server *RelayServer) handlePublish(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 8*1024*1024)
	var batch RelayFeedBatch
	if err := decodeJSONBody(request, &batch); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := server.relay.Publish(batch)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if len(result.MissingBlobCIDs) > 0 {
		writeJSON(writer, http.StatusConflict, result)
		return
	}
	writeJSON(writer, http.StatusCreated, result)
}

func (server *RelayServer) handlePull(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	request.Body = http.MaxBytesReader(writer, request.Body, 256*1024)
	var relayRequest RelayFeedRequest
	if err := decodeJSONBody(request, &relayRequest); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	batch, err := server.relay.Pull(relayRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(writer, http.StatusOK, batch)
}

func (server *RelayServer) handleBlob(writer http.ResponseWriter, request *http.Request) {
	cid := strings.Trim(strings.TrimPrefix(request.URL.Path, relayRoutePrefix+"/blobs/"), "/")
	if cid == "" {
		http.NotFound(writer, request)
		return
	}
	switch request.Method {
	case http.MethodGet:
		body, err := server.relay.Blob(cid)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		writeNoStoreHeaders(writer)
		writer.Header().Set("Content-Type", "application/octet-stream")
		if _, err := writer.Write(body); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPut:
		request.Body = http.MaxBytesReader(writer, request.Body, maxEvidenceAttachmentBytes)
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err := server.relay.StoreBlob(cid, body); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, map[string]any{"cid": cid, "stored": true})
	default:
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
	}
}
