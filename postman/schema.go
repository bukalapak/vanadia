package postman

// https://schema.getpostman.com/json/collection/v2.1.0/docs/index.html

type Collection struct {
	Info  Information `json:"info"`
	Items []*Item     `json:"item,omitempty"`
}

type Information struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

// Item represents both `Folder` and `Item` in Postman schema.
type Item struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Request     *Request   `json:"request,omitempty"`
	Response    []Response `json:"response,omitempty"`
	Items       []*Item    `json:"item,omitempty"`
}

type Request struct {
	Url         Url      `json:"url,omitempty"`
	Method      string   `json:"method,omitempty"`
	Description string   `json:"description,omitempty"`
	Header      []Header `json:"header,omitempty"`
	Body        Body     `json:"body,omitempty"`
}

type Response struct {
	Name   string   `json:"name"` // undocumented (not in the schema)
	Header []Header `json:"header,omitempty"`
	Body   string   `json:"body,omitempty"`
	Status string   `json:"status,omitempty"`
	Code   int      `json:"code,omitempty"`
}

type Url struct {
	Protocol string       `json:"protocol,omitempty"`
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

type Bearer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
