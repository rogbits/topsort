package lib

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Api struct {
	Logger *Logger
}

func NewApi(logger *Logger) *Api {
	a := new(Api)
	a.Logger = logger
	return a
}

func (api *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.Logger.Log("incoming", r.Method, r.RequestURI)
	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		api.Logger.Log(err)
		http.Error(w, "bad URI", http.StatusBadRequest)
		return
	}

	switch {
	case u.Path == "/health":
		if r.Method != http.MethodGet {
			http.Error(w, "unsupported method for /health", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("ok\n"))
	case u.Path == "/sort":
		if r.Method != http.MethodPost {
			http.Error(w, "unsupported method for /sort", http.StatusBadRequest)
			return
		}

		// decode input
		var err error
		var edges [][]string
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&edges)
		if err != nil {
			api.Logger.Log(err)
			http.Error(w, "error decoding edges input", http.StatusBadRequest)
			return
		}

		// build graph
		graph := NewGraph[string]()
		for _, edge := range edges {
			graph.AddEdge(edge[0], edge[1])
			if graph.Undirected {
				api.Logger.Log(errors.New("seeing an undirected graph"))
				http.Error(w, "seeing an undirected graph", http.StatusBadRequest)
				return
			}
		}

		// check if empty
		api.Logger.Log("graph size in request:", graph.Vertices.Size)
		if graph.Vertices.Size == 0 {
			api.Logger.Log(errors.New("seeing empty graph"))
			http.Error(w, "seeing empty graph", http.StatusBadRequest)
			return
		}

		// run top sort
		var response []string
		for graph.HasNextLevel() {
			vertices, err := graph.GetLevel()
			if err != nil {
				api.Logger.Log(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			response = append(response, vertices...)
		}

		// write response
		api.Logger.Log("sorted result", response)
		enc := json.NewEncoder(w)
		err = enc.Encode(response)
		if err != nil {
			api.Logger.Log(err)
			http.Error(w, "unexpected error while encoding response", http.StatusInternalServerError)
			return
		}
	default:
		http.NotFound(w, r)
	}
}
