// Package dummy Dummy API
//
// the111wewewqeqweqweqweqweqweqweqweqweq purpose of this application is to provide an application
// that is using plain go code to define an API
//
//     Schemes: http
//     Host: localhost:9090
//     BasePath: /api/v1
//     Version: 1.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package dummy

// A dummy list
// swagger:response dummyListResponse
type dummyListResponseWrapper struct {
	// All current dummy list
	// in: body
	Body []TDummy
}
