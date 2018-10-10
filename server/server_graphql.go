package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

type Album struct {
	ID     string `json:"id,omitempty"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Year   string `json:"year"`
	Genre  string `json:"genre"`
	Type   string `json:"type"`
}

type Artist struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Song struct {
	ID       string `json:"id,omitempty"`
	Album    string `json:"album"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Type     string `json:"type"`
}

var albums []Album = []Album{
	Album{
		ID:     "ts-fearless",
		Artist: "1",
		Title:  "Fearless",
		Year:   "2008",
		Type:   "album",
	},
}

var artists []Artist = []Artist{
	Artist{
		ID:   "1",
		Name: "Taylor Swift",
		Type: "artist",
	},
}

var songs []Song = []Song{
	Song{
		ID:       "1",
		Album:    "ts-fearless",
		Title:    "Fearless",
		Duration: "4:01",
		Type:     "song",
	},
	Song{
		ID:       "2",
		Album:    "ts-fearless",
		Title:    "Fifteen",
		Duration: "4:54",
		Type:     "song",
	},
}

func Start() {
	songType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Song",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"album": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"duration": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	artistType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Artist",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	albumType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Album",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"artist": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.String,
			},
			"genre": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	fields := graphql.Fields{
		"song": &graphql.Field{
			Type: graphql.NewList(songType),
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Find song by title",
				},
				"album": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Find song by album",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				//id, _ := strconv.Atoi(p.Args["id"].(string))
				title, _ := p.Args["title"].(string)
				album, _ := p.Args["album"].(string)
				var songsFound []Song
				if title != "" {
					for i := range songs {
						fmt.Println("songs: " + songs[i].Title)
						if songs[i].Title == title {
							songsFound = append(songsFound, songs[i])
						}
					}
				} else if album != "" {
					for i := range songs {
						fmt.Println("album: " + songs[i].Title)
						if songs[i].Album == album {
							songsFound = append(songsFound, songs[i])
						}
					}
				} else {
					songsFound = songs
				}

				return songsFound, nil
			},
		},
		"artist": &graphql.Field{
			Type: graphql.NewList(artistType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return songs, nil
			},
		},
		"album": &graphql.Field{
			Type: graphql.NewList(albumType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return albums, nil
			},
		},
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: fields})

	schemaConfig := graphql.SchemaConfig{Query: rootQuery}
	schema, _ := graphql.NewSchema(schemaConfig)
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":12345", nil)
}
