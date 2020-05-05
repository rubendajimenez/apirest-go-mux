package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Mensaje struct {
	Imprimir string `json:"imprimir"`
}

type Marcar struct {
	ID       string `json:"id"`
	FECHA    string `json:"fecha"`
	Fotourl  string `json:"fotourl"`
	Latitud  string `json:"latitud"`
	Longitud string `json:"longitud"`
}

var dbconexion *sql.DB
var err error

func uploadFile(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Uploading File\n")

	// 1. parse input, type multipart/form-data
	r.ParseMultipartForm(10 << 20)

	//2. retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving file from form-data")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File %+v\n", handler.Filename)
	fmt.Printf("File Size %+v\n", handler.Size)
	fmt.Printf("MIME Header : %+v\n", handler.Header)

	//3. write temporary file on our server
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)

	//4. return whether or not this has been successful
	//fmt.Fprintf(w, "Successfully Uploaded File\n")
	var men Mensaje
	men.Imprimir = "dd"
	json.NewEncoder(w).Encode(men)

}

func crearMarcacion(w http.ResponseWriter, r *http.Request) {
	resultado, err := dbconexion.Prepare("INSERT INTO marcar (fecha,fotourl,latitud,longitud) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	//fmt.Println("Esto llega " + r.Body.)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	fmt.Println("esto sale" + keyVal["fecha"] + " - " + keyVal["fotourl"] + " - " + keyVal["latitud"] + " - " + keyVal["longitud"])
	fecha := keyVal["fecha"]
	fotourl := keyVal["fotourl"]
	latitud := keyVal["latitud"]
	longitud := keyVal["longitud"]

	_, err = resultado.Exec(fecha, fotourl, latitud, longitud)
	if err != nil {
		panic(err.Error())
	}

	var men Mensaje
	men.Imprimir = "creado"
	json.NewEncoder(w).Encode(men)
}

func getMarcacion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var marcars []Marcar

	result, err := dbconexion.Query("SELECT id,fecha,fotourl,latitud,longitud from marcar")

	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var marcar_ Marcar
		err := result.Scan(&marcar_.ID, &marcar_.FECHA, &marcar_.Fotourl, &marcar_.Latitud, &marcar_.Longitud)
		if err != nil {
			panic(err.Error())
		}
		marcars = append(marcars, marcar_)
	}

	json.NewEncoder(w).Encode(marcars)
}

func setupRoutes() {

	dbconexion, err = sql.Open("mysql", "root:Asd123**@tcp(localhost:3306)/marcaciones")

	if err != nil {
		fmt.Println("error al abrir la base de datoss")
		panic(err.Error())
	}

	defer dbconexion.Close()

	router := mux.NewRouter()

	router.HandleFunc("/upload", uploadFile).Methods("POST")
	router.HandleFunc("/marcar", crearMarcacion).Methods("POST")
	router.HandleFunc("/marcar", getMarcacion).Methods("GET")

	http.ListenAndServe(":8000", router)

}

func main() {
	setupRoutes()
}
