package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/HarukiIdo/go_sample/model"
	"github.com/HarukiIdo/go_sample/service"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{
		svc: svc,
	}
}

var mutex = &sync.RWMutex{}

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	comments := make([]model.Comment, 0, 100)

	switch r.Method {
	case http.MethodGet:
		mutex.RLock()
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			return
		}
		mutex.RUnlock()

	case http.MethodPost:
		var c model.Comment
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
			return
		}
		mutex.Lock()
		comments = append(comments, c)
		mutex.Unlock()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "created"}`))

	default:
		http.Error(w, `{"status":"premits only Get or Post"}`, http.StatusMethodNotAllowed)
	}
}

