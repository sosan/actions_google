package notion

import (
	"fmt"
	"strings"
)

type ActionsNotion struct {
}

func NewActionsClient() *ActionsNotion {
	return &ActionsNotion{}
}

func (n *ActionsNotion) ProcessNotionData(results *[]interface{}) (*[]string, *[][]string) {
	headerMap := make(map[string]bool)
	var headers []string
	// unique headers
	for _, page := range *results {
		props, ok := page.(map[string]interface{})["properties"].(map[string]interface{})
		if !ok {
			continue
		}
		for propName := range props {
			if !headerMap[propName] {
				headerMap[propName] = true
				headers = append(headers, propName)
			}
		}
	}

	var csvData [][]string
	for _, page := range *results {
		props, ok := page.(map[string]interface{})["properties"].(map[string]interface{})
		if !ok {
			continue
		}
		record := make([]string, len(headers))
		for i, header := range headers {
			propData, exists := props[header]
			if !exists {
				record[i] = ""
				continue
			}
			record[i] = n.extractPropertyValue(propData.(map[string]interface{}))
		}
		csvData = append(csvData, record)
	}
	return &headers, &csvData
}

func (n *ActionsNotion) extractPropertyValue(prop map[string]interface{}) string {
	propType, ok := prop["type"].(string)
	if !ok {
		return ""
	}

	switch propType {
	case "title", "rich_text":
		if items, ok := prop[propType].([]interface{}); ok && len(items) > 0 {
			if firstItem, ok := items[0].(map[string]interface{}); ok {
				if plainText, ok := firstItem["plain_text"].(string); ok {
					return plainText
				}
			}
		}

	case "email", "phone_number", "url":
		if value, ok := prop[propType].(string); ok {
			return value
		}

	case "unique_id":
		if uidData, ok := prop[propType].(map[string]interface{}); ok {
			prefix, _ := uidData["prefix"].(string)
			number, _ := uidData["number"].(float64)
			return fmt.Sprintf("%s%v", prefix, number)
		}

	case "multi_select":
		if selectData, ok := prop[propType].([]interface{}); ok {
			var names []string
			for _, item := range selectData {
				if opt, ok := item.(map[string]interface{}); ok {
					if name, ok := opt["name"].(string); ok {
						names = append(names, name)
					}
				}
			}
			return strings.Join(names, ", ")
		}

	case "files":
		if filesData, ok := prop[propType].([]interface{}); ok {
			var files []string
			for _, f := range filesData {
				if file, ok := f.(map[string]interface{}); ok {
					if name, ok := file["name"].(string); ok {
						files = append(files, name)
					} else if url, ok := file["url"].(string); ok {
						files = append(files, url)
					}
				}
			}
			return strings.Join(files, ", ")
		}

	case "select":
		if selectData, ok := prop[propType].(map[string]interface{}); ok {
			if name, ok := selectData["name"].(string); ok {
				return name
			}
		}

	case "date":
		if dateData, ok := prop[propType].(map[string]interface{}); ok {
			start, _ := dateData["start"].(string)
			end, _ := dateData["end"].(string)
			if end != "" {
				return fmt.Sprintf("%s → %s", start, end)
			}
			return start
		}

	case "number":
		if number, ok := prop[propType].(float64); ok {
			return fmt.Sprintf("%v", number)
		}

	case "checkbox":
		if checked, ok := prop[propType].(bool); ok {
			return fmt.Sprintf("%t", checked)
		}

	default:
		// not impleemented
		return fmt.Sprintf("[not implemented: %s]", propType)
	}

	return ""
}
