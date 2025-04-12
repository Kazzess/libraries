package jwt

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	GrpcRefreshTokenName = "refresh-token"
	reasonCode           = "1"
	domain               = "2"
	authTag              = "authorization"
	refreshTag           = "refresh-token"

	roleIdCtxTag = "role_id"
	userIdCtxTag = "user_id"

	authorizationScheme = "bearer"
)

type AuthInterceptor struct {
	jwtHelper Helper
	roles     map[string][]uint64
}

func NewAuthInterceptor(jwtHelper Helper, roles map[string][]uint64) *AuthInterceptor {
	return &AuthInterceptor{jwtHelper: jwtHelper, roles: roles}
}

func (i *AuthInterceptor) AuthorizeHandler(ctx context.Context) (context.Context, error) {
	fromContext := grpc.ServerTransportStreamFromContext(ctx)
	method := fromContext.Method()

	accessibleRoles, ok := i.roles[method]
	if !ok {
		// everyone can access
		return ctx, nil
	}

	token, err := grpc_auth.AuthFromMD(ctx, authorizationScheme)
	if err != nil {
		st, errSt := createStatusDenied("no authorization metadata", reasonCode, domain)
		if errSt != nil {
			return nil, errSt
		}
		return nil, st.Err()
	}

	tokenMC, err := i.jwtHelper.ParseToken(token)
	if err != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			st, errSt := createStatusDenied("no metadata", reasonCode, domain)
			if errSt != nil {
				return nil, errSt
			}
			return nil, st.Err()
		}
		refreshT := md.Get(GrpcRefreshTokenName)
		if len(refreshT) == 0 {
			st, errSt := createStatusDenied("bad metadata. no refresh token", reasonCode, domain)
			if errSt != nil {
				return nil, errSt
			}
			return nil, st.Err()
		}

		mapClaims, err := i.jwtHelper.ParseToken(refreshT[0])
		if err != nil {
			st, errSt := createStatusDenied("bad refresh token", reasonCode, domain)
			if errSt != nil {
				return nil, errSt
			}
			return nil, st.Err()
		}
		claims := i.jwtHelper.ParseMapClaims(mapClaims)

		pair, err := i.jwtHelper.GeneratePair(claims.UserID, claims.IssuerName, claims.RoleID)
		if err != nil {
			st, errSt := createStatusDenied("generate token fail", reasonCode, domain)
			if errSt != nil {
				return nil, errSt
			}
			return nil, st.Err()
		}

		header := metadata.New(map[string]string{
			authTag:    authorizationScheme + " " + pair.AccessToken,
			refreshTag: pair.RefreshToken,
		})

		err = grpc.SendHeader(ctx, header)
		if err != nil {
			return nil, err
		}
		token = pair.AccessToken
		tokenMC, err = i.jwtHelper.ParseToken(token)
		if err != nil {
			st, errSt := createStatusDenied("bad authorization token", reasonCode, domain)
			if errSt != nil {
				return nil, errSt
			}
			return nil, st.Err()
		}
	}

	claims := i.jwtHelper.ParseMapClaims(tokenMC)

	grpc_ctxtags.Extract(ctx).Set(roleIdCtxTag, claims.RoleID)
	grpc_ctxtags.Extract(ctx).Set(userIdCtxTag, claims.UserID)

	ctx = context.WithValue(ctx, roleIdCtxTag, claims.RoleID)
	ctx = context.WithValue(ctx, userIdCtxTag, claims.UserID)
	for _, role := range accessibleRoles {
		if role == claims.RoleID {
			return ctx, nil
		}
	}

	return ctx, nil
}

func createStatusDenied(message, reason, domain string) (*status.Status, error) {
	st := status.New(codes.PermissionDenied, message)
	br := &errdetails.ErrorInfo{}
	br.Reason = reason
	br.Domain = domain
	st, err := st.WithDetails(br)
	if err != nil {
		return st, err
	}
	return st, nil
}

func GetRoleIdCtxKey() string {
	return roleIdCtxTag
}

func GetUserIdCtxKey() string {
	return userIdCtxTag
}
