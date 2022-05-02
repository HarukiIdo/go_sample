package handler

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"sync"

	"github.com/HarukiIdo/go_sample/model"
	"github.com/HarukiIdo/go_sample/service"

	"github.com/go-chi/chi/v5"
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

// go:embed vite-project/dist
var assets embed.FS

func tryRead(requestedPath string, w http.ResponseWriter) error {
	log.Println(requestedPath)
	f, err := assets.Open(path.Join("vite-project/dist", requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	// ディレクトリならエラー
	stat, _ := f.Stat()
	if stat.IsDir() {
		return errors.New("path is dir")
	}

	// MIMEタイプを設定
	ext := filepath.Ext(requestedPath)
	var contentType string
	if m := mime.TypeByExtension(ext); m != "" {
		contentType = m
	} else {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	io.Copy(w, f)
	return nil
}

func notFoundHanler(w http.ResponseWriter, r *http.Request) {
	// まずリクエストされた通りにファイルを探索
	err := tryRead(r.URL.Path, w)
	if err != nil {
		return
	}
	// 見つからなければindex.htmlを返す
	err = tryRead("index.html", w)
	if err != nil {
		panic(err)
	}
}

func NewHandler() http.Handler {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		})
	})

	// シングルページアプリケーションを配布するハンドラー
	router.NotFound(notFoundHanler)

	return router
}
