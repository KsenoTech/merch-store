package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/KsenoTech/merch-store/internal/database"
	"github.com/KsenoTech/merch-store/internal/handlers"
	"github.com/KsenoTech/merch-store/internal/middleware"
	"github.com/KsenoTech/merch-store/internal/migrations"
	"github.com/KsenoTech/merch-store/internal/repository"
	"github.com/KsenoTech/merch-store/internal/routes"
	service "github.com/KsenoTech/merch-store/internal/services"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() (http.Handler, *httptest.Server) {

	os.Setenv("DATABASE_URL", "postgres://postgres:590789@localhost:5432/merch_test?sslmode=disable")

	// Подключаемся к базе данных
	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Применяем миграции
	if err := migrations.ApplyMigrations(db); err != nil {
		panic(err)
	}

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)

	transactionDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	transactionRepo := repository.NewTransactionRepository(transactionDB)

	// Создаем сервисы
	authService := service.NewAuthService("SECRET_KEY", userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(authService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Создаем JWT-middleware
	jwtMiddleware := middleware.JWTMiddleware("SECRET_KEY")

	// Создаем роутер
	router := routes.SetupRoutes(authHandler, transactionHandler, jwtMiddleware)

	// Запускаем тестовый сервер
	server := httptest.NewServer(router)
	return router, server
}

func sendRequest(method, url string, body interface{}) (*http.Response, error) {
	payload, _ := json.Marshal(body)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func sendAuthorizedRequest(method, url, token string, body interface{}) (*http.Response, error) {
	payload, _ := json.Marshal(body)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	return client.Do(req)
}

func TestAuth(t *testing.T) {
	// Запускаем тестовый сервер
	_, server := setupTestServer()
	defer server.Close()

	// Регистрируем нового пользователя
	registerURL := fmt.Sprintf("%s/api/auth", server.URL)
	registerBody := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	resp, err := sendRequest("POST", registerURL, registerBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Читаем токен из ответа
	body, _ := ioutil.ReadAll(resp.Body)
	var authResponse map[string]string
	json.Unmarshal(body, &authResponse)
	token := authResponse["token"]
	assert.NotEmpty(t, token)
}

func TestSendCoins(t *testing.T) {
	// Запускаем тестовый сервер
	_, server := setupTestServer()
	defer server.Close()

	// Регистрируем двух пользователей
	registerURL := fmt.Sprintf("%s/api/auth", server.URL)
	for _, username := range []string{"user1", "user2"} {
		resp, err := sendRequest("POST", registerURL, map[string]string{
			"username": username,
			"password": "password123",
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	// Авторизуемся как user1
	loginResp, err := sendRequest("POST", registerURL, map[string]string{
		"username": "user1",
		"password": "password123",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	body, _ := ioutil.ReadAll(loginResp.Body)
	var authResponse map[string]string
	json.Unmarshal(body, &authResponse)
	token := authResponse["token"]
	assert.NotEmpty(t, token)

	// Отправляем монеты от user1 к user2
	sendCoinURL := fmt.Sprintf("%s/api/sendCoin", server.URL)
	sendCoinBody := map[string]interface{}{
		"toUser": "user2",
		"amount": 100,
	}
	resp, err := sendAuthorizedRequest("POST", sendCoinURL, token, sendCoinBody)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверяем баланс user1 после отправки
	infoURL := fmt.Sprintf("%s/api/info", server.URL)
	infoResp, err := sendAuthorizedRequest("GET", infoURL, token, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoResp.StatusCode)

	infoBody, _ := ioutil.ReadAll(infoResp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(infoBody, &userInfo)
	assert.Equal(t, float64(900), userInfo["coins"]) // Начальный баланс 1000 - 100 = 900
}

func TestBuyMerch(t *testing.T) {
	// Запускаем тестовый сервер
	_, server := setupTestServer()
	defer server.Close()

	// Регистрируем пользователя
	registerURL := fmt.Sprintf("%s/api/auth", server.URL)
	resp, err := sendRequest("POST", registerURL, map[string]string{
		"username": "testuser",
		"password": "password123",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Авторизуемся
	body, _ := ioutil.ReadAll(resp.Body)
	var authResponse map[string]string
	json.Unmarshal(body, &authResponse)
	token := authResponse["token"]
	assert.NotEmpty(t, token)

	// Покупаем t-shirt
	buyURL := fmt.Sprintf("%s/api/buy/t-shirt", server.URL)
	buyResp, err := sendAuthorizedRequest("GET", buyURL, token, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, buyResp.StatusCode)

	// Проверяем инвентарь после покупки
	infoURL := fmt.Sprintf("%s/api/info", server.URL)
	infoResp, err := sendAuthorizedRequest("GET", infoURL, token, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, infoResp.StatusCode)

	infoBody, _ := ioutil.ReadAll(infoResp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(infoBody, &userInfo)
	inventory := userInfo["inventory"].([]interface{})
	assert.Equal(t, 1, len(inventory))

	item := inventory[0].(map[string]interface{})
	assert.Equal(t, "t-shirt", item["type"])
	assert.Equal(t, float64(1), item["quantity"])
}
