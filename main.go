package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	redisClient *redis.Client
)

func init() {
	// Charger les variables d'environnement à partir du fichier .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, using default Redis address...")
	}

	// Récupérer l'adresse de Redis à partir de la variable d'environnement, en utilisant une valeur par défaut si elle n'est pas définie
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// Initialiser le client Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // Pas de mot de passe par défaut
		DB:       0,  // Utiliser la base de données par défaut
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	hello_msg := os.Getenv("HELLO_MSG")
	if hello_msg == "" {
		hello_msg = fmt.Sprintf("Server listening on port %s ...", port)
	}
	http.HandleFunc("/message", handleMessage)
	fmt.Println(hello_msg)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getMessage(w, r)
	case "POST":
		postMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	// Récupérer le message stocké dans Redis
	val, err := redisClient.Get(r.Context(), "message").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message := map[string]string{"message": val}
	json.NewEncoder(w).Encode(message)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, ok := data["message"]
	if !ok {
		http.Error(w, "Message field is required", http.StatusBadRequest)
		return
	}

	// Stocker le message dans Redis
	err = redisClient.Set(r.Context(), "message", message, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "Message stored successfully"}
	json.NewEncoder(w).Encode(response)
}