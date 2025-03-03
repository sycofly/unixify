package auth

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"time"

	"github.com/home/unixify/internal/models"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret generates a new TOTP secret for a user
func (s *Service) GenerateTOTPSecret(username string) (*models.TOTPSetupResponse, error) {
	issuer := s.config.Server.TOTPIssuer
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
	})
	if err != nil {
		return nil, err
	}

	// Get the QR code as an image
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	
	// Encode the image to PNG
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	
	// Convert to base64 for the frontend
	qrCode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	
	response := &models.TOTPSetupResponse{
		Secret: key.Secret(),
		QRCode: qrCode,
	}
	
	return response, nil
}

// VerifyTOTP verifies a TOTP code against a secret
func (s *Service) VerifyTOTP(secret, code string) bool {
	valid, err := totp.ValidateCustom(
		code,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	
	if err != nil {
		return false
	}
	
	return valid
}