package controllers

import (
	"encoding/json"
	"fmt"
	"halykTestTask/middlewares"
	"halykTestTask/models"
	"io"
	"net/http"
	"strings"
)

var requestsResponses []models.RequestResponse

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestID := r.Header.Get("X-Request-Id")
	if requestID == "" {
		r.Header.Set("X-Request-Id", middlewares.GenerateUUID())
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request models.Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(request)

	resp, err := middlewares.ProxyRequest(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	headerMap := make(map[string]string)
	for key, values := range resp.Header {
		headerMap[key] = strings.Join(values, ",")
	}
	fmt.Println(headerMap)
	response := models.Response{
		ID:      requestID,
		Status:  resp.StatusCode,
		Headers: headerMap,
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Body = string(bodyBytes)
	response.Length = int64(len(bodyBytes))

	requestResponse := models.RequestResponse{
		Request:  request,
		Response: response,
	}
	requestsResponses = append(requestsResponses, requestResponse)

	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
	defer middlewares.SaveRequestsResponses(requestsResponses)
}
