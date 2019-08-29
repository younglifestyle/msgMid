package httpClient

import "net/http"

type Client struct {
	httpClient http.Client
	retries    int
	timeout    int
}

func CreateHttpClient()  {
	
}