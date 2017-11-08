package postman

import (
	"fmt"
	"strings"
)

func SchemeToEnv(c *Collection, placeholder string) {
	items := getItemsFromCollection(c)

	for i := range items {
		if items[i].Request.Url.Protocol != "" {
			items[i].Request.Url.Protocol = fmt.Sprintf("{{%s}}", placeholder)
		}
	}
}

func HostToEnv(c *Collection, n int) {
	items := getItemsFromCollection(c)

	for i := range items {
		host := items[i].Request.Url.Host
		if host != "" {
			hostSegments := strings.Split(host, ".")

			newHostSegments := []string{}
			if len(hostSegments) > n {
				newHostSegments = append(
					newHostSegments,
					hostSegments[0:len(hostSegments)-2]...,
				)
			}
			newHostSegments = append(newHostSegments, "{{host}}")
			items[i].Request.Url.Host = strings.Join(newHostSegments, ".")
		}
	}
}

func AuthTokenToEnv(c *Collection) {
	items := getItemsFromCollection(c)

	for i := range items {
		for j := range items[i].Request.Header {
			header := items[i].Request.Header[j]

			if strings.ToLower(header.Key) != "authorization" {
				continue
			}

			values := strings.Split(header.Value, " ")
			values[len(values)-1] = "{{access_token}}"

			items[i].Request.Header[j].Value = strings.Join(values, " ")
		}
	}
}

func AddGlobalHeaders(c *Collection, headers []Header) {
	items := getItemsFromCollection(c)

	for i := range items {
		items[i].Request.Header = append(items[i].Request.Header, headers...)
	}
}

func getItemsFromCollection(c *Collection) []*Item {
	items := make([]*Item, 0)

	for i := range c.Items {
		items = append(items, c.Items[i])
		items = append(items, getItemsFromItem(c.Items[i])...)
	}

	return items
}

func getItemsFromItem(f *Item) []*Item {
	items := make([]*Item, 0)

	for i := range f.Items {
		items = append(items, f.Items[i])
		items = append(items, getItemsFromItem(f.Items[i])...)
	}

	return items
}
