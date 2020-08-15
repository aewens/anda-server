package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"database/sql"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Handler func(http.ResponseWriter, *http.Request)
type routeMap map[string]Handler
type routeMethodMap map[string]routeMap

type Server struct {
	Config	*Config
	DB		*sql.DB
	Router	*mux.Router
	Routes	routeMethodMap
}

type callback func(*Server) Handler

func OpenDatabase(config *Config) (*sql.DB) {
	var conn []string
	conn = append(conn, fmt.Sprintf("host=%s", config.DBHost))
	conn = append(conn, fmt.Sprintf("port=%s", config.DBPort))
	conn = append(conn, fmt.Sprintf("user=%s", config.DBUser))
	conn = append(conn, fmt.Sprintf("password=%s", config.DBPswd))
	conn = append(conn, fmt.Sprintf("dbname=%s", config.DBName))
	conn = append(conn, "sslmode=disable")
	//connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUser,
	//	config.DBPswd, config.DBHost, config.DBPort, config.DBName)

	db, err := sql.Open("postgres", strings.Join(conn, " ")); if err != nil {
		log.Println(err)
		os.Exit(99)
	}

	return db
}

func CreateServer(config *Config) (*Server) {
	db := OpenDatabase(config)
	router := mux.NewRouter().StrictSlash(true)
	routes := make(routeMethodMap)
	routes["GET"] = make(routeMap)
	routes["POST"] = make(routeMap)
	routes["PUT"] = make(routeMap)
	routes["DELETE"] = make(routeMap)

	server := &Server{
		Config: config,
		DB: db,
		Router: router,
		Routes: routes,
	}

	return server
}

func (server *Server) AddRoute(method string, route string, callback callback) {
	handle := callback(server)
	server.Routes[method][route] = handle
	server.Router.HandleFunc(route, handle).Methods(method)
}

func (server Server) Start() {
	defer server.DB.Close()

	// Start HTTP server
	httpPort := fmt.Sprintf(":%d", server.Config.ApiPort)
	fmt.Printf("Starting server on port %s\n", httpPort)
	log.Println(http.ListenAndServe(httpPort, server.Router))
}
