package postman

import (
	"encoding/json"

	"github.com/bukalapak/snowboard/api"
)

func CreateCollection(bp *api.API) ([]byte, error) {
	coll := Collection{
		Info: Information{
			Name:      bp.Title,
			PostmanID: bp.Title,
			Schema:    "http://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
	}

	return json.Marshal(coll)
}
