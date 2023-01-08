/*
Simple mock to get some data of nasa resource type,
so there's no need to use nasa api which is highly rate-limited.
*/
package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"url-collector/internal/collector/nasa"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	http.HandleFunc("/pictures", servePictureMetadata)

	log.Fatal(http.ListenAndServe(":7070", nil))
}

func servePictureMetadata(w http.ResponseWriter, r *http.Request) {
	var b []byte
	var err error

	randomMetadata := nasa.Resource{
		Title:       gofakeit.LoremIpsumSentence(4),
		Copyright:   gofakeit.Name(),
		URL:         gofakeit.URL(),
		HDURL:       gofakeit.URL(),
		Type:        gofakeit.Noun(),
		Date:        gofakeit.Date().String(),
		Explanation: gofakeit.LoremIpsumSentence(10),
	}

	if b, err = json.Marshal(randomMetadata); err != nil {
		log.Println(err)
		return
	}

	time.Sleep(time.Duration(rand.Intn(500)+300) * time.Millisecond)

	if _, err := w.Write(b); err != nil {
		log.Println(err)
		return
	}

	log.Println("Success: " + randomMetadata.URL)
}
