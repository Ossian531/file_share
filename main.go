package main

import (
	"fmt"
	"net/http"
	"file_share/api"
	"file_share/app_config"
)



func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Only serve index.html for exact "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./static/index.html")
}

func main(){

	app_config.LoadEnv()

	root := http.NewServeMux()

	api := api.Mux()

	// Mount API under /api
	root.Handle("/api/", http.StripPrefix("/api", api))

	// ---- Static files ----
	fs := http.FileServer(http.Dir("./static"))
	root.Handle("/static/", http.StripPrefix("/static/", fs))

	// ---- Index ----
	root.HandleFunc("/", indexHandler)

	port := fmt.Sprintf(":%d", app_config.AppConfig.Port)

	fmt.Printf("Listening on %s\n", port)
	http.ListenAndServe(port, root)
}

