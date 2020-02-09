package client

import (
	"encoding/json"
	"testing"
)

func Test_generateFindGeometryUrl(t *testing.T) {
	t.Run("space strings", func(t *testing.T) {
		name := "Mall of scandinavia"
		apiKey := "bob"
		findGeometryUrl := GenerateFindGeometryUrl(name, apiKey)
		expected := "https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=Mall+of+scandinavia&inputtype=textquery&fields=geometry&key=bob"
		if findGeometryUrl != expected {
			t.Errorf("expected: %s; got %s", expected, findGeometryUrl)
		}
	})
	t.Run("special characters", func(t *testing.T) {
		name := "Bianchi Café & Cycles"
		apiKey := "bob"
		findGeometryUrl := GenerateFindGeometryUrl(name, apiKey)
		expected := "https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=Bianchi+Caf%C3%A9+%26+Cycles&inputtype=textquery&fields=geometry&key=bob"
		if findGeometryUrl != expected {
			t.Errorf("expected: %s; got %s", expected, findGeometryUrl)
		}

	})
}

func Test_unmarshalGeometry(t *testing.T) {
	body := []byte(`
{
   "candidates" : [
      {
         "geometry" : {
            "location" : {
               "lat" : 59.33235649999999,
               "lng" : 18.0645449
            },
            "viewport" : {
               "northeast" : {
                  "lat" : 59.33375587989273,
                  "lng" : 18.06524522989272
               },
               "southwest" : {
                  "lat" : 59.33105622010729,
                  "lng" : 18.06254557010728
               }
            }
         }
      }
   ],
   "status" : "OK"
}`)
	lac := new(LocationAPICandidates)
	err := json.Unmarshal(body, &lac)
	if err != nil {
		t.Fatalf("couldn't unmarshal geometryAPI")
	}
	if len(lac.Candidates) < 1 {
		t.Fatalf("should have one geometry")
	}
	expectedLat := 59.33235649999999
	lat := lac.Candidates[0].Geometry.Location.Lat
	if lat != expectedLat {
		t.Fatalf("lat shoud have been %f, it was %f", expectedLat, lat)
	}
	expectedLng := 18.0645449

	lng := lac.Candidates[0].Geometry.Location.Lng
	if lng != expectedLng {
		t.Fatalf("lng shoud have been %f, it was %f", expectedLng, lng)
	}
}
