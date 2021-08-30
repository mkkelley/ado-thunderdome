package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	template2 "html/template"
	"net/http"
	"time"
)

func generationPage() http.HandlerFunc {
	template, err := template2.New("generation.html").ParseFiles("generation.html")
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		err := template.Execute(writer, nil)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			fmt.Println(err)
		}
	}
}

type BattlePageModel struct {
	QueryId           string
	ThunderdomeApiKey string
	BattlePrefix      string
}

func battlePage() http.HandlerFunc {
	template, err := template2.New("battle.html").ParseFiles("battle.html")
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		model := BattlePageModel{
			QueryId:           request.URL.Query().Get("queryId"),
			ThunderdomeApiKey: request.URL.Query().Get("thunderdomeApiKey"),
			BattlePrefix:      request.URL.Query().Get("battlePrefix"),
		}
		err := template.Execute(writer, model)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			fmt.Println(err)
		}
	}
}

func battlePageFormHandler(config *AppConfig) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte("unable to parse form"))
			return
		}
		queryId := request.Form.Get("queryId")
		apiKey := request.Form.Get("thunderdomeApiKey")
		battlePrefix := request.Form.Get("battlePrefix")

		battle, err := generateBattle(config, apiKey, queryId, battlePrefix)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			fmt.Println(err)
			return
		}
		newUrl := getUrlForPlan(battle, config)
		html := fmt.Sprintf("<a href=\"%s\">%s</a>", newUrl, newUrl)

		writer.Write([]byte(html))
	}
}

func RunHttpServer(config *AppConfig) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/generate", generationPage())
	r.Get("/battle", battlePage())
	r.Post("/battle", battlePageFormHandler(config))
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/generate", 302)
	})

	return http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
}
