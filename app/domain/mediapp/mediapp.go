package mediapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/kamogelosekhukhune777/lms/app/sdk/cloudinary"
	"github.com/kamogelosekhukhune777/lms/app/sdk/errs"
	"github.com/kamogelosekhukhune777/lms/foundation/web"
)

type app struct {
	client *cloudinary.CloudinaryService
}

func newApp(client *cloudinary.CloudinaryService) *app {
	return &app{
		client: client,
	}
}

func (a *app) uploadFile(ctx context.Context, r *http.Request) web.Encoder {
	file, handler, err := r.FormFile("file")
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	defer file.Close()

	result, err := a.client.UploadMedia(file, handler.Filename)
	if err != nil {
		//http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return errs.New(errs.Internal, err)
	}

	return toResult(result)
}

func (a *app) deleteFile(ctx context.Context, r *http.Request) web.Encoder {

	values := r.URL.Query()
	id := values.Get("id")

	if id == "" {
		//http.Error(w, "Asset ID is required", http.StatusBadRequest)
		return errs.New(errs.InvalidArgument, errors.New("asset id is required"))
	}

	err := a.client.DeleteMedia(id)
	if err != nil {
		//http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return errs.New(errs.Internal, errors.New("error deleting file"))
	}
	/*
		json.NewEncoder(w).Encode(map[string]string{
			"success": "true",
			"message": "Asset deleted successfully",
		})
	*/
	return nil
}

func (a *app) bulkUpload(ctx context.Context, r *http.Request) web.Encoder {
	r.ParseMultipartForm(10 << 20) // Max 10MB
	files := r.MultipartForm.File["files"]

	var results []interface{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			//http.Error(w, "Error opening file", http.StatusInternalServerError)
			return errs.New(errs.Internal, errors.New("error opening file"))
		}
		defer file.Close()

		result, err := a.client.UploadMedia(file, fileHeader.Filename)
		if err != nil {
			//http.Error(w, "Error uploading files", http.StatusInternalServerError)
			return errs.New(errs.Internal, errors.New("erroe uploading files"))
		}

		results = append(results, result)
	}
	/*
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    results,
		})
	*/
	return nil
}
