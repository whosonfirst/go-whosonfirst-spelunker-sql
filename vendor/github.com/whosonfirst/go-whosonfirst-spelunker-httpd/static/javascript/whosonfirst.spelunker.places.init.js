window.addEventListener("load", function load(event){

    var places = document.querySelectorAll(".whosonfirst-places-list li");

    if (! places){
	console.log("No places");
	return;
    }
    
    var count_places = places.length;

    if (count_places == 0){
	return;
    }
    
    var coords = [];
    var names = [];
    
    for (var i=0; i < count_places; i++) {

	var el = places[i];
	// console.log(el);
	
	var lat = parseFloat(el.getAttribute("data-latitude"));
	var lon = parseFloat(el.getAttribute("data-longitude"));	

	if ((! lat) || (!lon)){
	    console.log("Invalid coordinates", i, lat, lon);
	    continue;
	}

	var n = el.querySelector(".wof-place-name");

	if ((! n) || (n.innerText == "")){
	    console.log("Invalid name", i);
	    continue;
	}

	coords.push([ lon, lat ]);
	names[ JSON.stringify(coords[i]) ] = n.innerText;
    }

    var f = {
	"type": "Feature",
	"properties": {
	    "lflt:label_names": names,
	},
	"geometry": {
	    "type": "MultiPoint",
	    "coordinates": coords,
	},
    };

    var map_el = document.querySelector("#map");
    map_el.style.display = "block";
    
    const map = whosonfirst.spelunker.maps.map(map_el);

    switch (coords.length){
	case 0:
	    // Null Island
	    coords.push([ 0.0, 0.0 ]);
	    f.geometry.coords = coords;
	    map.setView([coords[0][1], coords[0][0]], 3);
	    break;
	case 1:
	    // TO DO: set zoom based on placetype or mz:min/max_zoom (requires fetching the record...)
	    map.setView([coords[0][1], coords[0][0]], 12);
	    break;
	default:
	    var bounds = whosonfirst.spelunker.geojson.derive_bounds(f);
	    map.fitBounds(bounds);
	    break;
    }
    
    var pt_handler_layer_args = {
	pane: whosonfirst.spelunker.maps.centroids_pane_name,
	tooltips_pane: whosonfirst.spelunker.maps.tooltips_pane_name,	
    };
		
    var pt_handler = whosonfirst.spelunker.leaflet.handlers.point(pt_handler_layer_args);
    var lbl_style = whosonfirst.spelunker.leaflet.styles.search_centroid();

    var points_layer_args = {
	style: lbl_style,
	pointToLayer: pt_handler,
	pane: whosonfirst.spelunker.maps.centroids_pane_name,
    }
    
    var points_layer = L.geoJSON(f, points_layer_args);
    points_layer.addTo(map);

});
