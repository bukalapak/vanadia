package postman

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/bukalapak/snowboard/api"
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
			Name:      bp.Title,
			PostmanID: bp.Title,
			Schema:    "https://schema.getpostman.com/json/collection/v2.1.0/docs/index.html",
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
		Name:  rg.Title,
		Items: items,
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
		Name:  rsc.Title,
		Items: items,
	}, nil
}

func itemFromTransition(tr *api.Transition) (*Item, error) {
	url, err := formalizeUrl(tr.URL)
	if err != nil {
		return nil, err
	}

	url = explainQueryParams(url, tr)

	item := &Item{
		Name: tr.Title,
		Request: Request{
			Url:    url,
			Method: tr.Method,
		},
	}

	headers := []Header{}
	for _, header := range tr.Transactions[0].Request.Headers {
		headers = append(headers, Header{
			Key:   header.Key,
			Value: header.Value,
		})
	}
	if len(headers) > 0 {
		item.Request.Header = headers
	}

	if tr.Transactions[0].Request.Body.Body != "" {
		item.Request.Body = Body{
			Mode: "raw",
			Raw:  tr.Transactions[0].Request.Body.Body,
		}
	}

	return item, nil
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
