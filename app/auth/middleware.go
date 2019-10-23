package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aristat/golang-example-app/app/logger"

	"github.com/aristat/golang-example-app/app/entrypoint"

	"github.com/dgrijalva/jwt-go"

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
	Log  logger.Logger
}

// keys
type keys struct {
	publicPemKey  []byte
	privatePemKey []byte
}

func SetLogger(m *Middleware, log logger.Logger) {
	log = log.WithFields(logger.Fields{"service": prefix})
	m.Log = log

	return
}

// Handler for check Bearer token
func (m Middleware) Handler(next http.Handler) http.Handler {
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
		}

		t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Claims.(*CustomClaims); ok {
				return jwt.ParseRSAPublicKeyFromPEM(m.keys.publicPemKey)
			}

			m.Log.Error("Public key not found: %s", logger.Args(errPublicNotFound.Error()))
			return nil, errPublicNotFound
		})

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			m.Log.Error("Parse error: %s", logger.Args(err.Error()))
			fmt.Fprintf(w, `{"message":%q}`, err)
			return
		}

		if t.Valid {
			if claims, ok = t.Claims.(*CustomClaims); ok {
				subject = claims.Subject
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			m.Log.Error("Validation Error: %s", logger.Args(errAuthJWT.Error()))
			fmt.Fprintf(w, `{"message":%q}`, errAuthJWT)
			return
		}

		serviceName, serviceId := m.Service(claims)
		r = r.WithContext(appContext.NewContext(r.Context(), appContext.Mapping{
			Subject:     subject,
			ServiceId:   serviceId,
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

func NewMiddleware(cfg Config) (*Middleware, func(), error) {
	rPath := strings.Trim(cfg.RelativePath, "/")
	m := &Middleware{cfg: cfg}

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

func NewTestMiddleware() (*Middleware, func(), error) {
	private := ``
	public := ``

	middleware := &Middleware{
		cfg: Config{},
		keys: keys{
			publicPemKey:  []byte(public),
			privatePemKey: []byte(private),
		},
	}

	return middleware, func() {}, nil
}
