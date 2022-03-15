package providers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/providers"
)

func TestIsURL(t *testing.T) {
	t.Run("Correct URL", func(t *testing.T) {
		urlStr := "http://localhost:9001"

		u, ok := providers.IsURL(urlStr)

		assert.EqualValues(t, true, ok)
		assert.EqualValues(t, "http", u.Scheme)
		assert.EqualValues(t, "localhost:9001", u.Host)
		assert.EqualValues(t, "9001", u.Port())
	})

	t.Run("Bad URLs", func(t *testing.T) {
		badURLs := []string{
			"http:::/not.valid/a//a??a?b=&&c#hi",
			"http//google.com",
			"google.com",
			"/foo/bar",
			"http://",
		}

		for i := range badURLs {
			_, ok := providers.IsURL(badURLs[i])

			assert.EqualValues(t, false, ok, badURLs[i])
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("Validate urls and 1 is generated correctly", func(t *testing.T) {
		urls := []string{
			"us=http:::/not.valid/a//a??a?b=&&c#hi",
			"us=http//google.com",
			"us=google.com",
			"us=http://localhost:9001",
			"us=/foo/bar",
			"us=http://",
		}

		p := providers.New(urls)

		// len must be 1
		assert.Len(t, p, 1)

		u := p["us"].URL

		assert.EqualValues(t, "http", u.Scheme)
		assert.EqualValues(t, "localhost:9001", u.Host)
		assert.EqualValues(t, "9001", u.Port())
	})

	t.Run("Validate urls and 3 is generated correctly", func(t *testing.T) {
		urls := []string{
			"us=http:::/not.valid/a//a??a?b=&&c#hi",
			"us=http//google.com",
			"us=google.com",
			"us=http://localhost:9001",
			"us=/foo/bar",
			"ur=http://localhost:9002",
			"us=http://",
			"mx=http://localhost:9003",
		}

		ps := providers.New(urls)

		// len must be 3
		assert.Len(t, ps, 3)

		p := ps["us"]
		// Testintg the ID
		assert.EqualValues(t, "us", p.ID)
		// Testing the URL
		assert.EqualValues(t, "http", p.URL.Scheme)
		assert.EqualValues(t, "localhost:9001", p.URL.Host)
		assert.EqualValues(t, "9001", p.URL.Port())
		// Testing the HTTP Client
		assert.NotNil(t, p.Client)

		p = ps["ur"]
		// Testintg the ID
		assert.EqualValues(t, "ur", p.ID)
		// Testing the URL
		assert.EqualValues(t, "http", p.URL.Scheme)
		assert.EqualValues(t, "localhost:9002", p.URL.Host)
		assert.EqualValues(t, "9002", p.URL.Port())
		// Testing the HTTP Client
		assert.NotNil(t, p.Client)

		p = ps["mx"]
		// Testintg the ID
		assert.EqualValues(t, "mx", p.ID)
		// Testing the URL
		assert.EqualValues(t, "http", p.URL.Scheme)
		assert.EqualValues(t, "localhost:9003", p.URL.Host)
		assert.EqualValues(t, "9003", p.URL.Port())
		// Testing the HTTP Client
		assert.NotNil(t, p.Client)
	})
}
