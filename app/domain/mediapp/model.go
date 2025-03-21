package mediapp

import (
	"encoding/json"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Result struct {
	*uploader.UploadResult
}

// Encode implements the encoder interface.
func (app Result) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toResult(r *uploader.UploadResult) Result {
	return Result{r}
}
