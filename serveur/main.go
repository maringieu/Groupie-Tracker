package main

import (
	"io/ioutil"
	"log"
	"errors"
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	"html/template"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"
const relationsURL string = "https://groupietrackers.herokuapp.com/api/relation"

type DataBase struct {
	Name string `json:"name"`
	Image string `json:"image"`
	ID int `json:"id"`
	Members string `json:"members"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`


}

type RelationBase struct {
	Index []IndexS `json:"index"`
}

type IndexS struct {
	ID int `json:"id"`
	DateLocation DateLocationS
}

type DateLocationS struct {
	DatesLocations []string `json:"dateslocations"`
}

var Artists []DataBase
var data DataBase
var DataToPrint []DataBase
var Relations []RelationBase

func GetData() {
	ArtistData()
	
	var template DataBase
	// var locate MyLocation
	for i := range Artists {
		template.ID = i + 1
		template.Image = Artists[i].Image
		template.Name = Artists[i].Name
		template.Members = Artists[i].Members
		template.CreationDate = Artists[i].CreationDate
		template.FirstAlbum = Artists[i].FirstAlbum

		DataToPrint = append(DataToPrint, template)
	}
	return
}

func ArtistData() error {

	resp, err := http.Get(baseURL + "/artists")
	if err != nil {
		return errors.New("Error by get")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Error by ReadAll")
	}
	json.Unmarshal(bytes, &Artists)
	return nil
}



func formHandler(w http.ResponseWriter, r *http.Request) {

	data := DataToPrint

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
	}


	t, _ := template.ParseFiles("./serveur/index.html")
	t.Execute(w, data)

 
}

func artist(w http.ResponseWriter, r *http.Request){

	id,_:= strconv.Atoi(r.URL.Path[8:])
	data := GetDataByID(id-1)

	t,_:= template.ParseFiles("./serveur/artist.html")
	t.Execute(w,data)
	fmt.Print(id)

}

func GetDataByID(id int)DataBase{

	var data DataBase
	for i:= range Artists {
		if i == id {
			data.ID = Artists[i].ID
			data.Image = Artists[i].Image
			data.Name = Artists[i].Name
			data.Members = Artists[i].Members
			data.CreationDate = Artists[i].CreationDate
			data.FirstAlbum = Artists[i].FirstAlbum	
		}

	}
	return data
}

func main() {
	GetData()
	fileServer := http.FileServer(http.Dir(""))
	http.Handle("/", fileServer)
	http.HandleFunc("/test", formHandler)
	http.HandleFunc("/artist/", artist )

	fmt.Printf("Starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

