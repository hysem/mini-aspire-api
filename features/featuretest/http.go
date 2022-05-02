package featuretest

import (
	"fmt"
	"os"

	"gopkg.in/h2non/baloo.v3"
)

var (
	url        = fmt.Sprintf("http://localhost:%s", os.Getenv("MINI_ASPIRE_API_APP_PORT"))
	httpClient = baloo.New(url)
)
