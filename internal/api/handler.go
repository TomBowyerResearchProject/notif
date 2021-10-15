package api

import (
	"encoding/json"
	"net/http"
	"notif/internal/db"
	"strconv"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"
	commonNotification "github.com/TomBowyerResearchProject/common/notification"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/go-chi/chi"
)

const (
	idParam       = "id"
	linkParam     = "link"
	usernameParam = "username"
)

func getNotificationList(w http.ResponseWriter, r *http.Request) {
	page := findPage(r)

	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		response.MessageResponseJSON(w, false, http.StatusUnprocessableEntity, response.Message{
			Message: "Failed to convert",
		})

		return
	}

	notifications := db.FindNotificationsByUsername(r.Context(), username, page)

	logger.Infof("Fetched notifications for %s on page %d", username, page)

	response.ResultResponseJSON(w, false, http.StatusOK, notifications)
}

func createNotification(w http.ResponseWriter, r *http.Request) {
	notification := &commonNotification.Notification{}

	if err := json.NewDecoder(r.Body).Decode(notification); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	notification.CreatedAt = time.Now()
	notification.Seen = false

	if err := db.CreateNotification(r.Context(), notification); err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof(
		"Created notification for %s TYPE - %s TITLE - %s MESSAGE %s",
		notification.Username,
		notification.Type,
		notification.Title,
		notification.Message,
	)

	response.ResultResponseJSON(w, false, http.StatusCreated, notification)
}

func updateNotificationsToSeen(w http.ResponseWriter, r *http.Request) {
	link := chi.URLParam(r, linkParam)
	username := chi.URLParam(r, usernameParam)
	db.UpdateNotificationsSeen(r.Context(), link, username)

	logger.Infof("Updated notifications to seen for %s link %s", username, link)

	response.MessageResponseJSON(w, false, http.StatusOK, response.Message{
		Message: "Complete",
	})
}

func updateNotificationToSeen(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idParam)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnprocessableEntity, response.Message{
			Message: err.Error(),
		})

		return
	}

	db.UpdateNotificationID(r.Context(), idInt)

	logger.Infof("Updated individual notification to seen %d", idInt)

	response.MessageResponseJSON(w, false, http.StatusOK, response.Message{
		Message: "Complete",
	})
}

func removeNotificationsByPostID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idParam)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.MessageResponseJSON(w, false, http.StatusUnprocessableEntity, response.Message{
			Message: err.Error(),
		})
	}

	db.DeleteNotificationByPostID(r.Context(), idInt)

	logger.Infof("Removed post notifications for id %d", idInt)

	response.MessageResponseJSON(w, false, http.StatusOK, response.Message{
		Message: "Complete",
	})
}
