package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/smart-cs/scheduler-backend/models"
	"github.com/smart-cs/scheduler-backend/schedules"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const logFormat = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Method}} {{.Path}}\n"

// Server runs the backend server.
type Server struct {
	Middleware      *negroni.Negroni
	ScheduleCreator schedules.ScheduleCreator
	AutoCompleter   schedules.AutoCompleter
}

// StandardResponse is the default response from the server.
type StandardResponse struct {
	OK     bool        `json:"OK"`
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// NewServer constructs a Server to listen on the given port.
func NewServer() Server {
	server := Server{
		Middleware:      negroni.New(),
		ScheduleCreator: schedules.NewScheduleCreator(),
		AutoCompleter:   schedules.NewAutoCompleter(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/schedules", server.SchedulesHandler).
		Methods("GET").
		Queries("courses", "{courses}")
	router.HandleFunc("/autocomplete", server.AutocompleteHandler).
		Methods("GET").
		Queries("text", "{text}")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	logger := negroni.NewLogger()
	logger.SetDateFormat(time.Stamp)
	logger.SetFormat(logFormat)
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	server.Middleware.Use(logger)
	server.Middleware.Use(cors)
	server.Middleware.Use(negroni.NewRecovery())
	server.Middleware.UseHandler(router)
	return server
}

// Run starts the server on $PORT or 8080 by default.
func (s *Server) Run() {
	s.Middleware.Run()
}

// SchedulesHandler handles the schedule endpoint
func (s *Server) SchedulesHandler(w http.ResponseWriter, r *http.Request) {
	courses := strings.Split(r.URL.Query().Get("courses"), ",")
	term := r.URL.Query().Get("term")
	lecturesOnly := r.URL.Query().Get("lectures_only")
	if term == "" {
		term = "1-2"
	}
	selectOptions := schedules.ScheduleSelectOptions{
		Term: term,
		SelectLabsAndTutorials: lecturesOnly == "false",
	}

	schedules := s.ScheduleCreator.Create(courses, selectOptions)
	if schedules == nil {
		// Make schedules into an array of size 0 for JSON serialization
		schedules = make([]models.Schedule, 0)
	}
	s.respOK(w, schedules)
}

// AutocompleteHandler handles the autocomplete endpoint
func (s *Server) AutocompleteHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	completes := s.AutoCompleter.CoursesWithPrefix(text)
	if completes == nil {
		// Make completes into an array of size 0 for JSON serialization
		completes = make([]string, 0)
	}
	s.respOK(w, completes)
}

func (s *Server) respOK(w http.ResponseWriter, body interface{}) {
	r := StandardResponse{
		OK:     true,
		Status: http.StatusOK,
		Body:   body,
	}

	j, err := json.Marshal(r)
	if err != nil {
		panic("can't marshal JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
