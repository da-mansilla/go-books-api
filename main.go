package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Libro struct {
	Id     uuid.UUID `json:"id"`
	Nombre string    `json:"nombre"`
}

var libros []Libro

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "Hola mundo")
		// w.Write([]byte("Hola Agustin"))
	})
	http.HandleFunc("/libro", libroHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func libroHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nombre := r.FormValue("nombre")
		if nombre == "" {
			http.Error(w, "Se debe proporcionar un nombre de libro", http.StatusBadRequest)
		}
		id, err := uuid.NewUUID()
		if err != nil {
			http.Error(w, "Error creando ID", http.StatusInternalServerError)
		}
		libro := Libro{
			Id:     id,
			Nombre: nombre,
		}

		libros = append(libros, libro)
		fmt.Fprintf(w, "Libro %s agregado exitosamente", libro.Nombre)
	} else if r.Method == http.MethodGet {
		librosJSON, err := json.Marshal(libros)
		if err != nil {
			http.Error(w, "Error serializando objeto", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(librosJSON)
	} else if r.Method == http.MethodPut {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Id no encontrado", http.StatusBadRequest)
		}
		var nuevoNombre struct {
			Nombre string `json:"nombre"`
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&nuevoNombre); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusInternalServerError)
		}

		var nuevoLibro Libro
		for i, libro := range libros {
			if libro.Id.String() == id {
				libros[i].Nombre = nuevoNombre.Nombre
				nuevoLibro = libros[i]
			}
		}
		respuesta, err := json.Marshal(nuevoLibro)
		if err != nil {
			http.Error(w, "Error al serializar respuesta", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respuesta)
	} else if r.Method == http.MethodDelete {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Id no encontrado", http.StatusBadRequest)
		}
		var libroBorrado Libro
		for i, libro := range libros {
			if libro.Id.String() == id {
				libroBorrado = libros[i]
				libros = append(libros[:i], libros[i+1:]...)
			}
		}
		respuesta, err := json.Marshal((libroBorrado))
		if err != nil {
			http.Error(w, "Error al serializar respuesta", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(respuesta)
	}

}
