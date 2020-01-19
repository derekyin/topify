mapboxgl.accessToken = 'pk.eyJ1Ijoia2VyZWRuaXkiLCJhIjoiY2s1bGNrN2lpMG5pNjNnbzc4dnh4dzRhYyJ9.i_XmhuMlaHmGGpyvaygBzw';

var dataArray;
var countryArray = {};

Papa.parse('https://storage.googleapis.com/topify-data/topify-list.csv', {
    download: true,
    complete: function(results) {

        dataArray = results.data;

        for(var i = 0; i < dataArray.length; i++){
            var country = {
                name: dataArray[i][0],
                link: dataArray[i][1],
            }
            countryArray[country.name] = country;
        }
    }
});

var map = new mapboxgl.Map({
    container: 'map', 
    style: 'mapbox://styles/keredniy/cjhs7mdsx6q8n2snwn7hrogq6?optimize=true', 
    zoom: 1.4,
    minZoom: 1.2,
    maxZoom: 2,
    pitchWithRotate: false,
    dragRotate: false,
});

map.on('load', function () {
    map.addSource("countries", {
        "type": "geojson",
        "data": 'https://storage.googleapis.com/topify-data/countries.geojson'
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

    map.on("click", "countries", function(e) {
        popup = new mapboxgl.Popup({
            closeButton: false,
            closeOnClick: true
        });

        if(e.features[0].properties.iso_a2 in countryArray){       
            popup.setLngLat(e.lngLat)
            .setHTML('<iframe src="https://open.spotify.com/embed?uri=spotify:track:'+countryArray[e.features[0].properties.iso_a2].link+'"width="300" height="80" frameborder="0" allowtransparency="true" allow="encrypted-media"></iframe>')
            .addTo(map);
        }else{
            popup.setLngLat(e.lngLat)
            .setHTML('<h1>NO DATA</h1>')
            .addTo(map)
        }
    });
});