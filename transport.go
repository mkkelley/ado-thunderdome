package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	template2 "html/template"
	"net/http"
)

func GenerationPage() http.HandlerFunc {
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
}

func BattlePage() http.HandlerFunc {
	template, err := template2.New("battle.html").ParseFiles("battle.html")
	if err != nil {
		panic(err)
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		model := BattlePageModel{
			QueryId:           request.URL.Query().Get("queryId"),
			ThunderdomeApiKey: request.URL.Query().Get("thunderdomeApiKey"),
		}
		err := template.Execute(writer, model)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			fmt.Println(err)
		}
	}
}

func BattlePageFormHandler(config *AppConfig) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte("unable to parse form"))
			return
		}
		queryId := request.Form.Get("queryId")
		apiKey := request.Form.Get("thunderdomeApiKey")

		battle, err := generateBattle(config, apiKey, queryId)
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

func RunHttpServer(config *AppConfig) {
	r := chi.NewRouter()
	r.Get("/generate", GenerationPage())
	r.Get("/battle", BattlePage())
	r.Post("/battle", BattlePageFormHandler(config))
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/generate", 302)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r)
	if err != nil {
		panic(err)
	}
}
