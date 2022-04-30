package auth

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aristat/golang-example-app/app/common"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/aristat/golang-example-app/app/entrypoint"

	"github.com/golang-jwt/jwt"

	appContext "github.com/aristat/golang-example-app/app/context"
	"github.com/pkg/errors"
)

const prefix = "app.auth"
const defaultSubject = "anonymous"
const defaultServiceName = "unknown"
const defaultServiceId = 0

var errPublicNotFound = errors.New("public key not found for issuer")
var errAuthJWT = errors.New("Authentication failed, JWT invalid")

// CustomClaims
type CustomClaims struct {
	jwt.StandardClaims
}

// Config
type Config struct {
	Services     map[string]uint64
	RelativePath string
}

// Middleware
type Middleware struct {
	keys keys
	cfg  Config
	log  logger.Logger
}

// keys
type keys struct {
	publicPemKey  []byte
	privatePemKey []byte
}

// Handler for check Bearer token
func (m Middleware) JWTHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var (
			token       string
			bearerToken = strings.Split(r.Header.Get("Authorization"), " ")
			subject     = defaultSubject
			claims      *CustomClaims
			ok          bool
		)

		if len(bearerToken) > 1 {
			token = bearerToken[1]
		} else {
			common.SendGraphqlErrorf(w, http.StatusUnauthorized, "Not found authorization token")
			return
		}

		t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Claims.(*CustomClaims); ok {
				return jwt.ParseRSAPublicKeyFromPEM(m.keys.publicPemKey)
			}

			m.log.Error("Public key not found: %s", logger.Args(errPublicNotFound.Error()))
			return nil, errPublicNotFound
		})

		if err != nil {
			m.log.Error("Parse error: %s", logger.Args(err.Error()))
			common.SendGraphqlErrorf(w, http.StatusUnauthorized, err.Error())
			return
		}

		if t.Valid {
			if claims, ok = t.Claims.(*CustomClaims); ok {
				subject = claims.Subject
			}
		} else {
			m.log.Error("Validation Error: %s", logger.Args(errAuthJWT.Error()))
			common.SendGraphqlErrorf(w, http.StatusUnauthorized, errAuthJWT.Error())
			return
		}

		serviceName, serviceID := m.Service(claims)
		r = r.WithContext(appContext.NewContext(r.Context(), appContext.Mapping{
			Subject:     subject,
			ServiceId:   serviceID,
			ServiceName: serviceName,
		}))

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Service returns service data as pair of name and id
func (m Middleware) Service(claims *CustomClaims) (string, uint64) {
	if claims == nil {
		return defaultServiceName, defaultServiceId
	}

	issuer := claims.Issuer

	if id, ok := m.cfg.Services[issuer]; ok {
		return claims.Issuer, id
	}
	return defaultServiceName, defaultServiceId
}

func NewMiddleware(cfg Config, log logger.Logger) (*Middleware, func(), error) {
	rPath := strings.Trim(cfg.RelativePath, "/")
	log = log.WithFields(logger.Fields{"service": prefix})
	m := &Middleware{cfg: cfg, log: log}

	publicKey, err := ioutil.ReadFile(entrypoint.WorkDir() + "/" + rPath + "/public_key.pem")
	if err != nil {
		return nil, func() {}, err
	}

	m.keys.publicPemKey = publicKey
	m.keys.privatePemKey, err = ioutil.ReadFile(entrypoint.WorkDir() + "/" + rPath + "/private_key.pem")
	if err != nil {
		return nil, func() {}, err
	}

	return m, func() {}, nil
}

func NewTestMiddleware(log logger.Logger) (*Middleware, func(), error) {
	log = log.WithFields(logger.Fields{"service": prefix})
	private := ``
	public := ``

	middleware := &Middleware{
		cfg: Config{},
		keys: keys{
			publicPemKey:  []byte(public),
			privatePemKey: []byte(private),
		},
		log: log,
	}

	return middleware, func() {}, nil
}
