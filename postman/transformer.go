package postman

func SchemeToVariable(c *Collection) {
	items := getItemsFromCollection(c)

	for i := range items {
		if items[i].Request.Url.Protocol != "" {
			items[i].Request.Url.Protocol = "{{scheme}}"
		}
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
