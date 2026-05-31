package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TESTES DE ENDPOINTS DA API REST

func TestAPI_PingEndpoint(t *testing.T) {
	mux := SetupRouter(nil, time.Now())

	// Cria uma requisição HTTP falsa simulando um cliente chamando o /ping
	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder() // Captura a resposta enviada pela API

	// Executa a chamada
	mux.ServeHTTP(w, req)

	// Asserts usando o Testify
	assert.Equal(t, http.StatusOK, w.Code, "O status code deveria ser 200 OK")
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"service":"go-engine","status":"pong"}`, w.Body.String())
}

func TestAPI_BossEndpoint(t *testing.T) {
	mux := SetupRouter(nil, time.Now())

	req := httptest.NewRequest("GET", "/boss", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*", "Deveria conter cabeçalho CORS")
}

func TestAPI_AdminRoomsEndpoint(t *testing.T) {
	mux := SetupRouter(nil, time.Now())

	req := httptest.NewRequest("GET", "/admin/rooms", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"pending","message":"Endpoint de criação de sala em construção!"}`, w.Body.String())
}

func TestAPI_MetricsEndpoint(t *testing.T) {
	mux := SetupRouter(nil, time.Now())

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain; version=0.0.4", w.Header().Get("Content-Type"))
}

func TestAPI_RankingEndpoint_WithNilDB(t *testing.T) {
	// Garante que o endpoint lida bem e retorna uma lista vazia sem quebrar se o DB for nulo
	mux := SetupRouter(nil, time.Now())

	req := httptest.NewRequest("GET", "/ranking", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Como passamos nil no DB, a lista de rankings deve retornar serializada como nula/vazia em JSON
	var bodyResponse []interface{}
	err := json.Unmarshal(w.Body.Bytes(), &bodyResponse)
	assert.NoError(t, err)
}

func TestAPI_OptionsPreflight(t *testing.T) {
	// Testa se o tratamento de requisições OPTIONS (CORS Preflight dos navegadores) responde imediatamente
	mux := SetupRouter(nil, time.Now())

	req := httptest.NewRequest("OPTIONS", "/boss", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String(), "Requisições OPTIONS devem retornar corpo vazio no preflight")
}
