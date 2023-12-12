package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

const PORT = "8000"

type MovieResponse struct {
	Results []struct {
		Title    string `json:"title"`
		Overview string `json:"overview"`
		Date     string `json:"release_date"`
		ImgPath  string `json:"poster_path"`
	} `json:"results"`
}

func main() {

	// static dir system
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/search", searchMovie)

	log.Println("App running on ", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, nil)
}

func searchMovie(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	baseURL := "https://api.themoviedb.org/3/search/movie"
	apiKey := "f6bd918a5d7b3ca70f0370902cdb9560"
	url := fmt.Sprintf("%s?query=%s&include_adult=false&language=en-US&page=1&api_key=%s", baseURL, query, apiKey)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var movieData MovieResponse

	// Unmarshal the JSON data into the struct
	if err := json.Unmarshal(body, &movieData); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	if len(movieData.Results) == 0 && query != "" {

		tmpl := template.Must(template.ParseFiles("./templates/no-result.html"))
		err := tmpl.Execute(w, query)
		if err != nil {
			log.Println(err)
		}

	} else {

		tmpl := template.Must(template.ParseFiles("./templates/search-result.html"))
		err := tmpl.Execute(w, movieData.Results)
		if err != nil {
			log.Println(err)
		}
	}

}
