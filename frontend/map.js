mapboxgl.accessToken = 'pk.eyJ1Ijoia2VyZWRuaXkiLCJhIjoiY2poczdicmh3MGxodjNkcHBndzZhMGpoeSJ9.PTnSTWrvBoUtTGvqxgtiLw';

var dataArray;
var countryArray = new Array();

Papa.parse('https://storage.googleapis.com/topify/data', {
    download:true,
    complete: function(results) {

                dataArray = results.data;
    
                for(var i = 0; i < dataArray.length; i++){
                    var country = {
                        name: dataArray[i][0],
                        link: dataArray[i][1],
                }

        countryArray[i]=country;
        }
        console.log(countryArray)
    }
});

var map = new mapboxgl.Map({
    container: 'map', 
    style: 'mapbox://styles/keredniy/cjhs7mdsx6q8n2snwn7hrogq6', 
    zoom: 1.4,
    minZoom: 1.2,
    maxZoom: 2,
    pitchWithRotate: false,
    dragRotate: false,
});

map.on('load', function () {
        map.addSource("countries", {
        "type": "geojson",
        "data": 'countries_v2.geojson'
    });


    map.addLayer({
        "id": "borders",   
        "type": "line",
        "source": "countries",
        "layout": {},
        "paint": {
        "line-color": "#1ED760",
        "line-width": 0.9
        }
    });

    map.addLayer({
        "id": "countries",
        "type": "fill",
        "source": "countries",
        "layout": {},
        "paint": {
        "fill-color": "#000",
        "fill-opacity": 0
        }
    });

    map.addLayer({
        "id": "countries_hover",
        "type": "fill",
        "source": "countries",
        "layout": {},
        "paint": {
            "fill-color": "#1ED760",
            "fill-opacity": 1
        },
        "filter": ["==", "name", ""]
    });

    var popup;

    map.on("mousemove", "countries", function(e) {

        for(var i = 0; i < countryArray.length; i++){
            if(countryArray[i].name == e.features[0].properties.iso_a2){
                map.setFilter("countries_hover", ["==", "name", e.features[0].properties.name]);
            }
        }
    });

    map.on("mouseleave", "countries", function() {
        map.setFilter("countries_hover", ["==", "name", ""]);
    });

    map.on("click", "countries", function(e) {

        popup = new mapboxgl.Popup({
            closeButton: false,
            closeOnClick: true
    });

        for(var i = 0; i < countryArray.length; i++){
            if(countryArray[i].name == e.features[0].properties.iso_a2){
                    map.setFilter("countries_hover", ["==", "name", e.features[0].properties.name]);
                
                    popup.setLngLat(e.lngLat)
                    .setHTML('<iframe src="https://open.spotify.com/embed?uri=spotify:track:'+countryArray[i].link+'"width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>')
                    .addTo(map);
            }
        }
    });
});