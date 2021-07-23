package api

import (
	"Shopee_UMS/usecases"
	"Shopee_UMS/utils"
	"net/http"
)

func (s *Server) handleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.handleGetProfile()(w, r)
		case "PUT":
			s.handleChangeNickname()(w, r)
		case "POST":
			s.handleUploadProfilePicture()(w, r)
		default:
			respondHTTPErr(w, r, http.StatusMethodNotAllowed)
		}
	}
}

func (s *Server) handleGetProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value("uid").(int)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		user, err := s.u.User.GetData(uid)
		if err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		respond(w, r, http.StatusOK, user)
	}
}

func (s *Server) handleChangeNickname() http.HandlerFunc {
	type request struct {
		Nickname string `json:"nickname"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value("uid").(int)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		var req request
		if err := decodeBody(r, &req); err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}

		nickname, err := s.u.User.ChangeNickname(uid, req.Nickname)
		if err != nil {
			if ve, ok := err.(utils.ValidationError); ok {
				respondErr(w, r, http.StatusBadRequest, ve)
				return
			}
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		user, err := s.u.User.GetData(uid)
		if err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}
		user.Nickname = nickname

		respond(w, r, http.StatusOK, user)
	}
}

func (s *Server) handleUploadProfilePicture() http.HandlerFunc {
	const MaxMemoryMultipartForm = 2000000 // 2MB
	allowedFileType := []string{"image/jpeg", "image/jpg", "image/png"}

	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := r.Context().Value("uid").(int)
		if !ok {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		if err := r.ParseMultipartForm(MaxMemoryMultipartForm); err != nil {
			respondErr(w, r, http.StatusBadRequest, "could not parse request form")
			return
		}

		file, header, err := r.FormFile("profilePicture")
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, "could not parse profile picture file")
			return
		}
		defer file.Close()

		buf := make([]byte, 512)
		if _, err := file.Read(buf); err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		fileType := http.DetectContentType(buf)
		if !utils.StringInSlice(fileType, allowedFileType) {
			respondErr(w, r, http.StatusBadRequest, "only support following file types", allowedFileType)
			return
		}

		photo := usecases.Photo{Data: file, Size: header.Size}
		uri, err := s.u.User.UploadProfilePicture(uid, &photo)
		if err != nil {
			if ve, ok := err.(utils.ValidationError); ok {
				respondErr(w, r, http.StatusBadRequest, ve)
				return
			}
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}

		user, err := s.u.User.GetData(uid)
		if err != nil {
			respondHTTPErr(w, r, http.StatusInternalServerError)
			return
		}
		user.ProfilePicUri = uri

		respond(w, r, http.StatusOK, user)
	}
}
