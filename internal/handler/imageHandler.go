package handler

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"strings"

	"imaging-service/internal/parser"
	"imaging-service/internal/processor"
	"imaging-service/pkg/utils"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		http.Error(w, "missing image path", http.StatusBadRequest)
		return
	}

	opts, imageURL, err := parser.ParseOptions(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse options: %v", err), http.StatusBadRequest)
		return
	}

	img, err := utils.FetchImage(imageURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to fetch image: %v", err), http.StatusBadGateway)
		return
	}

	processed, err := processor.ProcessImage(img, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to process image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, processed, &jpeg.Options{Quality: opts.Quality})
}
