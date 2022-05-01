package middleware_test

import (
	"net/http"
	"testing"

	"github.com/hysem/mini-aspire-api/app/core/context"
	"github.com/hysem/mini-aspire-api/app/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_GrantAccess(t *testing.T) {
	customer := &model.User{
		UserID:   1,
		Name:     "name",
		Email:    "email@yopmail.com",
		Password: "password",
		Role:     model.RoleCustomer,
	}
	admin := &model.User{
		UserID:   1,
		Name:     "name",
		Email:    "email@yopmail.com",
		Password: "password",
		Role:     model.RoleAdmin,
	}

	testCases := map[string]struct {
		authUser           *model.User
		roles              []model.Role
		expectedStatusCode int
	}{
		`error case: authenticated user is customer; restrict all access`: {
			authUser:           customer,
			expectedStatusCode: http.StatusForbidden,
		},
		`error case: authenticated user is customer; restrict access to admin`: {
			authUser:           customer,
			roles:              []model.Role{model.RoleAdmin},
			expectedStatusCode: http.StatusForbidden,
		},
		`error case: authenticated user is admin; restrict access to customer`: {
			authUser:           admin,
			roles:              []model.Role{model.RoleCustomer},
			expectedStatusCode: http.StatusForbidden,
		},
		`success case: authenticated user is customer; grant access to customer`: {
			authUser:           customer,
			roles:              []model.Role{model.RoleCustomer},
			expectedStatusCode: http.StatusOK,
		},
		`success case: authenticated user is customer; grant access to customer and admin`: {
			authUser:           customer,
			roles:              []model.Role{model.RoleCustomer, model.RoleAdmin},
			expectedStatusCode: http.StatusOK,
		},
		`success case: authenticated user is admin; grant access to admin`: {
			authUser:           admin,
			roles:              []model.Role{model.RoleAdmin},
			expectedStatusCode: http.StatusOK,
		},
		`success case: authenticated user is admin; grant access to customer and admin`: {
			authUser:           admin,
			roles:              []model.Role{model.RoleCustomer, model.RoleAdmin},
			expectedStatusCode: http.StatusOK,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mw, m := newMiddleware(t)
			defer m.assertExpectations(t)

			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			handler := func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}

			res := runMiddlewareTest(t, req, handler, mw.context, modifyContext(func(c *context.Context) {
				c.AuthUser = tc.authUser
			}), mw.grantAccess(tc.roles...))

			assert.Equal(t, tc.expectedStatusCode, res.Result().StatusCode)
		})
	}

}
