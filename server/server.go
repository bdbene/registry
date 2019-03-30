package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/bdbene/registry/storage"
	"github.com/gorilla/mux"
)

// RestServer exposes a REST API
type RestServer struct {
	port    string
	router  *mux.Router
	storage storage.Storage
}

type schema struct {
	Name   string
	Schema string
}

// NewServer creates a rest server that exposes a REST API for clients.
func NewServer(configs *ServerConfig, storage storage.Storage) (*RestServer, error) {
	if configs == nil {
		return nil, &ServerError{"No server configurations recieved."}
	}

	router := mux.NewRouter().StrictSlash(true)
	server := &RestServer{configs.Port, router, storage}

	router.HandleFunc("/schemas/{schema}", server.getSchemaVersions).Methods("GET")
	router.HandleFunc("/schemas/{schema}/version/{version}", server.getSchema).Methods("GET")
	router.HandleFunc("/create", server.createSchema).Methods("POST")

	return server, nil
}

// Listen starts the server.
func (server *RestServer) Listen() {
	port := server.port

	log.Printf("Running server on port %s.", port)
	log.Fatal(http.ListenAndServe(":"+port, server.router))
}

func (server *RestServer) getSchemaVersions(writer http.ResponseWriter, request *http.Request) {
	url, err := url.Parse(request.URL.Path)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	schemaName := path.Base(url.String())
	versions, err := server.storage.LookupVersions(schemaName)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Schema not found.", http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, versions)
}

func (server *RestServer) getSchema(writer http.ResponseWriter, request *http.Request) {
	url, err := url.Parse(request.URL.Path)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	eles := strings.Split(url.String(), "/")
	name := eles[2]
	version := eles[4]

	schema, err := server.storage.Lookup(name, version)
	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Failed to find schema.", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, schema)
}

func (server *RestServer) createSchema(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var s schema

	err := decoder.Decode(&s)
	if err != nil {
		log.Printf("Bad request: %s\n", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = server.storage.Propose(s.Name, s.Schema)
	if err != nil {
		log.Printf("Storage failure: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: raft consensus

	err = server.storage.Commit(s.Name)
	if err != nil {
		log.Printf("Storage failure: %s\n", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
