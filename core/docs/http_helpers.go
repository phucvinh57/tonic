package docs

import "net/http"

func JSONSchemaRef(ref string) MediaTypeOrReference {
	return MediaTypeOrReference{
		MediaTypeObject: &MediaTypeObject{
			Schema: &SchemaOrReference{
				ReferenceObject: &ReferenceObject{
					Ref: ref,
				},
			},
		},
	}
}

func ResponseDescription(statusCode int) string {
	if text := http.StatusText(statusCode); text != "" {
		return text
	}
	return "Response"
}

func JSONResponse(statusCode int, ref string) ResponseOrReference {
	return ResponseOrReference{
		ResponseObject: &ResponseObject{
			Description: ResponseDescription(statusCode),
			Content: map[string]MediaTypeOrReference{
				"application/json": JSONSchemaRef(ref),
			},
		},
	}
}
