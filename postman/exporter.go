package postman

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/bukalapak/snowboard/api"
	"net/http"
	"fmt"
)

const (
	modelTag  = "__model__"
	customTag = "__custom__"
)

var markdownListItem = regexp.MustCompile(`\n\n([-+*].+)`)

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
			Description: fixMarkdown(bp.Description),
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
		Description: fixMarkdown(rg.Description),
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
		Description: fixMarkdown(rsc.Description) + describeModel(rsc),
		Items:       items,
	}, nil
}

func describeModel(rsc *api.Resource) string {
	for _, tr := range rsc.Transitions {
		for _, tx := range tr.Transactions {
			var schema string
			if strings.Contains(tx.Request.Description, modelTag) {
				schema = tx.Request.Schema.Body
			} else  if strings.Contains(tx.Response.Description, modelTag) {
				schema = tx.Response.Schema.Body
			} else {
				continue
			}
			return "\n\n## Model\n" + DescribeJsonSchema([]byte(schema))
		}
	}
	return ""
}

func itemFromTransition(tr *api.Transition) (*Item, error) {
	url, err := formalizeUrl(tr.URL)
	if err != nil {
		return nil, err
	}

	explainQueryParams(url, tr)
	explainVariables(url, tr)

	item := &Item{
		Name: tr.Title,
		Request: &Request{
			Url:         url,
			Method:      tr.Method,
			Description: fixMarkdown(tr.Description),
		},
	}

	if len(tr.Transactions) == 0 {
		return item, nil
	}

	first := tr.Transactions[0]
	item.Request.Header = convertHeaders(first.Request.Headers)

	if first.Request.Body.Body != "" {
		item.Request.Body = Body{
			Mode: "raw",
			Raw:  first.Request.Body.Body,
		}
	}

	if strings.Contains(first.Request.Description, customTag) {
		item.Request.Description += "\n\n###### Request Attributes\n" +
			DescribeJsonSchema([]byte(first.Request.Schema.Body))
	}

	for _, tx := range tr.Transactions {
		body := tx.Response.Body.Body
		if body == "" {
			body = " " // add dummy body, otherwise Postman Docs will not show the response at all
		}
		code := tx.Response.StatusCode
		status := http.StatusText(code)
		item.Response = append(item.Response, Response{
			Name:   fmt.Sprintf("%d %s", code, status),
			Header: convertHeaders(tx.Response.Headers),
			Body:   body,
			Status: status,
			Code:   code,
		})
		if strings.Contains(tx.Response.Description, customTag) {
			item.Request.Description += "\n\n###### Response Attributes\n" +
				DescribeJsonSchema([]byte(first.Response.Schema.Body))
		}
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
	re := regexp.MustCompile(`\{\?([a-zA-Z_,.-]+)\}$`)
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
	re = regexp.MustCompile(`\{([a-zA-Z_.-]+)\}`)
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

func explainQueryParams(u Url, tr *api.Transition) {
	for _, param := range tr.Href.Parameters {
		for i, query := range u.Query {
			if query.Key == param.Key {
				u.Query[i].Value = param.Value
				u.Query[i].Description = param.Description
			}
		}
	}
}

func explainVariables(u Url, tr *api.Transition) {
	for _, param := range tr.Href.Parameters {
		for i, variable := range u.Variable {
			if variable.Key == param.Key {
				u.Variable[i].Value = param.Value
				u.Variable[i].Description = param.Description
			}
		}
	}
}

func fixMarkdown(md string) string {
	// remove double linebreaks before list items
	return string(markdownListItem.ReplaceAll([]byte(md), []byte("\n$1")))
}