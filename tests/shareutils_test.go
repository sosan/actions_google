package tests

import (
	"actions_google/mocks"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"google.golang.org/api/sheets/v4"
)

func TestCreateSheetsService(t *testing.T) {
	t.Run("Created SheetService", func(t *testing.T) {
		mockSheetUtils := &mocks.SheetUtils{}

		ctx := context.Background()
		client := &http.Client{}
		expectedService := &sheets.Service{}

		mockSheetUtils.On(
			"CreateSheetsService",
			ctx,
			client,
		).Return(expectedService, nil)

		result, err := mockSheetUtils.CreateSheetsService(ctx, client)

		assert.NoError(t, err)
		assert.Same(t, expectedService, result)
	})
}

func TestGetAllContentFromGoogleSheets2(t *testing.T) {
	t.Run("Obtain google sheets", func(t *testing.T) {
		mockSheetUtils := &mocks.SheetUtils{}

		expectedDoc := "mi-documento"
		expectedClient := &http.Client{}
		expectedActionID := "accion-123"
		expectedResponse := &sheets.ValueRange{Values: [][]interface{}{{"A1", "B1"}}}

		mockSheetUtils.On(
			"GetAllContentFromGoogleSheets",
			&expectedDoc,
			expectedClient,
			&expectedActionID,
		).Return(expectedResponse, nil)

		result, err := mockSheetUtils.GetAllContentFromGoogleSheets(
			&expectedDoc,
			expectedClient,
			&expectedActionID,
		)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockSheetUtils.AssertExpectations(t)
	})
}
