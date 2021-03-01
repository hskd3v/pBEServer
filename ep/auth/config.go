package auth

import "github.com/harriklein/pBE/pBEServer/utils"

const (
	envJWTSecretKey       = "BE_JWT_SECRET_KEY"
	envJWTIssuer          = "BE_JWT_ISSUER"
	envJWTExpirationHours = "BE_JWT_EXPIRATION_HOURS"

	defJWTSecretKey       = "MySecretKey2+.=;PUH9v)h=ezgQ2f5e]YhhSv(A;uu8XY@"
	defJWTIssuer          = "API"
	defJWTExpirationHours = 24
)

var (
	cfgJWTSecretKey       = utils.EnvStr(envJWTSecretKey, defJWTSecretKey)
	cfgJWTIssuer          = utils.EnvStr(envJWTIssuer, defJWTIssuer)
	cfgJWTExpirationHours = utils.EnvInt64(envJWTExpirationHours, defJWTExpirationHours)
)
