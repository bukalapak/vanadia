package postman

import (
	"github.com/bukalapak/snowboard/api"
)

func CreateCollection(bp *api.API) Collection {
	folders := []*Item{}
	for _, resourceGroup := range bp.ResourceGroups {
		folders = append(folders, itemFromResourceGroup(&resourceGroup))
	}

	coll := Collection{
		Info: Information{
			Name:      bp.Title,
			PostmanID: bp.Title,
			Schema:    "https://schema.getpostman.com/json/collection/v2.1.0/docs/index.html",
		},
		Items: folders,
	}

	return coll
}

func itemFromResourceGroup(rg *api.ResourceGroup) *Item {
	items := []*Item{}
	for _, resource := range rg.Resources {
		items = append(items, itemFromResource(resource))
	}

	return &Item{
		Name:  rg.Title,
		Items: items,
	}
}

func itemFromResource(rsc *api.Resource) *Item {
	items := []*Item{}
	for _, transition := range rsc.Transitions {
		items = append(items, itemFromTransition(transition))
	}

	return &Item{
		Name:  rsc.Title,
		Items: items,
	}
}

func itemFromTransition(tr *api.Transition) *Item {
	item := &Item{
		Name: tr.Title,
		Request: Request{
			Url:    tr.URL,
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

	return item
}
