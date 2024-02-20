package app

import (
	"fmt"
	"io"
	"net/http"
)

type OfacApi struct {
	client *http.Client
}

func (o *OfacApi) GetSdnEntries() (io.ReadCloser, error) {
	const fn = "app.ofac.GetSdnEntries"

	resp, err := o.client.Get("https://www.treasury.gov/ofac/downloads/sdn.xml")
	if err != nil {
		return nil, fmt.Errorf("%s :%w", fn, err)
	}
	return resp.Body, nil
}

func NewOfacAPI(client *http.Client) *OfacApi {
	return &OfacApi{client: client}
}
