package postman

import (
	"encoding/json"
	"strings"
	"html"
	"fmt"
)

func DescribeJsonSchema(schema string) string {
	var s map[string]interface{}
	json.Unmarshal([]byte(schema), &s)
	b := strings.Builder{}
	describeJsonType(&b, s, true)
	return b.String()
}

func describeJsonType(b *strings.Builder, jsonType map[string]interface{}, outerFrame bool) {
	writeType := func() {
		if typ, ok := jsonType["type"]; ok {
			b.WriteString(" _")
			b.WriteString(html.EscapeString(fmt.Sprintf("%v", typ)))
			b.WriteString("_\n")
		}
	}
	writeDesc := func() {
		if desc, ok := jsonType["description"]; ok {
			b.WriteString(desc.(string))
		}
	}

	if props, ok := jsonType["properties"]; ok { // object
		if !outerFrame {
			writeType()
		}
		writeDesc()
		describeObject(b, props.(map[string]interface{}), requiredFunc(jsonType))
	} else {
		if outerFrame {
			b.WriteString("<table><tr><td>")
		}
		writeType()
		writeDesc()
		if items, ok := jsonType["items"]; ok { // array
			describeArray(b, items)
		}
		if outerFrame {
			b.WriteString("</td></tr></table>")
		}
	}
}

func describeObject(b *strings.Builder, props map[string]interface{}, required func(string) bool) {
	b.WriteString("<table>")
	for name, v := range props {
		b.WriteString("<tr><td>`")
		b.WriteString(html.EscapeString(name))
		if required(name) {
			b.WriteString("` \\*</td><td>")
		} else {
			b.WriteString("`</td><td>")
		}
		describeJsonType(b, v.(map[string]interface{}), false)
		b.WriteString("</td></tr>")
	}
	b.WriteString("</table>")
}

func describeArray(b *strings.Builder, items interface{}) {
	b.WriteRune('\n')
	switch it := items.(type) {
	case map[string]interface{}:
		describeJsonType(b, it, true)
	case []interface{}:
		for _, item := range it {
			describeJsonType(b, item.(map[string]interface{}), true)
		}
	}
}

func requiredFunc(jsonType map[string]interface{}) func(string) bool {
	if req, ok := jsonType["required"]; ok {
		return func(name string) bool {
			for _, attr := range req.([]interface{}) {
				if attr.(string) == name {
					return true
				}
			}
			return false
		}
	}
	return func(string) bool {
		return false
	}
}
