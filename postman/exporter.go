package postman

import (
	"github.com/bukalapak/snowboard/api"
)

func CreateCollection(bp *api.API) Collection {
	folders := make([]Folder, 0)
	for _, resourceGroup := range bp.ResourceGroups {
		folders = append(folders, folderFromResourceGroup(&resourceGroup))
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

func folderFromResourceGroup(rg *api.ResourceGroup) Folder {
	items := make([]interface{}, 0)
	for _, resource := range rg.Resources {
		items = append(items, folderFromResource(resource))
	}

	return Folder{
		Name:  rg.Title,
		Items: items,
	}
}

func folderFromResource(rsc *api.Resource) Folder {
	items := make([]interface{}, 0)
	for _, transition := range rsc.Transitions {
		items = append(items, itemFromTransition(transition))
	}

	return Folder{
		Name:  rsc.Title,
		Items: items,
	}
}

func itemFromTransition(tr *api.Transition) Item {
	item := Item{
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
