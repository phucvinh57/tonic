package adapterutil

import "reflect"

type RouteTypes struct {
	Request  reflect.Type
	Response reflect.Type
}

func NewRouteTypes[D any, R any]() RouteTypes {
	return RouteTypes{
		Request:  reflect.TypeOf(new(D)),
		Response: reflect.TypeOf(new(R)),
	}
}
