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
	Url    Url      `json:"url,omitempty"`
	Method string   `json:"method,omitempty"`
	Header []Header `json:"header,omitempty"`
	Body   Body     `json:"body,omitempty"`
}

type Url struct {
	Protocol string       `json:"protocol"`
	Host     string       `json:"host"`
	Path     string       `json:"path"`
	Query    []QueryParam `json:"query,omitempty"`
	Variable []Variable   `json:"variable,omitempty"`
}

type Header struct {
	Key         string `json:"key" yaml:"Key"`
	Value       string `json:"value" yaml:"Value,omitempty"`
	Description string `json:"description,omitempty" yaml:"Description,omitempty"`
}

type Body struct {
	Mode string `json:"mode,omitempty"`
	Raw  string `json:"raw,omitempty"`
}

type QueryParam struct {
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type Variable struct {
	Key         string `json:"key"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}
