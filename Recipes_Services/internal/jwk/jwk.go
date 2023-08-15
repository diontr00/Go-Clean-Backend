package jwk

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"net/url"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	RS256 = SignatureAlgorithm("RS256")
	RS384 = SignatureAlgorithm("RS384")
	RS512 = SignatureAlgorithm("RS512")

	ES256 = SignatureAlgorithm("ES256")
	ES384 = SignatureAlgorithm("ES384")
	ES512 = SignatureAlgorithm("ES512")

	PS256 = SignatureAlgorithm("PS256")
	PS384 = SignatureAlgorithm("PS384")
	PS512 = SignatureAlgorithm("PS512")

	HS256 = SignatureAlgorithm("HS256")
	HS384 = SignatureAlgorithm("HS384")
	HS512 = SignatureAlgorithm("HS512")

	EdDSA = SignatureAlgorithm("EdDSA")
)

type SignatureAlgorithm string

// hold the secret key  and hashFunc and exposes sign for signing payload  and Verify for verify the expected Signature
// this is intended to work with jwt middleware that supplies with symmetric signing scheme
type hmacKey struct {
	secretKey []byte
	hashFunc  func() hash.Hash
}

func (k *hmacKey) Sign(payload []byte) ([]byte, error) {
	h := hmac.New(k.hashFunc, k.secretKey)
	_, err := h.Write(payload)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (k *hmacKey) Verify(payload, signature []byte) error {
	expectedSignature, err := k.Sign(payload)
	if err != nil {
		return err
	}

	if !hmac.Equal(expectedSignature, signature) {
		return fmt.Errorf("Invalid signature")
	}
	return nil

}

// Either implement by provider with out without cache
//
//go:generate  mockery --name JWKSProvider
type JWKSProvider interface {
	//  Key Func return JWKS object that can be used with jwt middle to validate
	KeyFunc(SignatureAlgorithm) func(*jwt.Token) (interface{}, error)
}

// Getting JWKS from issuerURL
type Provider struct {
	IssuerURL *url.URL
	Client    *fiber.Client
}

type ProviderOptions func(*Provider)

func NewJWKProvider(issuerURL *url.URL, opts ...ProviderOptions) *Provider {
	provider := &Provider{
		IssuerURL: issuerURL,
		Client:    fiber.AcquireClient(),
	}
	for _, opt := range opts {
		opt(provider)
	}
	return provider
}

// Custom Fiber http client options
func WithCustomClient(c *fiber.Client) func(*Provider) {
	return func(p *Provider) { p.Client = c }
}

// Return the pulled jwks
func (j *Provider) KeyFunc(
	signingMethod SignatureAlgorithm,
) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {

		if err := validateSigningMethod(string(signingMethod), t.Method.Alg()); err != nil {
			return nil, err
		}

		request_agent := j.Client.Get(j.IssuerURL.String())
		defer fiber.ReleaseAgent(request_agent)
		request_agent.Request().SetTimeout(time.Second * 5)

		var jwks jose.JSONWebKeySet
		_, _, errs := request_agent.Struct(&jwks)
		if errs != nil {
			return nil, fmt.Errorf("Unable to get JWKS due to : %v", errs)
		}

		keyID := t.Header["kid"].(string)
		keys := jwks.Key(keyID)
		if len(keys) == 0 {
			return nil, fmt.Errorf("Key not found for kid: %s", keyID)
		}

		switch signingMethod {

		case RS256, RS384, RS512:
			signKey, ok := keys[0].Key.(*rsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("Key is not an RSA public key : %T", signKey)
			}
			return signKey, nil

		case ES256, ES384, ES512:
			signKey, ok := keys[0].Key.(*ecdsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("Key is not an ECDSA public key : %T", signKey)
			}

		case PS256, PS384, PS512:
			signKey, ok := keys[0].Key.(*rsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("Key is not an RSA-PSS padding schem : %T", signKey)
			}

			var hash crypto.Hash
			switch signingMethod {
			case PS256:
				hash = crypto.SHA256
			case PS384:
				hash = crypto.SHA384
			case PS512:
				hash = crypto.SHA512
			default:
				return nil, fmt.Errorf("Unsupported Signing method: %s", signingMethod)
			}

			hashFunc := rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthEqualsHash,
				Hash:       hash,
			}.Hash.HashFunc()

			return &hashFunc, nil

		case HS256, HS384, HS512:
			signKey, ok := keys[0].Key.([]byte)
			if !ok {
				return nil, fmt.Errorf("Key is not in the expected format")
			}
			var hashFunc func() hash.Hash

			switch signingMethod {
			case HS256:
				hashFunc = sha256.New
			case HS384:
				hashFunc = sha512.New384
			case HS512:
				hashFunc = sha512.New
			default:
				return nil, fmt.Errorf("Unsupported signing method: %s", signingMethod)

			}
			return &hmacKey{
				secretKey: signKey,
				hashFunc:  hashFunc,
			}, nil
		default:
			return nil, fmt.Errorf("Unsupported signing methodn : %s", signingMethod)

		}
		return nil, fmt.Errorf("KeyFunc cannot be generated due to unknown error")

	}
}

type cacheJWKS struct {
	token   *rsa.PublicKey
	expires time.Time
}

// Getting JWKS from the issuer with caching support
type ProviderWithCache struct {
	Provider  *Provider
	mu        sync.RWMutex
	cache     map[string]cacheJWKS
	CacheTime time.Duration
}

// If cache time = 0  then it will set to 1 minute
func NewJWKProviderWithCache(
	issuerURL *url.URL,
	cacheTime time.Duration,
	opts ...ProviderOptions,
) *ProviderWithCache {
	if cacheTime == 0 {
		cacheTime = 1 * time.Minute
	}

	return &ProviderWithCache{
		Provider:  NewJWKProvider(issuerURL, opts...),
		cache:     map[string]cacheJWKS{},
		CacheTime: cacheTime,
	}

}

func (c *ProviderWithCache) KeyFunc(
	signingMethod SignatureAlgorithm,
) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {

		issuer := c.Provider.IssuerURL.Hostname()

		c.mu.RLock()
		defer c.mu.RUnlock()
		if cache, ok := c.cache[issuer]; ok {
			if !time.Now().After(cache.expires) {
				return cache.token, nil
			}
		}

		key, err := c.Provider.KeyFunc(signingMethod)(t)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		defer c.mu.Unlock()
		c.cache[issuer] = cacheJWKS{
			token:   key.(*rsa.PublicKey),
			expires: time.Now().Add(c.CacheTime),
		}

		return key, nil
	}
}

func NewCommonURL(issuerUrl string) (*url.URL, error) {
	url, err := url.Parse(issuerUrl)
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	if err != nil {
		return nil, fmt.Errorf("Invalid JWKS issuer URL: %s", issuerUrl)
	}

	url = url.JoinPath("/.well-known/jwks.json")
	return url, nil
}

func validateSigningMethod(validAlg, tokenAlg string) error {
	if validAlg != tokenAlg {
		return fmt.Errorf("Expected signing algorithm: %q , got :%q from token", validAlg, tokenAlg)
	}
	return nil
}
