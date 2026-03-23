package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID, domain.RoleUser, testSecret)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	userID := uuid.New()

	type expected struct {
		userID uuid.UUID
		role   domain.Role
	}

	tests := []struct {
		name      string
		makeToken func() string
		wantErr   bool
		expected  *expected
	}{
		{
			name: "valid token",
			makeToken: func() string {
				token, _ := GenerateToken(userID, domain.RoleUser, testSecret)
				return token
			},
			expected: &expected{
				userID: userID,
				role:   domain.RoleUser,
			},
		},
		{
			name: "wrong secret",
			makeToken: func() string {
				token, _ := GenerateToken(userID, domain.RoleUser, "other-secret")
				return token
			},
			wantErr: true,
		},
		{
			name: "malformed token",
			makeToken: func() string {
				return "not.a.token"
			},
			wantErr: true,
		},
		{
			name: "empty token",
			makeToken: func() string {
				return ""
			},
			wantErr: true,
		},
		{
			name: "expired token",
			makeToken: func() string {
				claims := Claims{
					UserID: userID,
					Role:   domain.RoleUser,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signed, _ := token.SignedString([]byte(testSecret))
				return signed
			},
			wantErr: true,
		},
		{
			name: "token expires now",
			makeToken: func() string {
				claims := Claims{
					UserID: userID,
					Role:   domain.RoleUser,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now()),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signed, _ := token.SignedString([]byte(testSecret))
				return signed
			},
			wantErr: true,
		},
		{
			name: "no expiration (allowed by current impl)",
			makeToken: func() string {
				claims := Claims{
					UserID: userID,
					Role:   domain.RoleUser,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signed, _ := token.SignedString([]byte(testSecret))
				return signed
			},
			expected: &expected{
				userID: userID,
				role:   domain.RoleUser,
			},
		},
		{
			name: "tampered token",
			makeToken: func() string {
				token, _ := GenerateToken(userID, domain.RoleUser, testSecret)
				parts := strings.Split(token, ".")
				parts[1] = "tamperedpayload"
				return strings.Join(parts, ".")
			},
			wantErr: true,
		},
		{
			name: "wrong signing method",
			makeToken: func() string {
				claims := Claims{
					UserID: userID,
					Role:   domain.RoleUser,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				signed, _ := token.SignedString([]byte("fake-key"))
				return signed
			},
			wantErr: true,
		},
		{
			name: "zero UUID",
			makeToken: func() string {
				token, _ := GenerateToken(uuid.Nil, domain.RoleUser, testSecret)
				return token
			},
			expected: &expected{
				userID: uuid.Nil,
				role:   domain.RoleUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.makeToken(), testSecret)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, claims)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, claims)

			assert.Equal(t, tt.expected.userID, claims.UserID)
			assert.Equal(t, tt.expected.role, claims.Role)

			if claims.ExpiresAt != nil {
				assert.True(t, claims.ExpiresAt.After(time.Now()))
			}
			if claims.IssuedAt != nil {
				assert.False(t, claims.IssuedAt.Time.IsZero())
			}
		})
	}
}
