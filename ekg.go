package ekg

import (
	"encoding/json"
	"github.com/adarqui/ekg-core-go"
	"net/http"
	"strings"
	"time"
    "fmt"
)

type Server struct {
	// serverThreadId
	// serverMetricStore
	store *ekg_core.Store
	io    *http.Server
}

type serverHandler struct {
	server *Server
	h      func(*Server, http.ResponseWriter, *http.Request)
}

func ForkServer(bind string) (*Server, error) {
	store := ekg_core.New()
	store.RegisterGCMetrics()
	return ForkServerWith(store, bind)
}

func ForkServerWith(store *ekg_core.Store, bind string) (*Server, error) {
	server := new(Server)

	io := &http.Server{
		Addr:           bind,
		Handler:        serverHandler{server, serve},
		MaxHeaderBytes: 1 << 20,
	}

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets", fs)

	server.io = io
	server.store = store

	server.store.RegisterCounter("ekg.server_timestamp_ms", getTimeMs)
	go server.io.ListenAndServe()

	return server, nil
}

func (sh serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve(sh.server, w, r)
}

func serve(server *Server, w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("/root/projects/go/src/github.com/adarqui/ekg-go/assets/")))
	url := r.URL.Path
	if strings.HasPrefix(url, "/index.html") {
		http.Redirect(w, r, "/assets/", 0)
	} else if strings.HasPrefix(url, "/assets/") {
		fs.ServeHTTP(w, r)
	} else {
		serveMetrics(server, w, r)
		//        http.NotFound(w, r)
	}
}

func serveMetrics(server *Server, w http.ResponseWriter, r *http.Request) {
	v := EncodeAll(server.store.SampleAll())
	switch r.Header.Get("Content-Type") {
	case "application/json":
		{
			js, err := json.Marshal(v)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, string(js))
		}
	case "application/xml":
		{
			fmt.Fprintf(w, "xml")
		}
	case "application/html":
		{
			fmt.Fprintf(w, "html")
		}
	default:
		{
			fmt.Fprintf(w, "error")
		}
	}
}

/*
func (server *Server) serveAll() {
//    metrics := server.store.SampleAll()
}

func (server *Server) serveOne(pathInfo string) {
}
*/

func (server *Server) GetStore() *ekg_core.Store {
	return server.store
}

func (server *Server) GetCounter(name string) *ekg_core.Counter {
	return server.store.CreateCounter(name)
}

func (server *Server) GetGauge(name string) *ekg_core.Gauge {
	return server.store.CreateGauge(name)
}

func (server *Server) GetLabel(name string) *ekg_core.Label {
	return server.store.CreateLabel(name)
}

func (server *Server) GetDistribution(name string) *ekg_core.Distribution {
	return server.store.CreateDistribution(name)
}

func getTimeMs(void interface{}) interface{} {
	return time.Now().UnixNano() / 1000000
}
