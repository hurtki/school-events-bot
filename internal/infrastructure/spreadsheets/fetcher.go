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
}

func NewDocsFetcher(cl *http.Client) *DocsFetcher {
	return &DocsFetcher{
		cl: cl,
	}
}

func (f *DocsFetcher) FetchXLSX(ctx context.Context, documentID string) (io.ReadCloser, error) {
	url := baseSpreadsheetURL + documentID + xlsxExportSuffix
	req, err := http.NewRequestWithContext(ctx,
		"GET",
		url,
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
		return nil, fmt.Errorf("not expected status code: %d url: %s", res.StatusCode, url)
	}
	return res.Body, err
}
