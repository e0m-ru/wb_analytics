package frontend

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/e0m-ru/wb_analitics/config"
)

type Server struct {
	cfg    *config.Config
	router *http.ServeMux
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		cfg:    cfg,
		router: http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Статические файлы
	staticDir := filepath.Join("frontend", "static")
	s.router.Handle("/", http.FileServer(http.Dir(staticDir)))

	// Проксирование API запросов
	s.router.HandleFunc("/api/", s.reverseProxy)
}

func (s *Server) reverseProxy(w http.ResponseWriter, r *http.Request) {
	// Создаем новый запрос к API серверу
	proxyReq, err := http.NewRequest(r.Method, s.cfg.APIBaseURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error contacting API server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем ответ
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
