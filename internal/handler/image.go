package handler

import (
	"fmt"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gojek/darkroom/pkg/service"
	"github.com/gojek/darkroom/pkg/storage"
	"image"
	"net/http"
)

const (
	// ContentLengthHeader is the response header key used to set content length
	ContentLengthHeader = "Content-Length"
	// CacheControlHeader is the response header key used to set cache control
	CacheControlHeader = "Cache-Control"
	// VaryHeader is the response header key used to indicate the CDN that the response should depend on client's accept header
	// Ref: https://tools.ietf.org/html/rfc7231#section-7.1.4
	VaryHeader = "Vary"
	// StorageGetErrorKey is the key used while pushing metrics update to statsd
	StorageGetErrorKey = "storage_get_error"
	// ProcessorErrorKey is the key used while pushing metrics update to statsd
	ProcessorErrorKey = "processor_error"
)

// ImageHandler is responsible for fetching the path from the storage backend and processing it if required
func ImageHandler(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.SugaredWithRequest(r)
		res := deps.Storage.GetPartially(r.Context(), r.URL.Path, &storage.GetPartiallyRequestOptions{})
		if res.Error() != nil {
			l.Errorf("error from Storage.Get: %s", res.Error())
			deps.MetricService.CountImageHandlerErrors(StorageGetErrorKey)
			w.WriteHeader(res.Status())
			return
		}
		var data []byte
		data = res.Data()

		params := make(map[string]string)
		values := r.URL.Query()
		if len(values) > 0 || deps.Manipulator.HasDefaultParams() {
			var err error
			for v := range values {
				if len(values.Get(v)) != 0 {
					params[v] = values.Get(v)
				}
			}
			data, err = deps.Manipulator.Process(service.NewSpecBuilder().WithImageData(data).WithParams(params).Build())

			if err == image.ErrFormat {
				deps.MetricService.CountImageHandlerErrors(ProcessorErrorKey)
				writeResponse(w, res.Data())
				return
			}

			if err != nil {
				l.Errorf("error from Manipulator.Process: %s", err)
				deps.MetricService.CountImageHandlerErrors(ProcessorErrorKey)
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		writeResponse(w, data)
	}
}

func writeResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set(CacheControlHeader, fmt.Sprintf("public,max-age=%d", config.CacheTime()))
	// Ref to Google CDN we support: https://cloud.google.com/cdn/docs/caching#cacheability
	w.Header().Set(VaryHeader, "Accept")

	cl, _ := w.Write(data)
	w.Header().Set(ContentLengthHeader, fmt.Sprintf("%d", cl))
}
