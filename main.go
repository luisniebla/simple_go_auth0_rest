package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/simple_rest/middleware"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBDeck struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Title     string         `json:"title"`
	Words     []DBWord       `gorm:"many2many:deck_words;" json:"words"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (d *DBDeck) TableName() string {
	return "decks"
}

type DBWord struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Word           string         `json:"word"`
	Synonyms       string         `json:"-"`
	SynonymsOutput []string       `gorm:"-" json:"synonyms"`
	Decks          []DBDeck       `gorm:"many2many:deck_words;" json:"-"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}

func (d *DBWord) toResponse() {
	var output []string
	err := json.Unmarshal([]byte(d.Synonyms), &output)
	if err != nil {
		d.SynonymsOutput = []string{}
		return
	}
	d.SynonymsOutput = output
}

func GetDeckWithWords(db *gorm.DB, deckID string) (*DBDeck, error) {
	deckUUID, err := uuid.Parse(deckID)
	if err != nil {
		return nil, errors.New("ERROR")
	}
	var deck *DBDeck
	result := db.Preload("Words").First(&deck, deckUUID)
	return deck, result.Error
}

// New junction model
type DeckWord struct {
	DeckID    uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	WordID    uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	CreatedAt time.Time
}

func (d *DBWord) TableName() string {
	return "words"
}

type Routes struct {
	DB *gorm.DB
}

func encodeJSONError(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err,
	})
}

func (ro Routes) getWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	word := mux.Vars(r)["word"]

	var dbWord *DBWord
	err := ro.DB.Where("word = ?", word).First(&dbWord).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			encodeJSONError(w, "not found", http.StatusNotFound)
			return
		}
		encodeJSONError(w, "invalid", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(dbWord)
}

func (ro Routes) getDeck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	deckID := mux.Vars(r)["deckID"]

	dbDeck, err := GetDeckWithWords(ro.DB, deckID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			encodeJSONError(w, "not found", http.StatusNotFound)
			return
		}
		encodeJSONError(w, "invalid", http.StatusBadRequest)
		return
	}

	for i := range dbDeck.Words {
		dbDeck.Words[i].toResponse()
	}

	json.NewEncoder(w).Encode(dbDeck)
}

func (ro Routes) checkWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("Thanks!")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Access-Control-Allow-Origin, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&DBWord{})
	db.AutoMigrate(&User{})
	routes := Routes{DB: db}

	r := mux.NewRouter()

	r.Use(corsMiddleware)

	log.Printf("hello world")
	// https://github.com/gorilla/mux?tab=readme-ov-file#middleware
	var dir = "./static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/words/{wordID}", routes.getWord).Methods("GET")
	r.HandleFunc("/decks/{deckID}", routes.getDeck).Methods("GET")

	// Apply auth middleware to specific endpoints
	r.Handle("/private", middleware.EnsureValidToken()(http.HandlerFunc(routes.checkWord)))

	r.Handle("/api/login", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*middleware.CustomClaims)

			if !claims.HasScope("read:current_user") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"message":"Insufficient scope."}`))
				return
			}
			userInfo, err := userInfoFromAuth0(r.Header.Get("Authorization"))
			if err != nil {
				fmt.Print(err)
				return
			}
			userInfoJson, err := json.Marshal(userInfo)
			if err != nil {
				fmt.Print(err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(userInfoJson)

			db.Create(&User{
				Name: userInfo["name"].(string),
			})

			w.WriteHeader(http.StatusOK)
		}),
	))

	r.HandleFunc("/decks/{deckID}", routes.getDeck).Methods("GET")
	fmt.Println("Server is starting at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
