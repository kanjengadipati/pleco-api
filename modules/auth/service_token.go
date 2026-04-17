package auth

import (
	"errors"
	"time"

	"go-auth-app/utils"

	"gorm.io/gorm"
)

func (s *authService) Logout(userID uint, deviceID string) error {
	token, err := s.RefreshTokenRepo.FindByUserAndDevice(userID, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return s.RefreshTokenRepo.DeleteByID(token.ID)
}

func (s *authService) LogoutAll(userID uint) error {
	return s.RefreshTokenRepo.DeleteByUser(userID)
}

func (s *authService) RefreshToken(oldRefreshToken string) (*AuthTokens, error) {
	claims, err := s.JWT.ValidateToken(oldRefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims["type"] != TokenRefresh {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user id in token claims")
	}
	uid := uint(userID)

	tokens, err := s.RefreshTokenRepo.FindByUser(uid)
	if err != nil {
		return nil, err
	}
	oldHash := utils.HashToken(oldRefreshToken)

	var matchedIndex = -1
	for i := range tokens {
		if tokens[i].TokenHash == oldHash {
			matchedIndex = i
			break
		}
	}
	if matchedIndex == -1 {
		return nil, errors.New("invalid refresh token")
	}

	matchedToken := tokens[matchedIndex]
	if time.Now().After(matchedToken.ExpiredAt) {
		return nil, errors.New("refresh token expired")
	}

	if err := s.RefreshTokenRepo.DeleteByID(matchedToken.ID); err != nil {
		return nil, err
	}

	user, err := s.UserRepo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(uid, user.Role, matchedToken.DeviceID, matchedToken.UserAgent, matchedToken.IPAddress)
}
