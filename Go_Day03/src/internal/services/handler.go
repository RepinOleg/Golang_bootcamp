package services

import (
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"loaderData/internal/db"
	"loaderData/utils/errors"
	"log"
	"net/http"
	"strings"
)

// секретный ключ
var sampleSecretKey = []byte("MyKey")

// VerifyJWT проверка валдиности введенного токена ex04
func VerifyJWT(w http.ResponseWriter, r *http.Request) {
	// Забираем строку из хедера с ключом Auth
	tokenStr := r.Header.Get("Authorization")

	// Если нет токена возвращаем ошибку 401
	if tokenStr == "" {
		errors.HandleError(w, "Token is missing", http.StatusUnauthorized)
		return
	}

	//разбиваем токен на две строки
	tokenString := strings.Split(tokenStr, " ")

	// Если их не две или первая не равна Bearer возвращаем ошибку 401
	if len(tokenString) != 2 || tokenString[0] != "Bearer" {
		errors.HandleError(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	//проверяем валиден ли токен
	token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	//если нет возвращаем ошибку 401
	if err != nil || !token.Valid {
		errors.HandleError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Готовим ближайшие рестораны
	es, err := db.CreateClient()
	if err != nil {
		log.Fatal(err)
	}

	places, err := es.PrepareNearestRest(r)
	if err != nil {
		WriteJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Кладем их в структуру
	jsonRequest := db.JsonClosestData{
		Name:   "Recommendation",
		Places: places,
	}

	// Отправляем JSON ответ и выставляем 200 статус
	WriteJson(w, jsonRequest, http.StatusOK)
}

// HandlerJWT генерирует токен ex04
func HandlerJWT(w http.ResponseWriter, r *http.Request) {
	token, err := generateJWT()
	if err != nil {
		errors.HandleError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"token": token,
	}
	WriteJson(w, response, http.StatusOK)
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// HandlerNearestRest from ex03
func HandlerNearestRest(w http.ResponseWriter, r *http.Request) {
	esClient, err := db.CreateClient()
	if err != nil {
		log.Fatal(err)
	}
	places, err := esClient.PrepareNearestRest(r)

	if err != nil {
		WriteJson(w, errors.ErrorResponse{Err: err.Error()}, http.StatusBadRequest)
		return
	}

	jsonRequest := db.JsonClosestData{
		Name:   "Recommendation",
		Places: places,
	}
	WriteJson(w, jsonRequest, http.StatusOK)
}

// HandlerJson from ex02
func HandlerJson(w http.ResponseWriter, r *http.Request) {
	// создание соединения с базой
	esClient, err := db.CreateClient()
	if err != nil {
		log.Fatal(err)
	}
	// Получаем данные и подготавливаем их для отображения
	data, err := esClient.PrepareData(r)
	if err != nil {
		WriteJson(w, errors.ErrorResponse{Err: err.Error()}, http.StatusBadRequest)
		return
	}

	jsonRequest := db.JsonData{
		Name:     "Places",
		Total:    data.Total,
		Places:   data.Places,
		PrevPage: data.PrevPage,
		NextPage: data.NextPage,
		LastPage: data.LastPage,
	}
	WriteJson(w, jsonRequest, http.StatusOK)
}

// HandlerGetPlaces from ex01
func HandlerGetPlaces(writer http.ResponseWriter, request *http.Request) {
	// создание соединения с базой
	esClient, err := db.CreateClient()
	if err != nil {
		log.Fatal(err)
	}
	// Получаем данные и подготавливаем их для отображения
	data, err := esClient.PrepareData(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Загрузка шаблона HTML из файла
	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		http.Error(writer, "Failed to load template", http.StatusInternalServerError)
		return
	}

	// Отображение данных с использованием шаблона
	err = tmpl.Execute(writer, data)
	if err != nil {
		http.Error(writer, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
