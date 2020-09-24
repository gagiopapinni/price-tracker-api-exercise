package helper 


import (
	"fmt"
	"net/url"
	"math/rand"
	"log"
	"time" 
	"errors" 
	"os"
	"encoding/json"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
 )



type Config struct {
	PORT int
	DOMAIN string
	REFRESH_SECS int
	MONGOHQ_URL string
	SENDER_GMAIL_ADDR string
	SENDER_GMAIL_PASS string
}
func (c *Config) Load() error{
	file, _ := os.Open("config.json")
	defer file.Close()
	err := json.NewDecoder(file).Decode(c)
	if err != nil {
		errors.New("config reading error")
	}

        return nil
}


func ConnectDb(MONGOHQ_URL string) *mongo.Collection {
  
    client, err := mongo.NewClient(options.Client().ApplyURI(MONGOHQ_URL))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
            log.Fatal(err)
    }

    fmt.Println("Connected to db")

    collection := client.Database("price").Collection("ads")

    return collection
}


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
func RandSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func IsValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

