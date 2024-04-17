package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /totallyRandomImage", http.RedirectHandler(
		"https://avatars.githubusercontent.com/u/26349925?s=200&v=4",
		301,
	))

	fmt.Printf("%v", http.ListenAndServe(os.Getenv("HOST"), mux))
}
