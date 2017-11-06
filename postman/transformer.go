package postman

import (
	"encoding/json"

	"github.com/bukalapak/snowboard/api"
)

func CreateCollection(bp *api.API) ([]byte, error) {
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

	return json.Marshal(coll)
}

func folderFromResourceGroup(rg *api.ResourceGroup) Folder {
	items := make([]Folder, 0)
	for _, resource := range rg.Resources {
		items = append(items, folderFromResource(resource))
	}

	return Folder{
		Name:  rg.Title,
		Items: items,
	}
}

func folderFromResource(rsc *api.Resource) Folder {
	items := make([]Item, 0)
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

	if tr.Transactions[0].Request.Body.Body != "" {
		item.Request.Body = Body{
			Mode: "raw",
			Raw:  tr.Transactions[0].Request.Body.Body,
		}
	}

	return item
}
