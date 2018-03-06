package ripple

import (
	"errors"
	"reflect"

	"gopkg.in/labstack/echo.v3"
)

// fieldTagKey is the field tag key for ripple
const fieldTagKey = "ripple"

// Controller is the interface for a Controller to be applied to an echo Group
type Controller interface {
	// Path is the namespace ripple will create the Group at, eg /posts
	Path() string
}

// Group applies the Controller to the echo via a new Group using the
// Controller's ripple tags as a manifest to properly associate methods/path and
// handler.
func Group(c Controller, echoMux *echo.Echo) *echo.Group {
	cvof, ctyp, err := reflectCtrl(c)
	if err != nil {
		panic(err)
	}

	grp := echoMux.Group(c.Path())

	i := 0
	n := ctyp.NumField() // returns the number of controller actions declared
	for ; i < n; i++ {
		res, err := newResource(ctyp.Field(i), cvof)
		if err != nil {
			panic(err)
		}
		if res == nil {
			continue // if there is no route
		}

		res.Set(grp)
	}
	return grp
}

// reflectCtrl is passed a user-defined controller object
// it reflects through it in order to find the name of the type that was
// declared in the user's code, going through pointers if necessary
func reflectCtrl(c Controller) (reflect.Value, reflect.Type, error) {
	vof := reflect.ValueOf(c)
	typ := vof.Type()

	if typ.Kind() == reflect.Ptr {
		vof = vof.Elem()
		typ = vof.Type()
	}

	var err error
	if typ.Kind() != reflect.Struct {
		err = errNotStruct
	}

	return vof, typ, err
}

var errNotStruct = errors.New("invalid controller type: requires a struct type")
