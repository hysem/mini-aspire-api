package featuretest

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type UserContext struct {
	jwtTokens map[string]string
}

func NewUserContext() *UserContext {
	return &UserContext{
		jwtTokens: map[string]string{},
	}
}

func (c *UserContext) DoLogin(who string, id int) error {
	email := c.getEmail(who, id)
	resp, err := httpClient.Post(endpointToken).JSON(map[string]interface{}{
		"email":    email,
		"password": password,
	}).Send()

	if err != nil {
		return errors.Wrap(err, "failed to login user")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("login failed")
	}

	c.jwtTokens[email] = gjson.GetBytes(resp.Bytes(), "data.token").String()
	return nil
}

func (c *UserContext) getEmail(who string, id int) string {
	if who == "customer" {
		who = "cstmr"
	}
	return fmt.Sprintf("%s%d@yopmail.com", who, id)
}

func (c *UserContext) getAuthHeader(who string, id int) string {
	return c.jwtTokens[c.getEmail(who, id)]
}
