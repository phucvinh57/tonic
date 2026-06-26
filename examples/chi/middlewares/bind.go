package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type dataKey struct{}

var validate = validator.New(validator.WithRequiredStructEnabled())

func Bind[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req T
		if allowsBody(r.Method) && r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				returnJSONError(w, http.StatusBadRequest, "invalid json", err)
				return
			}
		}
		if err := bindTaggedValues(&req, "path", func(name string) string {
			return chi.URLParam(r, name)
		}); err != nil {
			returnJSONError(w, http.StatusBadRequest, "invalid path", err)
			return
		}
		if err := bindTaggedValues(&req, "query", func(name string) string {
			return r.URL.Query().Get(name)
		}); err != nil {
			returnJSONError(w, http.StatusBadRequest, "invalid query", err)
			return
		}
		if err := bindTaggedValues(&req, "header", func(name string) string {
			return r.Header.Get(name)
		}); err != nil {
			returnJSONError(w, http.StatusBadRequest, "invalid headers", err)
			return
		}
		if err := validate.Struct(&req); err != nil {
			returnJSONError(w, http.StatusBadRequest, "validation failed", err)
			return
		}
		ctx := context.WithValue(r.Context(), dataKey{}, req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Data[T any](r *http.Request) T {
	value := r.Context().Value(dataKey{})
	if value == nil {
		var zero T
		return zero
	}
	return value.(T)
}

func bindTaggedValues(dst any, tag string, source func(string) string) error {
	value := reflect.ValueOf(dst)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return nil
	}
	return bindStruct(value, tag, source)
}

func bindStruct(value reflect.Value, tag string, source func(string) string) error {
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := valueType.Field(i)
		if fieldType.Anonymous {
			if err := bindStruct(indirect(field), tag, source); err != nil {
				return err
			}
			continue
		}
		name := strings.Split(fieldType.Tag.Get(tag), ",")[0]
		if name == "" || name == "-" {
			continue
		}
		raw := source(name)
		if raw == "" || !field.CanSet() {
			continue
		}
		if err := setField(field, raw); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func indirect(value reflect.Value) reflect.Value {
	if value.Kind() != reflect.Ptr {
		return value
	}
	if value.IsNil() {
		value.Set(reflect.New(value.Type().Elem()))
	}
	return value.Elem()
}

func setField(field reflect.Value, raw string) error {
	field = indirect(field)
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Bool:
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		field.SetBool(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(raw, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(raw, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(value)
	default:
		return fmt.Errorf("unsupported kind %s", field.Kind())
	}
	return nil
}

func allowsBody(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch
}

func returnJSONError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":  message,
		"detail": err.Error(),
	})
}
