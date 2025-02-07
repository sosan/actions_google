package transform

import (
	"reflect"
	"testing"
)

func TestProcessNotionData(t *testing.T) {
	mockResults := []interface{}{
		map[string]interface{}{
			"properties": map[string]interface{}{
				"Name": map[string]interface{}{
					"type": "title",
					"title": []interface{}{
						map[string]interface{}{
							"plain_text": "Sample Title",
						},
					},
				},
				"Tags": map[string]interface{}{
					"type": "multi_select",
					"multi_select": []interface{}{
						map[string]interface{}{
							"name": "Tag1",
						},
						map[string]interface{}{
							"name": "Tag2",
						},
					},
				},
				"Date": map[string]interface{}{
					"type": "date",
					"date": map[string]interface{}{
						"start": "2023-01-01",
						"end":   "2023-01-02",
					},
				},
			},
		},
		map[string]interface{}{
			"properties": map[string]interface{}{
				"Name": map[string]interface{}{
					"type": "title",
					"title": []interface{}{
						map[string]interface{}{
							"plain_text": "Another Title",
						},
					},
				},
				"Tags": map[string]interface{}{
					"type": "multi_select",
					"multi_select": []interface{}{
						map[string]interface{}{
							"name": "Tag3",
						},
					},
				},
				"Date": map[string]interface{}{
					"type": "date",
					"date": map[string]interface{}{
						"start": "2023-02-01",
					},
				},
			},
		},
	}

	repo := NewActionsClient()
	headers, csvData := repo.ProcessNotionData(&mockResults)

	expectedHeaders := []string{"Name", "Tags", "Date"}
	expectedCSVData := [][]string{
		{"Sample Title", "Tag1, Tag2", "2023-01-01 → 2023-01-02"},
		{"Another Title", "Tag3", "2023-02-01"},
	}

	if !reflect.DeepEqual(*headers, expectedHeaders) {
		t.Errorf("Expected headers %v, got %v", expectedHeaders, *headers)
	}

	if !reflect.DeepEqual(*csvData, expectedCSVData) {
		t.Errorf("Expected CSV data %v, got %v", expectedCSVData, *csvData)
	}
}

func TestExtractPropertyValue(t *testing.T) {
	repo := NewActionsClient()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name: "Title Property",
			input: map[string]interface{}{
				"type": "title",
				"title": []interface{}{
					map[string]interface{}{
						"plain_text": "Sample Title",
					},
				},
			},
			expected: "Sample Title",
		},
		{
			name: "Multi-Select Property",
			input: map[string]interface{}{
				"type": "multi_select",
				"multi_select": []interface{}{
					map[string]interface{}{
						"name": "Tag1",
					},
					map[string]interface{}{
						"name": "Tag2",
					},
				},
			},
			expected: "Tag1, Tag2",
		},
		{
			name: "Date Property",
			input: map[string]interface{}{
				"type": "date",
				"date": map[string]interface{}{
					"start": "2023-01-01",
					"end":   "2023-01-02",
				},
			},
			expected: "2023-01-01 → 2023-01-02",
		},
		{
			name: "Checkbox Property",
			input: map[string]interface{}{
				"type":     "checkbox",
				"checkbox": true,
			},
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.extractPropertyValue(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
