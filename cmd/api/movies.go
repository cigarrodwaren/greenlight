package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.example.org/internal/data"
	"greenlight.example.org/internal/validator"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title string `json:"title"`
		Year  int32  `json:"year"`
		// Make this field a data.Runtime type.
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Use the new readJSON() helper to decode the request body into the input struct.
	// If this returns an error we send the client the error message along with a 400
	// Bad Request status code, just like before.
	err := app.readJSON(w, r, &input)
	if err != nil {
		// Use the new badRequestResponse() helper.
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". in the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")

	// note that we're using the Unique helper in the line below to check that all
	// values in the input.Genres slice are unique
	v.Check(validator.Unique(input.Genres), "genres", "must not contains duplicate values")

	// Use the Valid() method to see if any of the checks failed. If they did, the use
	// the failedValidationResponse() helper to send a response to the client, passing
	// in the v.Errors map.
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Dump the contents of the input struct in a HTTP response
	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movie/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.
	movie := data.Movies{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
