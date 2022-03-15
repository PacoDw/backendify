package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cast"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/cache"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/providers"
)

const (
	// V1 represents the content type to point to v1 of the provider endpoint.
	HeaderV1 = "application/x-company-v1"

	// V2 represents the content type to point to v2 of the provider endpoint.
	HeaderV2 = "application/x-company-v2"
)

// containLegacyHeaders validates that the response of the legacy service contains
// the legacy headers.
func containLegacyHeaders(headers []string) bool {
	var (
		legacyHeaders = []string{HeaderV1, HeaderV2}
		found         bool
	)

	for i := range headers {
		for k := range legacyHeaders {
			if headers[i] == legacyHeaders[k] {
				found = true

				break
			}
		}
	}

	return found
}

func CompanyRoute(pdrs providers.Providers, c *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			// So far at this point we know that these values are filled.
			id  = cast.ToString(r.Context().Value(CompanyID))
			iso = cast.ToString(r.Context().Value(CountryCode))
		)

		p, ok := pdrs[iso]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
		}

		// Adding the companies path and id of the current request to preparate the next request.
		p.URL.Path = fmt.Sprintf("companies/%s", id)

		// preparing the request
		req, _ := http.NewRequestWithContext(r.Context(), http.MethodGet, p.URL.String(), http.NoBody)

		// Making request to the legacy services
		res, err := p.Client.Do(req)
		// if there is an error then get the last known data from the cache
		// but if the cache doesnt contains data then return the error.
		// NOTE: there is .50 second to wait until the legacy service responds if not response
		// then error is going to trigger to get data from cache.
		if err != nil {
			v, found := c.Get(id)
			if !found {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			if _, err := w.Write(v.([]byte)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
		defer res.Body.Close()

		// verify if the response contains the correct headers if not return an error 500.
		// NOTE: if this error appears a lot means that the legacy headers has changed.
		if ok := containLegacyHeaders(res.Header.Values("Content-Type")); !ok {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("Couldn't parse response body. %+v", err)
		}

		cresp := &CompanyResponse{}
		if err := json.Unmarshal(body, cresp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		// convert it in bytes (json)
		result := cresp.ToJSON()

		// store the new value from the service into the cache
		c.StoreOrLoad(id, result)

		// return the value
		if _, err := w.Write(result); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
