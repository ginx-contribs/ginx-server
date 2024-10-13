package token

import (
	"context"
	"errors"
	"github.com/ginx-contribs/ginx-server/pkg/utils/idx"
	"github.com/ginx-contribs/ginx-server/pkg/utils/ts"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/jwtx"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
	"time"
)

var (
	ErrTokenNeedsRefresh      = errors.New("token needs refresh")
	ErrRefreshTokenExpired    = errors.New("refresh token expired")
	ErrAccessTokenExpired     = errors.New("access token expired")
	ErrTokenNotSupportRefresh = errors.New("token not support refresh")
	ErrMisMatchTokenPair      = errors.New("mismatch token pair")
)

// Options is configuration for token resolver
type Options struct {
	// id generator
	IdGen func() string
	// token cache
	Cache         Cache
	AccessPrefix  string
	RefreshPrefix string

	// token issuer name
	Issuer string
	// access token signing key

	AccessSecret string
	AccessMethod jwt.SigningMethod
	// duration to expire access token
	AccessExpired time.Duration
	// delay duration after expiration
	AccessDelay time.Duration

	// duration to expire refresh token
	RefreshSecret string
	RefreshMethod jwt.SigningMethod
	// refresh token signing key
	RefreshExpired time.Duration
}

// Claims consisted of jwt.RegisteredClaims and custom Payload.
type Claims struct {
	// payload information
	Payload any
	// whether to need refresh
	Remember bool
	jwt.RegisteredClaims
}

// Token is Token Information, it could be from issued, parsed.
type Token struct {
	// raw token string
	Raw    string
	Claims Claims
	Token  *jwt.Token
}

// Pair is an issued token pair that consist of access-token and refresh-token
type Pair struct {
	Access  Token
	Refresh Token
}

func NewResolver(options Options) *Resolver {
	if options.IdGen == nil {
		options.IdGen = func() string {
			return idx.ULID()
		}
	}
	if options.Cache == nil {
		options.Cache = NewMemoryCache()
	}
	if options.AccessPrefix == "" {
		options.AccessPrefix = "access"
	}
	if options.RefreshPrefix == "" {
		options.RefreshPrefix = "refresh"
	}
	if options.Issuer == "" {
		options.Issuer = "ginx-server"
	}
	if options.AccessSecret == "" {
		options.AccessSecret = idx.ULID()
	}
	if options.AccessMethod == nil {
		options.AccessMethod = jwt.SigningMethodHS512
	}
	if options.AccessExpired == 0 {
		options.AccessExpired = 2 * time.Hour
	}
	if options.AccessDelay == 0 {
		options.AccessDelay = 10 * time.Minute
	}
	if options.RefreshSecret == "" {
		options.RefreshSecret = idx.ULID()
	}
	if options.RefreshMethod == nil {
		options.RefreshMethod = jwt.SigningMethodHS512
	}
	if options.RefreshExpired == 0 {
		options.RefreshExpired = 144 * time.Hour
	}
	return &Resolver{opt: options}
}

// Resolver is responsible for resolving jwt token
type Resolver struct {
	opt Options
}

// Issue return a new issued token pair with given payload, it will return refresh token if refresh is true.
func (r *Resolver) Issue(ctx context.Context, payload any, refresh bool) (Pair, error) {
	var (
		tokenCache = r.opt.Cache
		issuedAt   = ts.Now()
		// consider network latency
		latency = 5 * time.Second
		pair    Pair
	)

	// issued access token
	accessToken, err := r.createToken(r.opt.AccessSecret, refresh, payload, r.opt.AccessMethod, issuedAt, r.opt.AccessExpired)
	if err != nil {
		return Pair{}, err
	}
	accessTTL := r.opt.AccessExpired + latency
	if refresh {
		accessTTL += r.opt.AccessDelay
	}
	// store in cache
	err = tokenCache.Set(ctx, r.opt.AccessPrefix, accessToken.Claims.ID, accessToken.Claims.ID, accessTTL)
	if err != nil {
		return Pair{}, statuserr.InternalError(err)
	}
	pair.Access = accessToken

	// just return if no need to issue refresh-token
	if !refresh {
		return pair, nil
	}

	// issue refresh-token
	refreshToken, err := r.createToken(r.opt.RefreshSecret, refresh, payload, r.opt.RefreshMethod, issuedAt, r.opt.RefreshExpired)
	if err != nil {
		return Pair{}, err
	}
	// associate access-token with refresh-token by access-id
	err = tokenCache.Set(ctx, r.opt.RefreshPrefix, refreshToken.Claims.ID, accessToken.Claims.ID, r.opt.RefreshExpired)
	if err != nil {
		return Pair{}, statuserr.InternalError(err)
	}
	pair.Refresh = refreshToken

	return pair, nil
}

