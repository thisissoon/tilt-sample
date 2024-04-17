package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type AddName struct {
	Name string `json:"name"`
}

func main() {
	for _, name := range []string{"Tilt", "Dlang", "SOON_"} {
		if err := addName(context.Background(), name); err != nil {
			fmt.Printf("could not add name: %v", err)
			return
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /name", func(w http.ResponseWriter, r *http.Request) {
		name, err := getName(r.Context())
		if err != nil {
			w.WriteHeader(500)
			_, _ = io.WriteString(w, err.Error())
			return
		}

		w.WriteHeader(200)
		_, err = io.WriteString(w, name)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
		}
	})
	mux.HandleFunc("POST /name", func(w http.ResponseWriter, r *http.Request) {
		var payload AddName
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&payload); err != nil {
			w.WriteHeader(422)
			return
		}

		if err := addName(r.Context(), payload.Name); err != nil {
			w.WriteHeader(500)
			_, _ = io.WriteString(w, err.Error())
		} else {
			w.WriteHeader(200)
		}
	})

	fmt.Printf("%v", http.ListenAndServe(os.Getenv("HOST"), mux))
}

type Entity struct {
	Value string
}

func addName(ctx context.Context, name string) error {
	ds, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if err != nil {
		return fmt.Errorf("when making Datastore client: %v", err)
	}

	e := &Entity{Value: name}
	k := datastore.NameKey("Name", name, nil)
	if _, err := ds.Put(ctx, k, e); err != nil {
		return fmt.Errorf("when writing to Datastore: %v", err)
	}

	return nil
}

func getName(ctx context.Context) (string, error) {
	ds, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if err != nil {
		return "", fmt.Errorf("when making Datastore client: %v", err)
	}

	names := []string{}
	q := datastore.NewQuery("Name")
	it := ds.Run(ctx, q)

	e := &Entity{}
	result, err := it.Next(e)
	for !errors.Is(err, iterator.Done) {
		if err != nil {
			return "", fmt.Errorf("when reading next Datastore value: %v", err)
		}
		names = append(names, result.Name)
		result, err = it.Next(e)
	}

	return names[rand.Intn(len(names))], nil
}
