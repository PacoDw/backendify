package routes_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/routes"
)

func TestUnmarshalJSONWithV1(t *testing.T) {
	t.Run("Success with active as true", func(t *testing.T) {
		expectedTime := time.Now().AddDate(2, 0, 0)

		expected := &routes.CompanyResponse{
			Name:        "Company Name",
			Actived:     pointy.Bool(true),
			ActiveUntil: &expectedTime,
		}

		blob := []byte(fmt.Sprintf(`{"cn":"Company Name","created_on":"2012-11-01T22:08:41+00:00","closed_on":%q}`, expectedTime.Format(time.RFC3339Nano)))

		got := &routes.CompanyResponse{}

		err := json.Unmarshal(blob, got)
		assert.Nil(t, err, err)

		assert.EqualValues(t, expected, got)
	})

	t.Run("Success with active as false", func(t *testing.T) {
		expectedTime := time.Now().AddDate(-2, 0, 0)

		expected := &routes.CompanyResponse{
			Name:        "Company Name",
			Actived:     pointy.Bool(false),
			ActiveUntil: &expectedTime,
		}

		blob := []byte(fmt.Sprintf(`{"cn":"Company Name","created_on":"2012-11-01T22:08:41+00:00","closed_on":%q}`, expectedTime.Format(time.RFC3339Nano)))

		got := &routes.CompanyResponse{}

		err := json.Unmarshal(blob, got)
		assert.Nil(t, err, err)

		assert.EqualValues(t, expected, got)
	})
}

func TestUnmarshalJSONWithV2(t *testing.T) {
	t.Run("Success with active as true", func(t *testing.T) {
		expectedTime := time.Now().AddDate(2, 0, 0)

		expected := &routes.CompanyResponse{
			Name:        "Company Name",
			Actived:     pointy.Bool(true),
			ActiveUntil: &expectedTime,
		}

		blob := []byte(fmt.Sprintf(`{"company_name":"Company Name","tin":"V1234785","dissolved_on":%q}`, expectedTime.Format(time.RFC3339Nano)))

		got := &routes.CompanyResponse{}

		err := json.Unmarshal(blob, got)
		assert.Nil(t, err, err)

		assert.EqualValues(t, expected, got)
	})

	t.Run("Success with active as false", func(t *testing.T) {
		expectedTime := time.Now().AddDate(-2, 0, 0)

		expected := &routes.CompanyResponse{
			Name:        "Company Name",
			Actived:     pointy.Bool(false),
			ActiveUntil: &expectedTime,
		}

		blob := []byte(fmt.Sprintf(`{"company_name":"Company Name","tin":"V1234785","dissolved_on":%q}`, expectedTime.Format(time.RFC3339Nano)))

		got := &routes.CompanyResponse{}

		err := json.Unmarshal(blob, got)
		assert.Nil(t, err, err)

		assert.EqualValues(t, expected, got)
	})
}
