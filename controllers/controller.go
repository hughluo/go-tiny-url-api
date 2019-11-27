package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/hughluo/go-tiny-url-api/models"
	"github.com/julienschmidt/httprouter"
)

var REDIS_CLIENT *redis.Client

func SetRedisClient(client *redis.Client) {
	REDIS_CLIENT = client
}

func CreateTinyURL(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	if val, ok := r.Form["longURL"]; !ok || len(val) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request, should contain key-value pair in form such as (longURL: <URL>)")
		return
	}
	longURL := r.Form["longURL"][0]
	ok, message, tinyURL := models.CreateTinyURL(REDIS_CLIENT, longURL, time.Hour)
	if !ok {
		log.Printf("create tiny url failed! %s", message)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", message)
	}

	fmt.Fprintf(w, "tinyurl created %s!", tinyURL)
}

func RetrieveLongURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ok, message, longURL := models.RetrieveLongURL(REDIS_CLIENT, ps.ByName("tinyurl"))
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "retrieval failed, message: %s", message)
	} else {
		fmt.Fprintf(w, "longURL success retrieved %s", longURL)
	}
}

func UpdateTinyURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "update, %s!\n", ps.ByName("tinyurl"))
	w.WriteHeader(http.StatusNotImplemented)
}

func DeleteTinyURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "delete, %s!\n", ps.ByName("tinyurl"))
	w.WriteHeader(http.StatusNotImplemented)
}
