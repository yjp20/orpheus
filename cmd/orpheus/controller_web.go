package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"
)

type App struct {
	Session *discordgo.Session
	Addr    string
	Origins []string
}

func (app *App) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	w.Write(b)
	return nil
}

func (app *App) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	return app.parseJSON(http.MaxBytesReader(w, r.Body, int64(maxBytes)), dst)
}

func (app *App) readParams(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	fmt.Printf("%v+\n", r.URL.Query())
	marshalled, err := json.Marshal(r.URL.Query())
	fmt.Printf("%v+\n", string(marshalled))
	if err != nil {
		return fmt.Errorf("query string contains malformed data")
	}
	return app.parseJSON(bytes.NewReader(marshalled), dst)
}

func (app *App) parseJSON(reader io.Reader, dst interface{}) error {
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError) && unmarshalTypeError.Field != "":
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		case errors.As(err, &unmarshalTypeError) && unmarshalTypeError.Field == "":
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

type errorResponse struct {
	Error interface{} `json:"error"`
}

func (a *App) writeError(w http.ResponseWriter, status int, message interface{}) {
	payload := errorResponse{Error: message}
	err := a.writeJSON(w, status, payload)
	if err != nil {
		w.WriteHeader(500)
	}
}

func (app *App) badRequestResponse(w http.ResponseWriter, err error) {
	app.writeError(w, http.StatusBadRequest, err.Error())
}

func serverAPI(session *discordgo.Session, addr, cors string) *http.Server {
	app := &App{
		Addr:    addr,
		Origins: strings.Split(cors, ","),
		Session: session,
	}

	router := httprouter.New()
	router.GET("/api/queue", app.getServer)
	router.POST("/api/queue", app.addQueue)

	srv := &http.Server{
		Addr:    addr,
		Handler: app.enableCORS(router),
	}

	return srv
}

func (a *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")

		for i := range a.Origins {
			if a.Origins[i] == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				if r.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) addQueue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		GuildID string `json:"guild_id"`
		Url     string `json:"url"`
		UserId  string `json:"user_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
	}

	server := getServer(input.GuildID)
	songs, _ := fetchSongsFromURL(input.Url, false)
	song := server.Add(songs, input.UserId, false)
	app.writeJSON(w, 200, &(song[0]))
}

func (app *App) getServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		GuildID string `json:"guild_id"`
	}

	input.GuildID = r.URL.Query().Get("guild_id")
	app.writeJSON(w, 200, getServer(input.GuildID))
}
