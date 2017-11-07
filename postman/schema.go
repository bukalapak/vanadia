package postman

// https://schema.getpostman.com/json/collection/v2.1.0/docs/index.html

type Collection struct {
	Info  Information `json:"info"`
	Items []*Item     `json:"item,omitempty"`
}

type Information struct {
	Name      string `json:"name"`
	PostmanID string `json:"_postman_id,omitempty"`
	Schema    string `json:"schema"`
}

// Item represents both `Folder` and `Item` in Postman schema.
type Item struct {
	Name    string  `json:"name,omitempty"`
	Items   []*Item `json:"item,omitempty"`
	Request Request `json:"request"`
}

type Request struct {
	Url    string   `json:"url,omitempty"`
	Method string   `json:"method,omitempty"`
	Header []Header `json:"header,omitempty"`
	Body   Body     `json:"body,omitempty"`
}

type Header struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type Body struct {
	Mode string `json:"mode,omitempty"`
	Raw  string `json:"raw,omitempty"`
}
