package postman

import (
	json "github.com/buger/jsonparser"
	"strings"
)

func DescribeJsonSchema(schema []byte) string {
	b := strings.Builder{}
	describeJsonType(&b, schema, true)
	return b.String()
}

func describeJsonType(b *strings.Builder, schema []byte, outerFrame bool) {
	writeType := func() {
		typ, dataType, _, _ := json.Get(schema, "type")
		b.WriteString(" _") // TODO fix when "anyOf"
		switch dataType {
		case json.String:
			b.WriteString(string(typ))
		case json.Array:
			json.ArrayEach(schema, func(value []byte, _ json.ValueType, _ int, _ error) {
				b.WriteString(string(value))
				b.WriteRune(' ')
			})
		case json.NotExist:
			b.WriteString("any")
		}
		b.WriteString("_")

		if enum, dataType, _, _ := json.Get(schema, "enum"); dataType == json.Array {
			// generated enum values might contain duplicates in some situations
			values := map[string]bool{}
			b.WriteString(", one of:")
			json.ArrayEach(enum, func(value []byte, _ json.ValueType, _ int, _ error) {
				v := string(value)
				if !values[v] {
					values[v] = true
					b.WriteString(" `")
					b.WriteString(v)
					b.WriteString("`")
				}
			})
		}

		b.WriteRune('\n')
	}
	writeDesc := func() {
		if desc, dataType, _, _ := json.Get(schema, "description"); dataType == json.String {
			b.Write(desc)
		}
	}

	if props, dataType, _, _ := json.Get(schema, "properties"); dataType == json.Object { // object
		if !outerFrame {
			writeType()
		}
		writeDesc()
		describeObject(b, props, buildRequired(schema))
	} else {
		if outerFrame {
			b.WriteString("<table><tr><td>")
		}
		writeType()
		writeDesc()
		items, dataType, _, _ := json.Get(schema, "items")
		if dataType == json.Object {
			describeJsonType(b, items, true)
		} else if dataType == json.Array {
			describeArray(b, items)
		}
		if outerFrame {
			b.WriteString("</td></tr></table>")
		}
	}
}

func describeObject(b *strings.Builder, props []byte, required map[string]bool) {
	b.WriteString("<table>")
	json.ObjectEach(props, func(key []byte, value []byte, _ json.ValueType, _ int) error {
		b.WriteString("<tr><td>`")
		b.WriteString(string(key))
		if required[string(key)] {
			b.WriteString("` \\*</td><td>")
		} else {
			b.WriteString("`</td><td>")
		}
		describeJsonType(b, value, false)
		b.WriteString("</td></tr>")
		return nil
	})
	b.WriteString("</table>")
}

func describeArray(b *strings.Builder, items []byte) {
	json.ArrayEach(items, func(value []byte, _ json.ValueType, _ int, _ error) {
		describeJsonType(b, value, true)
	})
}

func buildRequired(schema []byte) map[string]bool {
	var reqs = map[string]bool{}
	json.ArrayEach(schema, func(value []byte, _ json.ValueType, _ int, _ error) {
		reqs[string(value)] = true
	}, "required")
	return reqs
}
