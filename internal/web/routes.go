package web

import (
	"github.com/aewens/anda-server/pkg/reading"
)

func Welcome(server *Server) *Response {
	return &Response{
		Error: false,
		Name:  "welcome",
		Data:  "Hello, world!",
	}
}

func GetEntries(server *Server) *Response {
	name := "get-entries"
	entries, err := reading.Entities(server.DB)

	if err != nil {
		return &Response{
			Error: true,
			Name:  name,
			Data:  err.Error(),
		}
	}

	return &Response{
		Error: false,
		Name:  name,
		Data:  &entries,
	}
}
