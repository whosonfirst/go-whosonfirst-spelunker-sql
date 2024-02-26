package httpd

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type URIs struct {
	// WWW/human-readable
	Id               string   `json:"id"`
	IdAlt            []string `json:"id_alt"`
	Concordances     string   `json:"concordances"`
	Concordance      string   `json:"concordance"`
	Descendants      string   `json:"descendants"`
	DescendantsAlt   []string `json:"descendants_alt"`
	DescendantsFacet string   `json:"descendants_facet"`
	Index            string   `json:"index"`
	Placetypes       string   `json:"placetypes"`
	Placetype        string   `json:"placetype"`
	Recent           string   `json:"recent"`
	Search           string   `json:"search"`
	About            string   `json:"about"`

	// Static assets
	Static string `json:"static"`

	// API/machine-readable
	GeoJSON      string   `json:"geojson"`
	GeoJSONAlt   []string `json:"geojson_alt"`
	GeoJSONLD    string   `json:"geojsonld"`
	GeoJSONLDAlt []string `json:"geojsonld_alt"`
	NavPlace     string   `json:"navplace"`
	NavPlaceAlt  []string `json:"navplace_alt"`
	Select       string   `json:"select"`
	SelectAlt    []string `json:"select_alt"`
	SPR          string   `json:"spr"`
	SPRAlt       []string `json:"spr_alt"`
	SVG          string   `json:"svg"`
	SVGAlt       []string `json:"svg_alt"`
}

func (u *URIs) ApplyPrefix(prefix string) error {

	val := reflect.ValueOf(*u)

	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		v := field.String()

		if v == "" {
			continue
		}

		if strings.HasPrefix(v, prefix) {
			continue
		}

		new_v, err := url.JoinPath(prefix, v)

		if err != nil {
			return fmt.Errorf("Failed to assign prefix to %s, %w", v, err)
		}

		reflect.ValueOf(u).Elem().Field(i).SetString(new_v)
	}

	return nil
}

func DefaultURIs() *URIs {

	uris_table := &URIs{

		// WWW/human-readable

		Index:            "/",
		Search:           "/search",
		About:            "/about",
		Placetypes:       "/placetypes",
		Placetype:        "/placetypes/{placetype}",
		Concordances:     "/concordances/",
		Concordance:      "/concordances/{concordance}",
		Recent:           "/recent/",
		Id:               "/id/{id}",
		Descendants:      "/id/{id}/descendants",
		DescendantsFacet: "/id/{id}/descendants/facet",

		// Static Assets
		Static: "/static/",

		// API/machine-readable
		GeoJSON: "/geojson/",
		GeoJSONAlt: []string{
			"/id/{id}/geojson",
		},
		GeoJSONLD: "/geojsonld/",
		GeoJSONLDAlt: []string{
			"/id/{id}/geojsonld",
		},
		NavPlace: "/navplace/",
		NavPlaceAlt: []string{
			"/id/{id}/navplace",
		},
		Select: "/select/",
		SelectAlt: []string{
			"/id/{id}/select",
		},
		SPR: "/spr/",
		SPRAlt: []string{
			"/id/{id}/spr",
		},
		SVG: "/svg/",
		SVGAlt: []string{
			"/id/{id}/svg",
		},
	}

	return uris_table
}

func URIForId(uri string, id int64) string {
	return ReplaceAll(uri, "{id}", id)
}

func ReplaceAll(input string, pattern string, value any) string {
	str_value := fmt.Sprintf("%v", value)
	return strings.Replace(input, pattern, str_value, -1)
}
