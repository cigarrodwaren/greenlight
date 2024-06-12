package data

import "time"

type Movies struct {
	ID        int64     `json:"id"`             // Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`              // Timestamp for when the movie is added to ourdatabase
	Title     string    `json:"title"`          // Movie title
	Year      int32     `json:"year,omitempty"` // Movie release year
	// Use the Runtime type instead of int32. Note that the omitempty directive wiil
	// still work on this: if the Runtime field has the underlying value 0, the it will
	// be considered empty and omitted -- and the MarshalJSON() method we just made
	// won't be called at all.
	Runtime Runtime  `json:"runtime,omitempty"` // Movie runtime (in minutes)
	Genres  []string `json:"genres,omitempty"`  // Slice of gender for the movie (romace, comedy, etc.)
	Version int32    `json:"version"`           // The version number starts at 1 and will be incremented
	// each time the movie informaction is updated
}
