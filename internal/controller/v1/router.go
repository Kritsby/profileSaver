package v1

import (
	"context"
	"dev/profileSaver/internal/repository"
	"fmt"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"net/http"
)

type Handler struct {
	repo *repository.DB
}

func New(repo *repository.DB) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) InitRouter() *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
		bunrouter.Use(h.authMidleware),
	)

	swagHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)
	bswag := bunrouter.HTTPHandlerFunc(swagHandler)
	router.GET("/swagger/:*", bswag)

	router.WithGroup("/v1", func(g *bunrouter.Group) {
		g.WithGroup("/user", func(g *bunrouter.Group) {
			g.POST("", h.createUser)
			g.GET("", h.getAllUsers)
			g.GET("/:id", h.getUser)
			g.PATCH("/:id", h.updateUser)
			g.DELETE("/:id", h.deleteUser)
		})
	})

	return router
}

func (h *Handler) authMidleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		username, password, ok := req.BasicAuth()
		if !ok {
			askPassword(w)
		}

		user, err := h.repo.GetUserByName(username)
		if err != nil {
			internalError(w, err)
			return nil
		}

		userPasswrod, _ := h.repo.HashPass([]byte(password), user.Salt)

		if user.Password != fmt.Sprintf("%x", userPasswrod) && user.Username != "admin" {
			askPassword(w)
			return nil
		}

		ctx := context.WithValue(req.Context(), "is_admin", user.Admin)

		w.Header().Set("Content-Type", "application/json")
		return next(w, req.WithContext(ctx))
	}
}

func askPassword(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	w.WriteHeader(http.StatusUnauthorized)
}

func internalError(w http.ResponseWriter, err error) {
	log.Error().Msgf("unable to get user from the store error - %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("something went wrong"))
}

func (h *Handler) responseJSON(w http.ResponseWriter, req bunrouter.Request, code int, value interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if code != http.StatusOK {
		log.Warn().Msgf("route: %s, http code: %d, error: %v", req.Route(), code, value)
		return bunrouter.JSON(w, bunrouter.H{
			"route":  req.Route(),
			"params": req.Params().Map(),
			"error":  value,
		})
	}

	return bunrouter.JSON(w, bunrouter.H{
		"route":  req.Route(),
		"params": req.Params().Map(),
		"data":   value,
	})
}
