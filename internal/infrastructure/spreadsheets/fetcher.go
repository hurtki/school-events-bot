package spreadsheets

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var (
	baseSpreadsheetURL = "https://docs.google.com/spreadsheets/d/"
	xlsxExportSuffix   = "/export?format=xlsx"
)

type DocsFetcher struct {
	cl *http.Client

	documentID string
}

func NewDocsFetcher(documentID string, cl *http.Client) *DocsFetcher {
	return &DocsFetcher{
		cl:         cl,
		documentID: documentID,
	}
}

func (f *DocsFetcher) FetchXLSX(ctx context.Context) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx,
		"GET",
		baseSpreadsheetURL+f.documentID+xlsxExportSuffix,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}
	res, err := f.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't execute request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not expected status code: %d", res.StatusCode)
	}
	return res.Body, err
}
