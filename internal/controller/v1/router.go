package v1

import (
	"dev/profileSaver/internal/repository"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"net/http"
)

type Handler struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) InitRouter() *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
		bunrouter.Use(h.authMiddleware),
	)

	swagHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)
	bswag := bunrouter.HTTPHandlerFunc(swagHandler)
	router.GET("/swagger/:*", bswag)

	router.WithGroup("/v1", func(g *bunrouter.Group) {
		g.WithGroup("/user", func(g *bunrouter.Group) {
			g.WithMiddleware(h.isAdminMiddleware).POST("", h.createUser)
			g.WithMiddleware(h.isAdminMiddleware).PATCH("/:id", h.updateUser)
			g.WithMiddleware(h.isAdminMiddleware).DELETE("/:id", h.deleteUser)
			g.GET("", h.getAllUsers)
			g.GET("/:id", h.getUser)
		})
	})

	return router
}

func (h *Handler) authMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		username, password, ok := req.BasicAuth()
		if !ok {
			askPassword(w)
		}

		if !h.repo.IsAuthorized(username, password) {
			askPassword(w)
			return nil
		}

		w.Header().Set("Content-Type", "application/json")
		return next(w, req)
	}
}

func (h *Handler) isAdminMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		username, _, ok := req.BasicAuth()
		if !ok {
			askPassword(w)
		}

		user, err := h.repo.GetUserByName(username)
		if err != nil {
			internalError(w, err)
			return nil
		}

		if !user.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		w.Header().Set("Content-Type", "application/json")
		return next(w, req)
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
			"error": value,
		})
	}

	return bunrouter.JSON(w, bunrouter.H{
		"data": value,
	})
}
