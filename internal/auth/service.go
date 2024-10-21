package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

type service struct {
	gcpClientSecret string
}

func NewService(gcs string) Service {
	return &service{
		gcpClientSecret: gcs,
	}
}

func (s *service) CreateJWT(user *usermodel.User) (string, error) {
	jwtClaims := jwt.MapClaims{}
	jwtClaims = addUserToJWTClaims(jwtClaims, user)
	jwtClaims = addExpiryToJWTClaims(jwtClaims, 30*24*time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	jwtString, err := token.SignedString([]byte(s.gcpClientSecret))
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

func addUserToJWTClaims(jwtClaims jwt.MapClaims, user *usermodel.User) jwt.MapClaims {
	jwtClaims["displayName"] = user.DisplayName
	jwtClaims["displayImageURL"] = user.DisplayImageURL
	return jwtClaims
}

func addExpiryToJWTClaims(jwtClaims jwt.MapClaims, jwtDuration time.Duration) jwt.MapClaims {
	now := time.Now().UTC()
	expiration := now.Add(jwtDuration)
	jwtClaims["exp"] = expiration.Unix()
	return jwtClaims
}

func (s *service) ValidateJWT(jwtString string) error {
	jwtToken, err := jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) {
		return s.gcpClientSecret, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return err
	}
	return isJWTExpired(jwtToken)
}

func isJWTExpired(jwtToken *jwt.Token) error {
	now := time.Now().UTC()
	exp, err := jwtToken.Claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if now.Unix() > exp.Unix() {
		return ierrors.NewError(ierrors.Unauthorized, errors.New("JWT is expired"))
	}
	return nil
}
