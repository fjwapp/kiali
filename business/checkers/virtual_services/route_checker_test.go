package virtual_services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/tests/data"
)

func TestServiceWellVirtualServiceValidation(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	validations, valid := RouteChecker{fakeIstioObjects()}.Check()

	// Well configured object
	assert.True(valid)
	assert.Empty(validations)
}

func TestServiceMultipleChecks(t *testing.T) {
	assert := assert.New(t)

	validations, valid := RouteChecker{fakeMultipleChecks()}.Check()

	// wrong weight'ed route rule
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 2)
	assert.Equal(validations[0].Message, "Weight should be between 0 and 100")
	assert.Equal(validations[0].Severity, "error")
	assert.Equal(validations[0].Path, "spec/http[0]/route[1]/weight/145")

	assert.Equal(validations[1].Message, "Weight sum should be 100")
	assert.Equal(validations[1].Severity, "error")
	assert.Equal(validations[1].Path, "spec/http[0]/route")

}

func TestServiceOver100VirtualService(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	validations, valid := RouteChecker{fakeOver100VirtualService()}.Check()

	// wrong weight'ed route rule
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 1)
	assert.Equal(validations[0].Message, "Weight sum should be 100")
	assert.Equal(validations[0].Severity, "error")
	assert.Equal(validations[0].Path, "spec/http[0]/route")
}

func TestServiceUnder100VirtualService(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	validations, valid := RouteChecker{fakeUnder100VirtualService()}.Check()

	// wrong weight'ed route rule
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 1)
	assert.Equal(validations[0].Message, "Weight sum should be 100")
	assert.Equal(validations[0].Severity, "error")
	assert.Equal(validations[0].Path, "spec/http[0]/route")
}

func TestOneRouteWithoutWeight(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	validations, valid := RouteChecker{fakeOneRouteWithoutWeight()}.Check()

	// wrong weight'ed route rule
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 2)

	assert.Equal(validations[0].Message, "Weight sum should be 100")
	assert.Equal(validations[0].Severity, "error")
	assert.Equal(validations[0].Path, "spec/http[0]/route")

	assert.Equal(validations[1].Message, "All routes should have weight")
	assert.Equal(validations[1].Severity, "error")
	assert.Equal(validations[1].Path, "spec/http[0]/route")
}

func TestSecondHTTPRouteHasNoWeight(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	validations, valid := RouteChecker{fake2HTTPRoutes()}.Check()

	// wrong weight'ed route rule
	assert.False(valid)
	assert.NotEmpty(validations)
	assert.Len(validations, 2)

	assert.Equal(validations[0].Message, "Weight sum should be 100")
	assert.Equal(validations[0].Severity, "error")
	assert.Equal(validations[0].Path, "spec/http[0]/route")

	assert.Equal(validations[1].Message, "All routes should have weight")
	assert.Equal(validations[1].Severity, "error")
	assert.Equal(validations[1].Path, "spec/http[0]/route")
}

func fakeIstioObjects() kubernetes.IstioObject {
	validVirtualService := data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 55),
		data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 45),
			data.CreateEmptyVirtualService("reviews-well", "test", []string{"reviews"}),
		),
	)

	return validVirtualService
}

func fakeUnder100VirtualService() kubernetes.IstioObject {
	virtualService := data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 45),
		data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 45),
			data.CreateEmptyVirtualService("reviews-100-minus", "test", []string{"reviews"}),
		),
	)

	return virtualService
}

func fakeOver100VirtualService() kubernetes.IstioObject {
	virtualService := data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 55),
		data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 55),
			data.CreateEmptyVirtualService("reviews-100-plus", "test", []string{"reviews"}),
		),
	)

	return virtualService
}

func fakeMultipleChecks() kubernetes.IstioObject {
	virtualService := data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 145),
		data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 55),
			data.CreateEmptyVirtualService("reviews-multiple", "test", []string{"reviews"}),
		),
	)

	return virtualService
}

func fakeOneRouteWithoutWeight() kubernetes.IstioObject {
	validVirtualService := data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", -1),
		data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 55),
			data.CreateEmptyVirtualService("reviews-well", "test", []string{"reviews"}),
		),
	)

	return validVirtualService
}

func fake2HTTPRoutes() kubernetes.IstioObject {
	validVirtualService := data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", -1),
		data.AddRoutesToVirtualService("http", data.CreateRoute("reviews", "v1", 55),
			data.CreateEmptyVirtualService("reviews-well", "test", []string{"reviews"}),
		),
	)

	if routeTypeInterface, found := validVirtualService.GetSpec()["http"]; found {
		if routeTypeCasted, ok := routeTypeInterface.([]interface{}); ok {
			duplicateRoute := data.CreateRoute("reviews", "v1", -1)
			routeTypeCasted = append(routeTypeCasted, duplicateRoute)
			validVirtualService.GetSpec()["http"] = routeTypeCasted
		}
	}

	return validVirtualService
}
