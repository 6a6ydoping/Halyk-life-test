package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"halykTestTask/models"
	"log"
	"net/http"
	"net/url"
	"os"
)

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func ProxyRequest(request models.Request) (*http.Response, error) {
	u, err := url.Parse(request.URL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(request.Method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SaveRequestsResponses(r []models.RequestResponse) {
	file, err := os.Create("requests_and_responses.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fmt.Fprintln(file)
	if err != nil {
		log.Fatal(err)
	}
}
