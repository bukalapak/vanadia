package postman

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/bukalapak/snowboard/api"
	"net/http"
)

func CreateCollection(bp *api.API) (Collection, error) {
	folders := []*Item{}
	for _, resourceGroup := range bp.ResourceGroups {
		item, err := itemFromResourceGroup(&resourceGroup)
		if err != nil {
			return Collection{}, err
		}

		folders = append(folders, item)
	}

	coll := Collection{
		Info: Information{
			Name:        bp.Title,
			Description: bp.Description,
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Items: folders,
	}

	return coll, nil
}

func itemFromResourceGroup(rg *api.ResourceGroup) (*Item, error) {
	items := []*Item{}
	for _, resource := range rg.Resources {
		item, err := itemFromResource(resource)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &Item{
		Name:        rg.Title,
		Description: rg.Description,
		Items:       items,
	}, nil
}

func itemFromResource(rsc *api.Resource) (*Item, error) {
	items := []*Item{}
	for _, transition := range rsc.Transitions {
		item, err := itemFromTransition(transition)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &Item{
		Name:        rsc.Title,
		Description: rsc.Description,
		Items:       items,
	}, nil
}

func itemFromTransition(tr *api.Transition) (*Item, error) {
	url, err := formalizeUrl(tr.URL)
	if err != nil {
		return nil, err
	}

	url = explainQueryParams(url, tr)
	url = explainVariables(url, tr)

	item := &Item{
		Name: tr.Title,
		Request: &Request{
			Url:         url,
			Method:      tr.Method,
			Description: tr.Description,
		},
	}

	first := tr.Transactions[0]
	item.Request.Header = convertHeaders(first.Request.Headers)

	if first.Request.Body.Body != "" {
		item.Request.Body = Body{
			Mode: "raw",
			Raw:  first.Request.Body.Body,
		}
	}

	if first.Request.Schema.Body != "" {
		item.Request.Description += "\n\n###### Request Attributes\n" + DescribeJsonSchema(first.Request.Schema.Body)
	}

	if first.Response.Schema.Body != "" {
		item.Request.Description += "\n\n###### Response Attributes\n" + DescribeJsonSchema(first.Response.Schema.Body)
	}

	for _, tx := range tr.Transactions {
		body := tx.Response.Body.Body
		if body == "" {
			body = " " // add dummy body, otherwise Postman Docs will not show the response at all
		}
		item.Response = append(item.Response, Response{
			Name:   " ", // suppress "Untitled Response" in Postman Docs
			Header: convertHeaders(tx.Response.Headers),
			Body:   body,
			Status: http.StatusText(tx.Response.StatusCode),
			Code:   tx.Response.StatusCode,
		})
	}

	return item, nil
}

func convertHeaders(apiHeaders []api.Header) []Header {
	var headers []Header
	for _, header := range apiHeaders {
		headers = append(headers, Header{
			Key:   header.Key,
			Value: header.Value,
		})
	}
	return headers
}

func formalizeUrl(urlString string) (Url, error) {
	var queries []QueryParam
	var variables []Variable

	// Get query parameters
	re := regexp.MustCompile(`\{\?([a-z_,-]+)\}$`)
	allMatch := re.FindAllStringSubmatch(urlString, -1)

	for _, innerMatch := range allMatch {
		for _, match := range strings.Split(innerMatch[1], `,`) {
			query := QueryParam{
				Key: match,
			}
			queries = append(queries, query)
		}
	}

	urlString = re.ReplaceAllString(urlString, "")

	// Get variables
	re = regexp.MustCompile(`\{([a-z_-]+)\}`)
	allMatch = re.FindAllStringSubmatch(urlString, -1)

	for _, innerMatch := range allMatch {
		for _, match := range strings.Split(innerMatch[1], `,`) {
			variable := Variable{
				Key: match,
			}
			variables = append(variables, variable)
		}
	}

	urlString = re.ReplaceAllString(urlString, ":$1")

	// Get scheme, host, and path
	urlObject, err := url.ParseRequestURI(urlString)
	if err != nil {
		return Url{}, err
	}

	return Url{
		Protocol: urlObject.Scheme,
		Host:     urlObject.Host,
		Path:     urlObject.Path,
		Query:    queries,
		Variable: variables,
	}, nil
}

func explainQueryParams(u Url, tr *api.Transition) Url {
	queryMap := make(map[string]QueryParam)

	for _, query := range u.Query {
		queryMap[query.Key] = query
	}

	for _, param := range tr.Href.Parameters {
		if query, found := queryMap[param.Key]; found {
			query.Value = param.Value
			query.Description = param.Description

			queryMap[param.Key] = query
		}
	}

	querySlice := []QueryParam{}
	for _, query := range queryMap {
		querySlice = append(querySlice, query)
	}

	u.Query = querySlice
	return u
}

func explainVariables(u Url, tr *api.Transition) Url {
	variableMap := make(map[string]Variable)

	for _, variable := range u.Variable {
		variableMap[variable.Key] = variable
	}

	for _, param := range tr.Href.Parameters {
		if variable, found := variableMap[param.Key]; found {
			variable.Value = param.Value
			variable.Description = param.Description

			variableMap[param.Key] = variable
		}
	}

	variableSclie := []Variable{}
	for _, variable := range variableMap {
		variableSclie = append(variableSclie, variable)
	}

	u.Variable = variableSclie
	return u
}
