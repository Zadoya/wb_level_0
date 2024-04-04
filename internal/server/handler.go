package server

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"net/http"
	"wb_level_0/internal/server/orderservice"
)

type httpHandler struct {
	router	*mux.Router
	service *orderservice.OrderService
}

func newHttpHandler(service *orderservice.OrderService) *httpHandler {
	h := &httpHandler{
		router:  mux.NewRouter(),
		service: service,
	}
	h.configureRouter()
	return h
}
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *httpHandler) configureRouter() {
	h.router.HandleFunc("/orders/{id}", h.GetOrderByUID).Methods("GET")
	h.router.HandleFunc("/", h.HandleStart).Methods("GET")
}

func (s *httpHandler) HandleStart(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../html/index.html")
}

func (s *httpHandler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["id"]
	response, err := s.service.GetByUid(uid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to find order with given uid")
	} else {
		respondWithJSON(w, http.StatusOK, response)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	formatingResponse(&response)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func formatingResponse(response *[]byte) {
	tab := []byte{}
	(*response) = append([]byte("ORDER:"), (*response)...)
	for i := 0; i < len(*response); i++ {	
		switch (*response)[i] {
		case '{':
			tab = append(tab, []byte("    ")...)
			(*response)[i] = '\n'
			(*response) = append((*response)[:i + 1], append(tab, (*response)[i + 1:]...)...)
		case '}':
			tab = tab[:len(tab) - 4]
			(*response)[i] = '\n'
			(*response) = append((*response)[:i + 1], append(tab, (*response)[i + 1:]...)...)
		case ',':
			(*response)[i] = '\n'
			(*response) = append((*response)[:i + 1], append(tab, (*response)[i + 1:]...)...)
		}
	}
}