// Refresh refreshes the access token lifetime with the given refresh token
func (r *Resolver) Refresh(ctx context.Context, accessTokenStr, refreshTokenStr string) (Pair, error) {
	var (
		now        = ts.Now()
		pair       = Pair{}
		tokenCache = r.opt.Cache
	)
	// parse access-token
	refreshToken, err := r.VerifyRefresh(ctx, refreshTokenStr)
	if err != nil {
		return pair, err
	}
	pair.Refresh = refreshToken

	// parse refresh-token
	accessToken, err := r.VerifyAccess(ctx, accessTokenStr)
	if errors.Is(err, ErrAccessTokenExpired) {
		if !accessToken.Claims.Remember {
			return pair, ErrTokenNotSupportRefresh
		}
		// check if is over delay
		if accessToken.Claims.ExpiresAt.Add(r.opt.AccessDelay).Sub(now) < 0 {
			return pair, ErrAccessTokenExpired
		}
	} else if err != nil {
		return pair, err
	}

	// access-token might be not expired, or expired but not over delay in there.

	// check token pair if is match
	accessId, e, err := tokenCache.Get(ctx, r.opt.RefreshPrefix, refreshToken.Claims.ID)
	if !e {
		return pair, ErrRefreshTokenExpired
	} else if err != nil {
		return pair, statuserr.InternalError(err)
	}
	if accessToken.Claims.ID != accessId {
		return pair, ErrMisMatchTokenPair
	}

	// issue a new access-token
	newAccessToken, err := r.createToken(r.opt.AccessSecret, true, accessToken.Claims.Payload, r.opt.AccessMethod, now, r.opt.AccessExpired)
	if err != nil {
		return pair, err
	}
	pair.Access = newAccessToken

	ttl, e, err := tokenCache.TTL(ctx, r.opt.AccessPrefix, accessToken.Claims.ID)
	if !e {
		return pair, ErrAccessTokenExpired
	} else if err != nil {
		return pair, statuserr.InternalError(err)
	}

	// for the access-token, the max ttl is 2 * AccessExpired
	ttl += r.opt.AccessExpired / 2
	if ttl >= 2*r.opt.AccessExpired {
		ttl = 2 * r.opt.AccessExpired
	}
	err = tokenCache.Set(ctx, r.opt.AccessPrefix, newAccessToken.Claims.ID, newAccessToken.Claims.ID, ttl)
	if err != nil {
		return pair, statuserr.InternalError(err)
	}

	// update token pair association
	err = tokenCache.Set(ctx, r.opt.RefreshPrefix, refreshToken.Claims.ID, newAccessToken.Claims.ID, -1)
	if err != nil {
		return pair, statuserr.InternalError(err)
	}
	return pair, nil
}

// VerifyAccess verify the access-token if is valid.
func (r *Resolver) VerifyAccess(ctx context.Context, tokenString string) (Token, error) {
	// check if is valid
	accessToken, err := r.parseToken(tokenString, r.opt.AccessSecret, r.opt.AccessMethod)
	// if it is expired
	if errors.Is(err, jwt.ErrTokenExpired) {
		// token need to refresh
		if accessToken.Claims.Remember && accessToken.Claims.ExpiresAt.Add(r.opt.AccessDelay).Before(ts.Now()) {
			return accessToken, ErrTokenNeedsRefresh
		}
		return accessToken, ErrAccessTokenExpired
	} else if err != nil {
		return accessToken, err
	}

	// check in cache
	if _, found, err := r.opt.Cache.Get(ctx, r.opt.AccessPrefix, accessToken.Claims.ID); err != nil {
		return accessToken, err
	} else if !found {
		return accessToken, ErrAccessTokenExpired
	}
	return accessToken, nil
}

func (r *Resolver) VerifyRefresh(ctx context.Context, tokenString string) (Token, error) {
	// check if is valid
	refreshToken, err := r.parseToken(tokenString, r.opt.RefreshSecret, r.opt.RefreshMethod)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return refreshToken, ErrRefreshTokenExpired
	} else if err != nil {
		return refreshToken, err
	}
	// check in cache
	if _, found, err := r.opt.Cache.Get(ctx, r.opt.RefreshPrefix, refreshToken.Claims.ID); err != nil {
		return refreshToken, statuserr.InternalError(err)
	} else if !found {
		return refreshToken, ErrRefreshTokenExpired
	}
	return refreshToken, nil
}

// issue a new token with given args
func (r *Resolver) createToken(secret string, refresh bool, payload any, method jwt.SigningMethod, at time.Time, ttl time.Duration) (Token, error) {
	// generate unique id
	id := r.opt.IdGen()
	claims := Claims{
		Remember: refresh,
		Payload:  payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    r.opt.Issuer,
			IssuedAt:  jwt.NewNumericDate(at),
			ExpiresAt: jwt.NewNumericDate(at.Add(ttl)),
			ID:        id,
		},
	}
	// issue token
	issuedToken, err := jwtx.IssueWithClaims([]byte(secret), method, claims)
	if err != nil {
		return Token{}, err
	}
	return Token{
		Raw:    issuedToken.SignedString,
		Claims: claims,
		Token:  issuedToken.Token,
	}, nil
}

// check a token if is valid, then return token info.
func (r *Resolver) parseToken(tokenString, secret string, method jwt.SigningMethod) (Token, error) {
	var token Token
	verifiedToken, err := jwtx.VerifyWithClaims(tokenString, []byte(secret), method, &Claims{})
	if verifiedToken != nil {
		token = Token{
			Raw:   verifiedToken.SignedString,
			Token: verifiedToken.Token,
		}
		if !reflect.ValueOf(verifiedToken.Claims).IsNil() {
			token.Claims = *verifiedToken.Claims.(*Claims)
		}
	}
	return token, err
}
