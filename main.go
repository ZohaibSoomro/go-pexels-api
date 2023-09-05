package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zohaibsoomro/gopexelsapi/config"
)

var (
	BaseApiUrl = "https://api.pexels.com/v1/search?per_page=5"
)

type Photo struct {
	Id     int `json:"id"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Src    struct {
		URL string `json:"medium"`
	} `json:"src"`
	Photographer string `json:"photographer"`
	AltText      string `json:"alt"`
}
type Photos struct {
	Pics []Photo `json:"photos"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", greet).Methods("GET")
	r.HandleFunc("/{searchQuery}", hitRequest).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}
func hitRequest(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)["searchQuery"]
	apiUrl := BaseApiUrl + "&query=" + query
	rq, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	rq.Header.Set("Authorization", config.ApiKey)
	res, err := http.DefaultClient.Do(rq)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var photos Photos
	err = json.Unmarshal(bytes, &photos)
	if err != nil {
		log.Fatal(err)
	}
	str := "<h1>Images</h1><br><br>"
	for _, photo := range photos.Pics {
		str += "<h3>" + photo.AltText + "</h3><br><img src=" + photo.Src.URL + " alt=" + photo.AltText + " >"
		// str += "<img src= " + photo.Src.URL + "  alt= " + photo.AltText + "  width= " + fmt.Sprint(photo.Width) + "  height= " + fmt.Sprint(photo.Height) + "  >"
		str += "<br><br>"
	}
	w.Header().Set("Content-type", "text/html")
	w.Write([]byte(str))
}
