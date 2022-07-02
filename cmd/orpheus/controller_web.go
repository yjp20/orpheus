package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/julienschmidt/httprouter"

	"github.com/yjp20/orpheus/pkg/music"
	"github.com/yjp20/orpheus/pkg/queue"
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
	router.GET("/api/guild/:ID/queue", app.getServer)
	router.POST("/api/guild/:ID/queue/add", app.addQueue)
	router.POST("/api/guild/:ID/queue/addlist", app.addList)
	router.POST("/api/guild/:ID/queue/queue", app.Queue)
	router.POST("/api/guild/:ID/queue/pause", app.Pause)
	router.POST("/api/guild/:ID/queue/resume", app.Resume)
	router.POST("/api/guild/:ID/queue/fastforward", app.fastForward)
	router.POST("/api/guild/:ID/queue/rewind", app.Rewind)
	router.POST("/api/guild/:ID/queue/seek", app.Seek)
	router.POST("/api/guild/:ID/queue/skip", app.Skip)
	router.POST("/api/guild/:ID/queue/goto", app.Goto)
	router.POST("/api/guild/:ID/queue/remove", app.Remove)
	router.POST("/api/guild/:ID/queue/shuffle", app.Shuffle)
	router.POST("/api/guild/:ID/queue/nowplaying", app.nowPlaying)
	router.POST("/api/guild/:ID/queue/loop", app.Loop)
	router.POST("/api/guild/:ID/queue/move", app.Move)
	router.POST("/api/guild/:ID/queue/clear", app.Clear)

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
		Url    string          `json:"url"`
		UserId string          `json:"user_id"`
		Policy queue.AddPolicy `json:"policy"`
	}
	input.Policy = queue.Smart

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	g := GetGuild(ps.ByName("ID"))
	songs, err := music.FetchFromURL(input.Url, false)
	if err != nil {
		app.writeError(w, http.StatusNotFound, "failed to add song")
		return
	}
	song := g.Queue.Add(songs, input.UserId, false, input.Policy)
	app.writeJSON(w, 200, &(song[0]))
}

func (app *App) addList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		Url     string          `json:"url"`
		UserID  string          `json:"user_id"`
		Policy  queue.AddPolicy `json:"policy"`
		Shuffle bool            `json:"shuffle"`
	}
	input.Policy = queue.Smart
	input.Shuffle = false

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	g := GetGuild(ps.ByName("ID"))
	songs, err := music.FetchFromURL(input.Url, true)
	if err != nil {
		app.writeError(w, http.StatusNotFound, "failed to add song")
	}
	queueItems := g.Queue.Add(songs, input.UserID, input.Shuffle, input.Policy)
	app.writeJSON(w, 200, &queueItems)
}

func (app *App) Queue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		Center int `json:"center"`
	}
	input.Center = g.Queue.Index

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if len(g.Queue.List) == 0 {
		app.writeJSON(w, 200, "Queue is empty")
		return
	}
	app.writeJSON(w, 200, PrintQueue(g, input.Center))
}

func (app *App) Pause(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	g.Player.Pause()
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Resume(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	if g.Queue.CurrentItem() == nil && len(g.Queue.List) > 0 {
		g.Queue.SkipTo(0)
	}
	g.Player.Resume()
	app.writeJSON(w, 202, "accepted")
}

func (app *App) fastForward(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		Seconds float64 `json:"seconds"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if g.Queue.CurrentItem() == nil {
		app.writeJSON(w, 202, "Not playing any song to fast-forward")
		return
	}
	g.Player.FastForward(input.Seconds)
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Rewind(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
    var input struct {
        Seconds float64 `json:"seconds"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, err)
        return
    }

	if g.Queue.CurrentItem() == nil {
		app.writeJSON(w, 202, "Not playing any song to rewind")
		return
	}
	g.Player.FastForward(-input.Seconds)
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Seek(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		Seconds float64 `json:"seconds"`
	}

	err := app.readJSON(w, r, &input)
    if err != nil {
    	app.badRequestResponse(w, err)
		return
    }	

	item := g.Queue.CurrentItem()
	if item == nil {
		app.writeJSON(w, 202, "Not playing any song to seek")
		return
	}
	if time.Duration(float64(time.Second)*input.Seconds) >= item.Song.Length || input.Seconds < 0 {
		app.writeJSON(w, 202, "Seek value out of range")
		return
	}
	g.Player.Seek(input.Seconds)
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Skip(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
    var input struct {
        Skip int `json:"skip"`
    }
	input.Skip = 1

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}
	index := (g.Queue.Index + input.Skip) % len(g.Queue.List)
    g.Queue.SkipTo(index)
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Goto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		Index int `json:"index"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if input.Index >= len(g.Queue.List) || input.Index < 0 {
		app.writeJSON(w, 202, "Index out of range")
		return
	}
    g.Queue.SkipTo(input.Index)
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Remove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		Index int `json:"index"`
	}
	input.Index = g.Queue.Index

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	queueItem, err := g.Queue.Remove(input.Index)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}
	app.writeJSON(w, 202, &queueItem)
}

func (app *App) Shuffle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	g.Queue.Shuffle()
	app.writeJSON(w, 202, "Shuffled queue")
}

func (app *App) nowPlaying(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	app.writeJSON(w, 202, formatCurrentSong("Currently Playing: ", g))
}

func (app *App) Loop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		NextPolicy queue.NextPolicy `json:"next_policy"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}
	
	g.Queue.NextPolicy = input.NextPolicy
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Move(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	var input struct {
		From int `json:"from"`
		To int `json:"to"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	_, err = g.Queue.Move(input.From, input.To)
	if err != nil {
		// TODO handle error
	}
	app.writeJSON(w, 202, "accepted")
}

func (app *App) Clear(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	g := GetGuild(ps.ByName("ID"))
	g.Queue.Clear()
	app.writeJSON(w, 202, "accepted")
}

func (app *App) getServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app.writeJSON(w, 200, GetGuild(ps.ByName("ID")))
}
