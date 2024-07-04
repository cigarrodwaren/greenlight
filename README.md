gh repo create greenlight --public --source=. --remote=upstream --push


http://localhost:4000/v1/healthcheck


BODY='{"title":"Moana","year":2016,"runtime":"107 mins","genres":["animation","adventure"]}'

curl -i -d "$BODY" localhost:4000/v1/movies

echo $GREENLIGHT_DB_DSN
postgres://greenlight:pass@localhost/greenlight?sslmode=disable

brew install golang-migrate
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints

migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up

CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);

DROP TABLE IF EXISTS movies;

ALTER TABLE movies ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now())); 
ALTER TABLE movies ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);

ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_runtime_check; 
ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_year_check; 
ALTER TABLE movies DROP CONSTRAINT IF EXISTS genres_length_check;


curl -i localhost:4000/v1/movies/5


# Update movie
BODY='{"title":"Black Panther","year":2018,"runtime":"134 mins","genres":["sci-fi","action","adventure"]}'
curl -X PUT -d "$BODY" localhost:4000/v1/movies/2

# Get movie
curl -i localhost:4000/v1/movies/2


# Get list movies with filter, page size, sort type
curl "localhost:4000/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=year"

curl "localhost:4000/v1/movies?title=black+panther"

curl "localhost:4000/v1/movies?genres=adventure"

curl "localhost:4000/v1/movies?title=moana&genres=animation,adventure"

# Get list movies with full-text search Postgresql

curl "localhost:4000/v1/movies?title=panther"

curl "localhost:4000/v1/movies?title=the+club"

# Adding Indexes serach full-text
migrate create -seq -ext .sql -dir ./migrations add_movies_indexes
migrate -path ./migrations -database $GREENLIGHT_DB_DSN up

$greenlight=>\d movies
                                        Table "public.movies"
   Column   |            Type             | Collation | Nullable |              Default
------------+-----------------------------+-----------+----------+------------------------------------
 id         | bigint                      |           | not null | nextval('movies_id_seq'::regclass)
 created_at | timestamp(0) with time zone |           | not null | now()
 title      | text                        |           | not null |
 year       | integer                     |           | not null |
 runtime    | integer                     |           | not null |
 genres     | text[]                      |           | not null |
 version    | integer                     |           | not null | 1
Indexes:
    "movies_pkey" PRIMARY KEY, btree (id)
    "movies_genres_idx" gin (genres)
    "movies_title_idx" gin (to_tsvector('simple'::regconfig, title))
Check constraints:
    "genres_length_check" CHECK (array_length(genres, 1) >= 1 AND array_length(genres, 1) <= 5)
    "movies_runtime_check" CHECK (runtime >= 0)
    "movies_year_check" CHECK (year >= 1888 AND year::double precision <= date_part('year'::text, now()))


# Sorting

curl "localhost:4000/v1/movies?sort=-title"
curl "localhost:4000/v1/movies?sort=-runtime"

# Paginating Lists

/ Return the 5 records on page 1 (records 1-5 in the dataset)
/v1/movies?page=1&page_size=5

// Return the next 5 records on page 2 (records 6-10 in the dataset)
/v1/movies?page=2&page_size=5

// Return the next 5 records on page 3 (records 11-15 in the dataset)
/v1/movies?page=3&page_size=5