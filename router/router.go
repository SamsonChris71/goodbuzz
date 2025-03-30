package router

import (
	"embed"
	"goodbuzz/router/admin"
	"goodbuzz/router/healthcheck"
	"goodbuzz/router/index"
	"goodbuzz/router/login"
	"goodbuzz/router/rooms"
	"goodbuzz/router/rooms/buzz"
	"goodbuzz/router/rooms/locks"
	"goodbuzz/router/rooms/moderator"
	"goodbuzz/router/rooms/player"
	"goodbuzz/router/tournaments"
	"net/http"
)

// this bit of go magic embeds everything in the /static directory
//
//go:embed all:static
var content embed.FS

func SetupRouter(mux *http.ServeMux) {
	mux.Handle("/static/", http.FileServer(http.FS(content)))
	mux.HandleFunc("GET /{$}", index.Get)
	mux.HandleFunc("GET /login", login.Get)
	mux.HandleFunc("POST /login", login.Post)
	mux.HandleFunc("DELETE /login", login.Delete)
	mux.HandleFunc("POST /login/player", login.PostPlayer)

	mux.HandleFunc("POST /tournaments", tournaments.Post)
	mux.Handle("GET /tournaments/{id}", tournaments.MiddlewareAny(tournaments.Get))
	mux.Handle("POST /tournaments/{id}", tournaments.MiddlewareAdmin(tournaments.PostRoom))
	mux.Handle("PUT /tournaments/{id}", tournaments.MiddlewareAdmin(tournaments.Put))
	mux.Handle("DELETE /tournaments/{id}", tournaments.MiddlewareAdmin(tournaments.Delete))

	mux.Handle("PUT /rooms/{id}", rooms.MiddlewareMod(rooms.Put))
	mux.Handle("PUT /rooms/{id}/lock", rooms.MiddlewareMod(rooms.Lock))
	mux.Handle("DELETE /rooms/{id}/lock", rooms.MiddlewareMod(rooms.Unlock))
	mux.Handle("PUT /rooms/{id}/description", rooms.MiddlewareMod(rooms.Description))
	mux.Handle("DELETE /rooms/{id}", rooms.MiddlewareAdmin(rooms.Delete))
	mux.Handle("GET /rooms/{id}/edit", rooms.MiddlewareMod(rooms.Get))

	mux.Handle("GET /rooms/{id}/player", rooms.MiddlewareAny(player.Get))
	mux.Handle("GET /rooms/{id}/player/live", rooms.MiddlewareAny(player.Live))
	mux.Handle("PUT /rooms/{id}/player", rooms.MiddlewareAny(player.Put))
	mux.Handle("PUT /rooms/{id}/player/{token}", rooms.MiddlewareAny(player.PutPlayer))
	mux.Handle("PUT /rooms/{id}/buzz", rooms.MiddlewareAny(buzz.Put))
	mux.Handle("DELETE /rooms/{id}/buzz", rooms.MiddlewareAny(buzz.Delete))

	mux.Handle("GET /rooms/{id}/moderator", rooms.MiddlewareMod(moderator.Get))
	mux.Handle("GET /rooms/{id}/moderator/live", rooms.MiddlewareMod(moderator.Live))

	mux.Handle("DELETE /rooms/{id}/players/{userToken}", rooms.MiddlewareMod(rooms.KickPlayer))
	mux.Handle("DELETE /rooms/{id}/locks/{userToken}", rooms.MiddlewareMod(locks.Delete))

	mux.Handle("GET /admin", admin.Middleware(admin.Get))
	mux.Handle("PUT /admin", admin.Middleware(admin.Put))

	mux.HandleFunc("GET /healthcheck", healthcheck.Get)
}
