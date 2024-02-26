package www

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func FiltersFromRequest(ctx context.Context, req *http.Request, params []string) ([]spelunker.Filter, error) {

	filters := make([]spelunker.Filter, 0)

	for _, p := range params {

		switch p {
		case "country":

			country, err := sanitize.GetString(req, "country")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?placetype= query parameter, %w", err)
			}

			if country != "" {

				country_f, err := spelunker.NewCountryFilterFromString(ctx, country)

				if err != nil {
					return nil, fmt.Errorf("Failed to create country filter from string '%s', %w", country, err)
				}

				filters = append(filters, country_f)
			}

		case "placetype":

			placetype, err := sanitize.GetString(req, "placetype")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?placetype= query parameter, %w", err)
			}

			if placetype != "" {

				placetype_f, err := spelunker.NewPlacetypeFilterFromString(ctx, placetype)

				if err != nil {
					return nil, fmt.Errorf("Failed to create placetype filter from string '%s', %w", placetype, err)
				}

				filters = append(filters, placetype_f)
			}

		default:
			return nil, fmt.Errorf("Invalid or unsupported parameter, %s", p)
		}
	}

	return filters, nil
}
