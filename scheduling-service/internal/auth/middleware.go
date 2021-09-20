package auth

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/entity"
)

// Handler returns a JWT-based authentication middleware.
func Handler(verificationKey string) routing.Handler {
	return func(c *routing.Context) error {
		header := c.Request.Header.Get("Authorization")
		// message := ""
		// parser := &jwt.Parser{
		// 	ValidMethods: []string{"RS256"},
		// }
		// block, _ := pem.Decode([]byte(PemString()))
		// key, err := jwt.ParseRSAPublicKeyFromPEM(block.Bytes)
		// if err != nil {
		// 	message = err.Error()
		// } else {
		// 	if strings.HasPrefix(header, "Bearer ") {
		// 		token, err := parser.Parse(header[7:], func(t *jwt.Token) (interface{}, error) { return key, nil })
		// 		if err == nil && token.Valid {
		// 			err = handleToken(c, token)
		// 		}
		// 		if err == nil {
		// 			return nil
		// 		}
		// 		message = err.Error()
		// 	}
		// }
		var message string
		if len(header) < 7 {
			message = "Unauthorized"
		} else {
			tokenStr := header[7:]
			// c.Set("JWT", token)
			// ctx := WithUser(
			// 	c.Request.Context(),
			// 	"aaec208d077e4f18a6cba93ed7875cc5",
			// 	"patient@gmail.com",
			// )
			// c.Request = c.Request.WithContext(ctx)
			// return nil
			parser := &jwt.Parser{}
			token, _, _ := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
			err := handleToken(c, token)
			if err != nil {
				message = err.Error()
			} else {
				return nil
			}
		}

		c.Response.Header().Set("WWW-Authenticate", `Bearer realm="`+"API"+`"`)
		if message != "" {
			return routing.NewHTTPError(http.StatusUnauthorized, message)
		}
		return routing.NewHTTPError(http.StatusUnauthorized)
	}
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func handleToken(c *routing.Context, token *jwt.Token) error {
	ctx := WithUser(
		c.Request.Context(),
		token.Claims.(jwt.MapClaims)["id"].(string),
		token.Claims.(jwt.MapClaims)["username"].(string),
		token.Claims.(jwt.MapClaims)["role"].(string),
	)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

type contextKey int

const (
	userKey contextKey = iota
)

type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
	GetRole() string
}

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, id, name, role string) context.Context {
	return context.WithValue(ctx, userKey, entity.User{ID: id, Name: name, Role: role})
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentUser(ctx context.Context) Identity {
	if user, ok := ctx.Value(userKey).(entity.User); ok {
		return user
	}
	return nil
}
