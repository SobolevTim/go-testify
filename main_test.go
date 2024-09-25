package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerWhenStatusOk(t *testing.T) {
	//ожидаемый статус
	expected := http.StatusOK

	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	//проверка на статус кода
	require.Equal(t, expected, responseRecorder.Code)
	//проверка что список не пустой
	assert.NotEmpty(t, responseRecorder.Body.String())
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenCityNotSupported(t *testing.T) {
	//ожидаемый статус
	expected := http.StatusBadRequest
	//ожидаемый ответ сервера
	responseAnswer := `wrong city value`

	//запрос с другим городом
	req := httptest.NewRequest("GET", "/cafe?count=4&city=kazan", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	//проверка на статус кода
	require.Equal(t, expected, responseRecorder.Code)
	//проверка текста ответа сервера
	assert.Equal(t, responseAnswer, responseRecorder.Body.String())
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	//ожидаемое кол-во
	totalCount := 4
	//ожидаемый статус
	expected := http.StatusOK
	//ожидаемый список
	expectedList := []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"}

	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	responseBody := responseRecorder.Body.String()
	responseList := strings.Split(responseBody, ",")

	//проверка на статус кода
	require.Equal(t, expected, responseRecorder.Code)
	//проверка соответствия списка
	assert.Equal(t, expectedList, responseList)
	//проверка количества
	assert.Equal(t, totalCount, len(responseList))
}
