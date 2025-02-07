package services

import (
	"actions_google/pkg/domain/models"
	"actions_google/pkg/domain/repos"
	"actions_google/tests"
	"encoding/json"
	"reflect"
	"testing"
)

func TestActionsServiceImpl_GetDatabaseID(t *testing.T) {
	type fields struct {
		RedisRepo             repos.ActionsRedisRepoInterface
		BrokerActionsRepo     repos.ActionsBrokerRepository
		BrokerCredentialsRepo repos.CredentialBrokerRepository
		HTTPRepo              repos.ActionsHTTPRepository
		CredentialHTTP        repos.CredentialHTTPRepository
		ActionsNotion         repos.TransformNotion
		TokenAuth             repos.TokenAuth
		SheetUtils            repos.SheetUtils
	}
	type args struct {
		documentURI string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *string
	}{
		{
			name:   "Valid hashed ID without view parameter",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/17f77312-e366-809e-9b3c-f94f2eb40700",
			},
			want: tests.StringPtr("17f77312-e366-809e-9b3c-f94f2eb40700"),
		},
		{
			name:   "InValid hashed ID without view parameter",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/17f77312x-e366-809e-9b3c-f94f2eb40700",
			},
			want: nil,
		},
		{
			name:   "Invalid HOST",
			fields: fields{},
			args: args{
				documentURI: "https://example.com/database?query=param",
			},
			want: nil,
		},
		{
			name:   "InValid Database ID",
			fields: fields{},
			args: args{
				documentURI: "https://notion.so/database?query=param",
			},
			want: nil,
		},
		{
			name:   "Invalid Document URI with additional segments",
			fields: fields{},
			args: args{
				documentURI: "https://notion.so/database/path//67890?other=param",
			},
			want: nil,
		},
		{
			name:   "Empty document URI",
			fields: fields{},
			args: args{
				documentURI: "",
			},
			want: nil,
		},

		{
			name:   "Valid Notion URI",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/workspace_name/123456789012345678901234567890123456?v=67890",
			},
			want: nil,
		},
		{
			name:   "Valid Notion URI with long workspace name",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/long_workspace_name_here/123456789012345678901234567890123456?v=54321",
			},
			want: nil,
		},
		{
			name:   "Empty document URI",
			fields: fields{},
			args: args{
				documentURI: "",
			},
			want: nil,
		},
		{
			name:   "Invalid Notion URI format (missing database ID)",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/workspace_name/?v=67890",
			},
			want: nil,
		},

		{
			name:   "Empty document URI",
			fields: fields{},
			args: args{
				documentURI: "",
			},
			want: nil,
		},
		{
			name:   "Valid Notion URI with workspace name",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/workspace_name/17f77312e366809e9b3cf94f2eb40700?v=67890",
			},
			want: tests.StringPtr("17f77312e366809e9b3cf94f2eb40700"),
		},
		{
			name:   "Valid Notion URI without version name",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/workspace_name/17f77312e366809e9b3cf94f2eb40700",
			},
			want: tests.StringPtr("17f77312e366809e9b3cf94f2eb40700"),
		},
		{
			name:   "Valid Notion URI with hashed database ID",
			fields: fields{},
			args: args{
				documentURI: "https://www.notion.so/17f77312e366809e9b3cf94f2eb40700?v=63606b9f1c9149d8ad092a38e148d579",
			},
			want: tests.StringPtr("17f77312e366809e9b3cf94f2eb40700"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ActionsServiceImpl{
				RedisRepo:             tt.fields.RedisRepo,
				BrokerActionsRepo:     tt.fields.BrokerActionsRepo,
				BrokerCredentialsRepo: tt.fields.BrokerCredentialsRepo,
				HTTPRepo:              tt.fields.HTTPRepo,
				CredentialHTTP:        tt.fields.CredentialHTTP,
				ActionsNotion:         tt.fields.ActionsNotion,
				TokenAuth:             tt.fields.TokenAuth,
				SheetUtils:            tt.fields.SheetUtils,
			}
			got := a.GetDatabaseID(tt.args.documentURI)
			if got == nil && tt.want == nil {
				return
			}
			if *got != *tt.want {
				t.Errorf("ActionsServiceImpl.GetDatabaseID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionsServiceImpl_SerializeNotionContent(t *testing.T) {
	type fields struct {
		RedisRepo             repos.ActionsRedisRepoInterface
		BrokerActionsRepo     repos.ActionsBrokerRepository
		BrokerCredentialsRepo repos.CredentialBrokerRepository
		HTTPRepo              repos.ActionsHTTPRepository
		CredentialHTTP        repos.CredentialHTTPRepository
		ActionsNotion         repos.TransformNotion
		TokenAuth             repos.TokenAuth
		SheetUtils            repos.SheetUtils
	}
	type args struct {
		headers    *[]string
		arrContent *[][]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *[]byte
	}{
		{
			name:   "Basic usage",
			fields: fields{},
			args: args{
				headers: &[]string{"Header1", "Header2", "Header3"},
				arrContent: &[][]string{
					{"Row1Col1", "Row1Col2", "Row1Col3"},
					{"Row2Col1", "Row2Col2", "Row2Col3"},
				},
			},
			want: func() *[]byte {
				data := &models.ProcessedNotionData{
					Headers: []string{"Header1", "Header2", "Header3"},
					ContentRows: [][]string{
						{"Row1Col1", "Row1Col2", "Row1Col3"},
						{"Row2Col1", "Row2Col2", "Row2Col3"},
					},
				}
				jsonData, _ := json.Marshal(data)
				return &jsonData
			}(),
		},
		{
			name:   "Empty headers and content",
			fields: fields{},
			args: args{
				headers:    &[]string{},
				arrContent: &[][]string{},
			},
			want: func() *[]byte {
				data := &models.ProcessedNotionData{
					Headers:     []string{},
					ContentRows: [][]string{},
				}
				jsonData, _ := json.Marshal(data)
				return &jsonData
			}(),
		},
		{
			name:   "One header and one row",
			fields: fields{},
			args: args{
				headers:    &[]string{"Header1"},
				arrContent: &[][]string{{"Row1Col1"}},
			},
			want: func() *[]byte {
				data := &models.ProcessedNotionData{
					Headers:     []string{"Header1"},
					ContentRows: [][]string{{"Row1Col1"}},
				}
				jsonData, _ := json.Marshal(data)
				return &jsonData
			}(),
		},
		{
			name:   "Headers and no content",
			fields: fields{},
			args: args{
				headers:    &[]string{"Header1", "Header2"},
				arrContent: &[][]string{},
			},
			want: func() *[]byte {
				data := &models.ProcessedNotionData{
					Headers:     []string{"Header1", "Header2"},
					ContentRows: [][]string{},
				}
				jsonData, _ := json.Marshal(data)
				return &jsonData
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ActionsServiceImpl{
				RedisRepo:             tt.fields.RedisRepo,
				BrokerActionsRepo:     tt.fields.BrokerActionsRepo,
				BrokerCredentialsRepo: tt.fields.BrokerCredentialsRepo,
				HTTPRepo:              tt.fields.HTTPRepo,
				CredentialHTTP:        tt.fields.CredentialHTTP,
				ActionsNotion:         tt.fields.ActionsNotion,
				TokenAuth:             tt.fields.TokenAuth,
				SheetUtils:            tt.fields.SheetUtils,
			}
			if got := a.SerializeNotionContent(tt.args.headers, tt.args.arrContent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActionsServiceImpl.SerializeNotionContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
