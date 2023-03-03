package main

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/martian/v3/log"
	"net/http"
	"server/database"
	"server/models"

	//"server/models"
	"server/web"
)

func main() {
	db := database.Connection()
	database.InitDatabase()
	models.NewUser(db).FindAll()
	//models.NewUser(db).CreateTable()
	log.Infof("abcd")
	r := chi.NewRouter()

	//set cors handler for all routes
	var cors = cors.New(cors.Options{
		//AllowedOrigins: []string{"https://foo.com"}, //use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		//AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With", "access-token", "accept-version", "Session", "Traceparent", "Tracecontext"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, //maximum value not ignored by any of major browsers
	})
	//set cors handler for all routes
	r.Use(cors.Handler)

	r.Use(middleware.Logger)

	web.Web(r)

	err := http.ListenAndServe(":3030", r)
	fmt.Println(err)
}
