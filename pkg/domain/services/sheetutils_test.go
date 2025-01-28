package services

import (
	"reflect"
	"testing"
)

func TestSheetUtilsImpl_GetSpreadsheetID(t *testing.T) {
	type args struct {
		documentURI *string
	}
	tests := []struct {
		name string
		s    *SheetUtilsImpl
		args args
		want *string
	}{
		{
			name: "Valid URI with query parameters",
			s:    &SheetUtilsImpl{},
			args: args{documentURI: func() *string {
				s := "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit?usp=sharing"
				return &s
			}()},
			want: func() *string { s := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"; return &s }(),
		},
		{
			name: "Valid URI without query parameters",
			s:    &SheetUtilsImpl{},
			args: args{documentURI: func() *string {
				s := "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit"
				return &s
			}()},
			want: func() *string { s := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"; return &s }(),
		},
		{
			name: "Invalid URI",
			s:    &SheetUtilsImpl{},
			args: args{documentURI: func() *string { s := "https://docs.google.com/spreadsheets/d/"; return &s }()},
			want: func() *string { s := ""; return &s }(),
		},
		{
			name: "Clean URI",
			s:    &SheetUtilsImpl{},
			args: args{documentURI: func() *string {
				s := "https://docs.google.com/spreadsheets/d/sdfasdfasdadsfasdff?asdfkjaslkdjfasdflkjasdf?ajsdfhaksdf"
				return &s
			}()},
			want: func() *string { s := "sdfasdfasdadsfasdff"; return &s }(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SheetUtilsImpl{}
			if got := s.GetSpreadsheetID(tt.args.documentURI); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SheetUtilsImpl.GetSpreadsheetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
