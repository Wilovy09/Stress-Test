package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqBody)
}

func main() {
	httpPort := 8080

	http.HandleFunc("/login", loginHandler)

	fmt.Printf("Servidor escuchando en el puerto %d...\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

