package addToUser

import (
	"errors"
	resp "github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/api/response"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/lib/logger/sl"
	"github.com/DanilaNik/avito-backend-trainee-assignment-2023/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type Request struct {
	SegmentsToSave   []string `json:"SegmentsToSave" validate:"required"`
	SegmentsToDelete []string `json:"SegmentsToDelete" validate:"required"`
	UserID           int64    `json:"UserID" validate:"required"`
}

type Response struct {
	resp.Response
	UserId           int64    `json:"UserId ,omitempty"`
	AddedSegments    []string `json:"AddedSegments,omitempty"`
	NotAddedSegments []string `json:"NotAddedSegments,omitempty"`
	DeletedSegments  []string `json:"DeletedSegments ,omitempty"`
}

type UserToSegmentsAdder interface {
	AddUserToSegments(segmentsToSave []string, segmentsToDelete []string, userId int64) (*storage.UserInSegmentDTO, error)
}

func New(log *slog.Logger, userToSegmentsAdder UserToSegmentsAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.segment.addToUser.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		segmentsToSave := req.SegmentsToSave
		segmentsToDelete := req.SegmentsToDelete
		userID := req.UserID

		res, err := userToSegmentsAdder.AddUserToSegments(segmentsToSave, segmentsToDelete, userID)
		if err != nil {
			log.Error("failed to change user segments", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to change user segments"))

			return
		}
		log.Info("user segments changed", slog.Int64("UserId", res.UserID))

		responseOK(w, r, res)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, result *storage.UserInSegmentDTO) {
	render.JSON(w, r, Response{
		Response:         resp.OK(),
		UserId:           result.UserID,
		AddedSegments:    result.AddedSegments,
		NotAddedSegments: result.NotAddedSegments,
		DeletedSegments:  result.DeletedSegments,
	})
}
