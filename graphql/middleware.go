package graphql

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"slices"
	"time"

	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UserAuth struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  int32  `json:"role"`
}

type TokenClaims struct {
	User UserAuth `json:"user"`
	jwt.StandardClaims
}

type ctxKey string

var (
	responseWriterKey ctxKey = "response_writer"
	userCtxKey        ctxKey = "X-Request-User"
	ErrUnauthHeader          = errors.New("the user is unauthorization")
)

type authMiddlewre struct {
	publicKey *rsa.PublicKey
}

func NewAuthMiddleware(path string) *authMiddlewre {
	pub, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Sprintln("failed to read public pem file", err)
		return nil
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		fmt.Sprintln("failed to parse public pem to rsa", err)
		return nil
	}

	return &authMiddlewre{publicKey: publicKey}
}

func ResponseWriterGetTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey, w)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		ctx = context.WithValue(ctx, "token", token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *authMiddlewre) HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role []RoleType) (interface{}, error) {
	tokenIn := ctx.Value("token")
	if tokenIn == nil {
		return nil, &gqlerror.Error{
			Message: "Missing authorization token",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		}
	}

	var tokenStr string
	tokenStr = tokenIn.(string)
	if tokenStr == "" {
		return nil, &gqlerror.Error{
			Message: "Missing authorization token",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		}
	}

	claims, err := m.ValidateToken(tokenStr)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Invalid or expired token",
			Extensions: map[string]interface{}{
				"code": "UNAUTHORIZED",
			},
		}
	}

	roleUser := MapIntToRole(claims.User.Role)
	if have := slices.Contains(role, roleUser); !have {
		return nil, &gqlerror.Error{
			Message: "Invalid role for the resource",
			Extensions: map[string]interface{}{
				"code": "FORBIDDEN",
			},
		}
	}

	ctx = context.WithValue(ctx, userCtxKey, claims.User)
	return next(ctx)
}

func (m *authMiddlewre) ValidateToken(tokenStr string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	// extract the claims and verify the signature
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return m.publicKey, nil
	})
	if err != nil {
		fmt.Println("failed to parse with claims token", err)
		return nil, err
	}

	if !token.Valid {
		fmt.Println("the token is invalid in middleware")
		return nil, ErrUnauthHeader
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		fmt.Println("the token's valid but failed to parse claims")
		return nil, ErrUnauthHeader
	}

	return claims, nil
}

func GetUserContext(ctx context.Context) (UserAuth, error) {
	u, ok := ctx.Value(userCtxKey).(UserAuth)
	if !ok {
		return u, &gqlerror.Error{Message: "user not found in context"}
	}

	return u, nil
}

func MetadataInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	user, ok := ctx.Value(userCtxKey).(*UserAuth)
	if ok {
		md := metadata.New(map[string]string{
			"id":    user.ID,
			"email": user.Email,
			"role":  fmt.Sprintf("%d", user.Role),
		})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
