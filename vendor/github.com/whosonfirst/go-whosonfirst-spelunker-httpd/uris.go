package httpd

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type URIs struct {
	// WWW/human-readable
	Id             string   `json:"id"`
	IdAlt          []string `json:"id_alt"`
	Descendants    string   `json:"descendants"`
	DescendantsAlt []string `json:"descendants_alt"`
	Index string `json:"index"`
	Recent string `json:"recent"`
	Search         string   `json:"search"`
	About          string   `json:"about"`

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
		
		Id:          "/id/",
		Recent: "/recent/",
		Descendants: "/descendants/",
		DescendantsAlt: []string{
			"/id/{id}/descendants/",
		},
		Index: "/",
		Search: "/search",
		About:  "/about",

		// Static Assets
		Static: "/static/",

		// API/machine-readable
		GeoJSON:   "/geojson/",
		GeoJSONLD: "/geojsonld/",
		NavPlace:  "/navplace/",
		Select:    "/select/",
		SPR:       "/spr/",
		SVG:       "/svg/",
	}

	return uris_table
}
