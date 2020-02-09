## Google Maps Place Api with GO

### API

#### Endpoints

The API has two endpoints, one to query the locations based on a central point, 
and type of location.

- `GET /?place=<string>&type=<string>`

    supported types:
    - restaurant
    - store
    - bicycle_store 
    
- `GET /all`

    Return all elements in the database

#### Design

The server looks for the environment variable **API_KEY** at startup, in case where can't find it
the server will terminate.

The code gets the place string, makes a request to get the coordinates (lat and lng)
and uses to get a list of the locations of the provided type.
The places are stored in the database (in memory) and next time the request comes
will fetch data from the database instead of sending request to the Google Maps API.


#### Known issues

If the API key is not a premium one, it could happen that it will not be able to fetch 
the items due to **OVER_QUERY_LIMIT**, this information is not given to the user, but logged.
The user will receive a 500 Status Code response.


### Improvements

#### Testing

There are tests in the code, but the client communication with the google places api
needs testing. Also the HttpHandlers.

#### OVER_QUERY_LIMIT Handling

Right now there is not a good handling when it comes to reaching the query limit. One 
solution might be adding the request in a channel and triggering every x seconds, and populating 
the database in the background.


### Build

#### locally 
    
    `make dev`
   or 
   
    `go build .` 
#### docker
    `make docker`
   or
   
    `docker build -t maps-api-go .`
    
### Run

- with docker:
    `docker run -e API_KEY='YOUR_API_KEY' -p 8080:8080 maps-api-go:latest`
- native
    `./google-maps-places-go`
