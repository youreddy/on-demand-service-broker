package credhubbroker

type CredentialStore interface {
	Set(key string, value map[string]interface{}) error
}
