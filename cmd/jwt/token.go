package jwt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cast"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

var (
	jwtFlagPrivateKey string
	jwtFlagFields     string
)

func init() {
	tokenCmd.PersistentFlags().StringVar(&jwtFlagPrivateKey, "key", "", "private key")
	tokenCmd.PersistentFlags().StringVar(&jwtFlagFields, "fields", "", "JWT fields in JSON format")
}

var tokenCmd = &cobra.Command{
	Use:           "token",
	Short:         "Generate JWT token with sign",
	Example:       `jwt --key='{path to private key}' --fields='{in json format {"key":"value"} }'`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {

		if jwtFlagPrivateKey == "" {
			fmt.Printf("Flag `key` is required\n")
			os.Exit(1)
		}

		privatePemKey, err := ioutil.ReadFile(jwtFlagPrivateKey)

		if err != nil {

			fmt.Printf("Error occurred: %v\n", err.Error())
			os.Exit(1)
		}

		fields := make(map[string]interface{})

		if jwtFlagFields != "" {

			if err = json.Unmarshal([]byte(jwtFlagFields), &fields); err != nil {

				fmt.Printf("Parse fields error %v\n", err.Error())
				os.Exit(1)
			}
		}

		jwtEncoded, err := GenerateJWT(privatePemKey, fields)

		if err != nil {

			fmt.Printf("Error occurred: %v\n", err.Error())
			os.Exit(1)
		}

		fmt.Println(jwtEncoded)
	},
}

// GenerateJWT returns token signed by private key with filled fields
func GenerateJWT(privateKey []byte, fields map[string]interface{}) (string, error) {

	type CustomClaims struct {
		UserId int64 `json:"user_id,omitempty"`
		jwt.StandardClaims
	}

	claims := &CustomClaims{}

	for k, v := range fields {
		switch k {
		case "aud":
			claims.Audience = cast.ToString(v)
		case "sub":
			claims.Subject = cast.ToString(v)
		case "iss":
			claims.Issuer = cast.ToString(v)
		case "id":
			claims.Id = cast.ToString(v)
		case "exp":
			claims.ExpiresAt = cast.ToInt64(v)
		case "user_id":
			claims.UserId = cast.ToInt64(v)
		case "nbf":
			claims.NotBefore = cast.ToInt64(v)
		case "iat":
			claims.IssuedAt = time.Now().Unix()
		}
	}

	var (
		sPrivateKey interface{}
		err         error
	)

	sPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(sPrivateKey)

	if err != nil {

		return "", err
	}
	return ss, nil
}
