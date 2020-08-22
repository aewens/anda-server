package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/aewens/anda-server/internal/storage"
	"github.com/aewens/anda-server/pkg/core"
)

type handler func(http.ResponseWriter, *http.Request)

type routeMap map[string]handler
type routeMethodMap map[string]routeMap

type Server struct {
	Config *core.Config
	DB     *sql.DB
	Router *mux.Router
	Routes routeMethodMap
}

type Response struct {
	Error bool        `json:"err"`
	Name  string      `json:"name"`
	Data  interface{} `json:"data"`
}

type callback func(*Server) *Response

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func routeWrapper(res *Response) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err.Error())
		}
	}
}

func (server *Server) AddRoute(method string, route string, callback callback) {
	handle := routeWrapper(callback(server))
	server.Routes[method][route] = handle
	server.Router.HandleFunc(route, handle).Methods(method)
}

func Create(config *core.Config) (*Server, error) {
	db, err := storage.OpenPostgreSQL(config)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Use(jsonMiddleware)

	routes := make(routeMethodMap)
	routes["GET"] = make(routeMap)
	routes["POST"] = make(routeMap)
	routes["PUT"] = make(routeMap)
	routes["DELETE"] = make(routeMap)

	server := &Server{
		Config: config,
		DB:     db,
		Router: router,
		Routes: routes,
	}

	return server, nil
}

func (server *Server) Start() {
	defer server.DB.Close()

	// Start HTTP server
	httpPort := fmt.Sprintf(":%d", server.Config.ApiPort)
	fmt.Printf("Starting server on port %s\n", httpPort)
	log.Println(http.ListenAndServe(httpPort, server.Router))
}
