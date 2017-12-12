package mockme

import (
	"reflect"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	fakeit "github.com/the0rem/go-fakeit"
)

// MockMe takes the api struct and looks through each handler in the
// given api and fakes the response struct required for the response data.
// If persistence is enabled, response data will persist
func MockMe(response middleware.Responder, model interface{}, params interface{}) middleware.Responder {

	// Fake the model
	fakeit.FakeIt(&model)

	if params != nil {
		UpdateValues(&model, params)
	}

	t := reflect.ValueOf(response)
	var payload = t.FieldByName("Payload")
	if payload.IsValid() && payload.CanSet() {
		payload.Set(reflect.ValueOf(model))
	}

	return response
}

// UpdateValues looks for any struct fields that are in both the param and model structs and
// updates the model accordingly
func UpdateValues(model *interface{}, params interface{}) {
	paramVal := reflect.ValueOf(params)
	modelVal := reflect.Indirect(reflect.ValueOf(model))
	for j := 0; j < paramVal.NumField(); j++ {

		paramField := paramVal.Field(j).Addr()
		paramName := paramVal.Type().Field(j).Name

		modelField := modelVal.FieldByName(paramName)
		if !modelField.IsValid() ||
			!modelVal.CanAddr() ||
			!modelVal.CanSet() ||
			paramField.Type().String() != modelField.Type().String() ||
			paramField.Type().PkgPath() != modelField.Type().PkgPath() ||
			paramField.Kind().String() != modelField.Kind().String() {
			continue
		}

		modelVal.Set(reflect.ValueOf(paramVal))
	}
}

type MockFlags struct {
	MockEnabled string `short:"m" long:"mock" description:"Run as mock server"`
}

var mockFlags MockFlags

// AddMockFlag
// func configureFlags(api *operations.PrincipleAPI) {
// 	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{dbHost}
// }
//
func AddMockFlag(options []swag.CommandLineOptionsGroup) {
	mockFlag := swag.CommandLineOptionsGroup{
		ShortDescription: "Enable Mockme server",
		LongDescription:  "Enable Mockme server",
		Options:          &mockFlags,
	}

	options = append(options, mockFlag)
}

type Mocker struct {
	MockFlags MockFlags
}
