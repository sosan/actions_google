package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetUtilsImpl struct{}

func NewSheetUtilsImpl() *SheetUtilsImpl {
	return &SheetUtilsImpl{}
}

func (s *SheetUtilsImpl) GetAllContentFromGoogleSheets(document *string, client *http.Client, actionID *string) (*sheets.ValueRange, error) {
	ctx := context.Background()
	sheetsService, err := s.CreateSheetsService(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("ERROR | not possible to initialize google sheets service: %v", err)
	}
	spreadsheetID := s.GetSpreadsheetID(document)
	if spreadsheetID == nil {
		return nil, fmt.Errorf("ERROR | cannot get spreadsheetID")
	}
	if *spreadsheetID == "" {
		return nil, fmt.Errorf("ERROR | cannot get spreadsheetID")
	}

	spreadsheet, err := s.GetSpreadsheet(ctx, sheetsService, *spreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("ERROR | cannot fetch spreadsheetID: %s error: %v", *spreadsheetID, err)
	}
	values, err := s.GetValuesFromSheet(spreadsheet, sheetsService, spreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("ERROR | cannot get values for spreadsheetID: %s error: %v  for actioid: %s", *spreadsheetID, err, *actionID)
	}
	return values, nil
}

func (s *SheetUtilsImpl) GetSpreadsheetID(documentURI *string) *string {
	splitted := strings.Split(*documentURI, "/")
	// clean and check if uri not contains ?
	// fixed position in array 5
	cleanedArray := strings.Split(splitted[5], "?")
	if len(cleanedArray) > 0 {
		return &cleanedArray[0]
	}
	return &splitted[5]
}

func (s *SheetUtilsImpl) GetValuesFromSheet(sheets *sheets.Spreadsheet, sheetsService *sheets.Service, spreadsheetID *string) (*sheets.ValueRange, error) {
	for _, sheet := range sheets.Sheets {
		// properties := sheet.Properties
		// gridProperties := properties.GridProperties
		// dinamic range
		// readRange := fmt.Sprintf("%s", sheetName)
		// readRange := fmt.Sprintf("%s!A1:%s%d", sheetName, getColumnName(gridProperties.ColumnCount), gridProperties.RowCount)
		sheetName := sheet.Properties.Title
		readRange := sheetName

		values, err := sheetsService.Spreadsheets.Values.Get(*spreadsheetID, readRange).Do()
		return values, err
	}
	return nil, fmt.Errorf("ERROR | not sheets")
}

func (s *SheetUtilsImpl) CreateSheetsService(ctx context.Context, client *http.Client) (*sheets.Service, error) {
	return sheets.NewService(ctx, option.WithHTTPClient(client))
}

func (s *SheetUtilsImpl) GetSpreadsheet(ctx context.Context, srv *sheets.Service, spreadsheetID string) (*sheets.Spreadsheet, error) {
	return srv.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
}
