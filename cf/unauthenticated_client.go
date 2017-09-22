package cf

import (
	"github.com/pivotal-cf/on-demand-service-broker/authorizationheader"
	"fmt"
	"log"
	"errors"
)

type UnauthClient struct {
	httpJsonClient
	url string
}


func NewUnauthenticated(
	url string,
	trustedCertPEM []byte,
	disableTLSCertVerification bool) (UnauthClient, error) {
	httpClient, err := newWrappedHttpClient(
		&authorizationheader.ClientTokenAuthHeaderBuilder{},
		trustedCertPEM,
		disableTLSCertVerification,
	)
	if err != nil {
		return UnauthClient{}, err
	}
	return UnauthClient{url: url, httpJsonClient : httpClient}, nil
}


func (c UnauthClient) GetAuthURL(logger *log.Logger) (string, error) {
	var infoResponse infoResponse
	err := c.getUnauthorized(fmt.Sprintf("%s/v2/info", c.url), &infoResponse, logger)
	if err != nil {
		return "", err
	}
	if(infoResponse.AuthorizationEndpoint == "") {
		return "", errors.New("Non-valid auth url.")
	}
	return infoResponse.AuthorizationEndpoint, nil
}