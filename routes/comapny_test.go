package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/cache"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/providers"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/routes"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/server"
)

func serverMock(t *testing.T, latency time.Duration, wrongLegacyHeaders bool) *httptest.Server {
	t.Helper()

	handler := http.NewServeMux()

	// We will test with the v1 version of the provider server
	handler.HandleFunc("/companies/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if wrongLegacyHeaders {
			w.Header().Add("Content-Type", "application/x-company-v343")
		} else {
			w.Header().Add("Content-Type", "application/x-company-v1")
		}

		time.Sleep(latency)

		_, err := w.Write([]byte(`{
			"cn": "Company Name",
			"created_on": "2012-03-14T16:46:45.019018-06:00",
			"closed_on": "2024-03-14T16:46:45.019018-06:00"
		  }`))

		assert.NoError(t, err)
	})

	// We will test with the v2 version of the provider server
	handler.HandleFunc("/companies/v2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if wrongLegacyHeaders {
			w.Header().Add("Content-Type", "application/x-company-v343")
		} else {
			w.Header().Add("Content-Type", "application/x-company-v2")
		}

		time.Sleep(latency)

		_, err := w.Write([]byte(`{
			"company_name":"Company Name",
			"tin":"V12345678",
			"dissolved_on":"2024-03-14T16:46:45.019018-06:00"
		 }`))

		assert.NoError(t, err)
	})

	srv := httptest.NewServer(handler)

	return srv
}

func TestCompanyRoute(t *testing.T) {
	var (
		latency                = 0 * time.Second
		withWrongLegacyHeaders = false
	)

	srv := serverMock(t, latency, withWrongLegacyHeaders)

	tests := []struct {
		name         string
		providers    providers.Providers
		cache        *cache.Cache
		rec          *httptest.ResponseRecorder
		req          *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success V1",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company?id=v1&county_iso=us", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`,
		},
		{
			name:         "Success V2",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company?id=v2&county_iso=us", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`,
		},
		{
			name:         "Bad request",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company/county_iso=us", nil),
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server.ValidateQueryParametersMiddleware([]routes.RequiredQueryParameter{routes.CompanyID, routes.CountryCode})(
				http.HandlerFunc(routes.CompanyRoute(test.providers, cache.New(0, 0))),
			).ServeHTTP(test.rec, test.req)

			// validate the body
			assert.EqualValues(t, test.expectedBody, test.rec.Body.String())

			// validate status code
			assert.EqualValues(t, test.expectedCode, test.rec.Code)
		})
	}
}

func TestCompanyRoute_WithCache(t *testing.T) {
	var (
		latency                = 0 * time.Second
		withWrongLegacyHeaders = false
	)

	srv := serverMock(t, latency, withWrongLegacyHeaders)

	tests := []struct {
		name         string
		providers    providers.Providers
		cache        *cache.Cache
		rec          *httptest.ResponseRecorder
		req          *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success V1",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0).ChainStoreOrLoad("v1", []byte(`{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`)),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company?id=v1&county_iso=us", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`,
		},
		{
			name:         "Success V2",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0).ChainStoreOrLoad("v2", []byte(`{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`)),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company?id=v2&county_iso=us", nil),
			expectedCode: http.StatusOK,
			expectedBody: `{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`,
		},
		{
			name:         "Bad request",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company/county_iso=us", nil),
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server.ValidateQueryParametersMiddleware([]routes.RequiredQueryParameter{routes.CompanyID, routes.CountryCode})(
				http.HandlerFunc(routes.CompanyRoute(test.providers, test.cache)),
			).ServeHTTP(test.rec, test.req)

			// validate status code
			assert.EqualValues(t, test.expectedCode, test.rec.Code)

			// validate the body
			assert.EqualValues(t, test.expectedBody, test.rec.Body.String())
		})
	}
}

func TestCompanyRoute_WithWrongLegacyHeaders(t *testing.T) {
	var (
		latency                = 0 * time.Second
		withWrongLegacyHeaders = true
	)

	srv := serverMock(t, latency, withWrongLegacyHeaders)

	tests := []struct {
		name         string
		providers    providers.Providers
		cache        *cache.Cache
		rec          *httptest.ResponseRecorder
		req          *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Success V1",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0).ChainStoreOrLoad("v1", []byte(`{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`)),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company?id=v1&county_iso=us", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:         "Success V2",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0).ChainStoreOrLoad("v2", []byte(`{"name":"Company Name","actived":true,"active_until":"2024-03-14T16:46:45.019018-06:00"}`)),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company?id=v2&county_iso=us", nil),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:         "Bad request",
			providers:    providers.New([]string{fmt.Sprintf("us=%s", srv.URL)}),
			cache:        cache.New(0, 0),
			rec:          httptest.NewRecorder(),
			req:          httptest.NewRequest("GET", "/company/county_iso=us", nil),
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server.ValidateQueryParametersMiddleware([]routes.RequiredQueryParameter{routes.CompanyID, routes.CountryCode})(
				http.HandlerFunc(routes.CompanyRoute(test.providers, test.cache)),
			).ServeHTTP(test.rec, test.req)

			// validate status code
			assert.EqualValues(t, test.expectedCode, test.rec.Code)

			// validate the body
			assert.EqualValues(t, test.expectedBody, test.rec.Body.String())
		})
	}
}
