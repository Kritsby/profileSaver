package v1

import (
	"dev/profileSaver/internal/model"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bunrouter"
	"net/http"
	"strings"
)

// createUser
// @Summary Create new user
// @Tags User
// @Description Create new user
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Param input body controller.UserRequest true "user"
// @Success 200
// @Failure 500
// @Router /v1/user [POST]
func (h *Handler) createUser(w http.ResponseWriter, req bunrouter.Request) error {
	body := req.Body
	defer body.Close()

	var newUser model.User
	if err := json.NewDecoder(body).Decode(&newUser); err != nil {
		log.Error().Err(err)
		return h.responseJSON(w, req, http.StatusBadRequest, err)
	}

	err := validate(newUser)
	if err != nil {
		return h.responseJSON(w, req, http.StatusBadRequest, err)
	}

	err = h.repo.CreateUser(newUser)
	if err != nil {
		log.Error().Err(err)
		return h.responseJSON(w, req, http.StatusInternalServerError, err)
	}

	return h.responseJSON(w, req, http.StatusOK, "user was created")
}

// getAllUsers
// @Summary Get all users
// @Tags User
// @Description Get all users
// @Accept  json
// @Produce  json
// @Success 200 {array} controller.UserResponse
// @Failure 500
// @Router /v1/user [GET]
func (h *Handler) getAllUsers(w http.ResponseWriter, req bunrouter.Request) error {
	users := h.repo.GetAllUsers()

	return h.responseJSON(w, req, http.StatusOK, users)
}

// getUser
// @Summary Get user by id
// @Tags User
// @Description Get user by id
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {array} controller.UserResponse
// @Failure 500
// @Router /v1/user/{id} [GET]
func (h *Handler) getUser(w http.ResponseWriter, req bunrouter.Request) error {
	id := req.Params().ByName("id")

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		return h.responseJSON(w, req, http.StatusInternalServerError, err)
	}

	return h.responseJSON(w, req, http.StatusOK, user)
}

// updateUser
// @Summary Update user
// @Tags User
// @Description Update user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body controller.UserRequest false "user"
// @Success 200
// @Failure 500
// @Router /v1/user/{id} [PATCH]
func (h *Handler) updateUser(w http.ResponseWriter, req bunrouter.Request) error {
	body := req.Body
	defer body.Close()

	var newUser model.User
	if err := json.NewDecoder(body).Decode(&newUser); err != nil {
		log.Error().Err(err)
		return h.responseJSON(w, req, http.StatusBadRequest, err)
	}

	id := req.Params().ByName("id")

	newUser.ID = id

	err := validate(newUser)
	if err != nil {
		return h.responseJSON(w, req, http.StatusBadRequest, err)
	}

	err = h.repo.UpdateUser(newUser)
	if err != nil {
		return h.responseJSON(w, req, http.StatusInternalServerError, err)
	}

	return h.responseJSON(w, req, http.StatusOK, "user was updated")
}

// deleteUser
// @Summary Delete user
// @Tags User
// @Description Delete user
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Param id path string true "user id"
// @Success 200 {array} controller.UserResponse
// @Failure 500
// @Router /v1/user/{id} [DELETE]
func (h *Handler) deleteUser(w http.ResponseWriter, req bunrouter.Request) error {
	id := req.Params().ByName("id")

	err := h.repo.DeleteUser(id)
	if err != nil {
		return h.responseJSON(w, req, http.StatusInternalServerError, err)
	}

	return h.responseJSON(w, req, http.StatusOK, "user was deleted")
}

func validate(newUser model.User) error {
	var reason []string

	if newUser.Username == "" {
		reason = append(reason, "empty username")
	}

	if newUser.Password == "" {
		reason = append(reason, "empty password")
	}

	if newUser.Email == "" {
		reason = append(reason, "empty email")
	}

	if len(reason) != 0 {
		return errors.New(strings.Join(reason, ", "))
	}

	return nil
}
