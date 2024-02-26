window.addEventListener("load", function load(event){

    
    var facets_wrapper = document.querySelector("#whosonfirst-facets");

    if (! facets_wrapper){
	console.log("NOPE");
	return;
    }

    var current_url = facets_wrapper.getAttribute("data-current-url");
    var facets_url = facets_wrapper.getAttribute("data-facets-url");    

    if ((! current_url) || (! facets_url)){
	return;
    }

    var draw_facets = function(rsp){

	var f = rsp.facet.property;
	var id = "#whosonfirst-facets-" + f;

	var el = document.querySelector(id);

	if (! el){
	    console.log("Unable to find facet wrapper", id);
	    return;
	}

	var label = document.createElement("h3");
	label.appendChild(document.createTextNode(f));
	
	var ul = document.createElement("ul");
	ul.setAttribute("class", "whosonfirst-facets");
	
	var results = rsp.results;
	var count = results.length;

	for (var i=0; i < count; i++){

	    var a = document.createElement("a");
	    // To do: Use proper query builder methods
	    a.setAttribute("href", current_url + "?" + encodeURIComponent(f) + "=" + encodeURIComponent(results[i].key));
	    a.setAttribute("class", "hey-look");
	    a.appendChild(document.createTextNode(results[i].key));

	    var sm = document.createElement("small");
	    sm.appendChild(document.createTextNode(results[i].count));
		
	    var item = document.createElement("li");
	    item.appendChild(a);
	    item.appendChild(sm);

	    ul.appendChild(item);
	}

	el.appendChild(label);
	el.appendChild(ul);
    };
    
    var fetch_facet = function(f){

	var url = facets_url + "?facet=" + f;

	fetch(url)
	    .then((rsp) => rsp.json())
	    .then((data) => {

		var count = data.length;

		for (var i=0; i < count; i++){
		    draw_facets(data[i]);
		}
		
	    }).catch((err) => {
		console.log("SAD", f, err);
	    });
    };
    
    var facets = facets_wrapper.getAttribute("data-facets");
    facets = facets.split(",");

    var count_facets = facets.length;

    for (var i=0; i < count_facets; i++){

	var f = facets[i];
	
	var el = document.createElement("div");
	el.setAttribute("id", "whosonfirst-facets-" + f);
	facets_wrapper.appendChild(el);

	fetch_facet(f);
    }

});
