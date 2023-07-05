package authentication

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/emil-1003/UserManagementBackendGolang/pkg/models"
)

type Claims struct {
	Uid  int    `json:"uid"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
}

var maxAge = time.Minute * 10

func CreateJwt(secretKey string, user models.User) string {
	// Create the header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	// Create the payload
	payload := &Claims{
		Uid:  user.ID,
		Role: user.Role.Name,
		Exp:  time.Now().Add(maxAge).Unix(), // Token will expire after 10 minutes
	}

	// Encode the header and payload to JSON
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	// Base64 encode the header and payload
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create the signature
	signature := signToken(headerBase64, payloadBase64, secretKey)

	// Create the JWT token
	tokenString := fmt.Sprintf("%s.%s.%s", headerBase64, payloadBase64, signature)

	return tokenString
}

// Helper function to sign the token
func signToken(header, payload, secretKey string) string {
	signatureInput := fmt.Sprintf("%s.%s", header, payload)
	signature := hmacSHA256(signatureInput, secretKey)
	return base64.RawURLEncoding.EncodeToString(signature)
}

// Helper function to compute HMAC-SHA256 signature
func hmacSHA256(message, secretKey string) []byte {
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return h.Sum(nil)
}

// Returns parsed token and bool thats true if token is valid
func GetToken(r *http.Request) (*Claims, bool) {
	// Get the value of the Authorization header from the HTTP request
	authHeader := r.Header.Get("Authorization")

	// Check that the Authorization header is not empty and starts with the string "Bearer "
	if authHeader == "" && !strings.HasPrefix(authHeader, "Bearer") {
		return nil, false
	}

	// If the Authorization header is valid, extract the token from it by removing the "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	jwtClaims, err := ParseJWT(tokenString, os.Getenv("TOKEN_SECRET_WORD"))
	if err != nil {
		fmt.Printf("failed to parse jwt: %v", err)
		return nil, false
	}

	return jwtClaims, true
}

func ParseJWT(jwt string, secretKey string) (*Claims, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	// Base64 decode the header and payload
	_, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT claims: %w", err)
	}

	// Verify the signature
	key := []byte(secretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expectedMAC := mac.Sum(nil)

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Compare the expected and actual signature
	if !hmac.Equal(signature, expectedMAC) {
		return nil, fmt.Errorf("token signature is invalid")
	}

	var claims Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWT claims: %w", err)
	}

	// Check if the token has expired
	now := time.Now().Unix()
	if claims.Exp < now {
		return nil, fmt.Errorf("token has expired")
	}

	return &claims, nil
}

func IsUserAdmin(tokenClaims *Claims) bool {
	return tokenClaims.Role == "admin"
}

func GetUserID(r *http.Request) (int, error) {
	// Check if token is valid
	tokenClaims, ok := GetToken(r)
	if !ok {
		return 0, fmt.Errorf("Unauthorized")
	}

	return tokenClaims.Uid, nil
}
