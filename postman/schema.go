package postman

// https://schema.getpostman.com/json/collection/v2.1.0/docs/index.html

type Collection struct {
	Info Information `json:"info"`
}

type Information struct {
	Name      string `json:"name"`
	PostmanID string `json:"_postman_id",omitempty`
	Schema    string `json:"schema"`
}
