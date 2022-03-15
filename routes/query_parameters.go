package routes

// QueryParameter represents the required parameter that could be use for some
// other endpoints.
type RequiredQueryParameter string

const (
	// CompanyID represents the id of a company.
	CompanyID RequiredQueryParameter = RequiredQueryParameter("id")

	// CountryCode represents the country iso, e.g.: us, ur, etc.
	CountryCode RequiredQueryParameter = RequiredQueryParameter("county_iso")
)
