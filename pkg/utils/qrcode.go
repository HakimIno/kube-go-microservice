package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"kube/internal/middleware"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type bufferCloser struct {
	*bytes.Buffer
}

func (bc *bufferCloser) Close() error {
	return nil
}

func GenerateQRCodeWithLogo(data string, logoPath string, size int) (string, error) {
	qrc, err := qrcode.New(data)

	if err != nil {
		middleware.LogError("Could not generate QRCode", err)
		return "", err
	}

	var buf bytes.Buffer
	bufferWriter := &bufferCloser{&buf}

	halftonePath := "assets/dog.jpg"
	if _, err := os.Stat(halftonePath); os.IsNotExist(err) {
		middleware.LogWarn(fmt.Sprintf("Halftone image file %s not found", halftonePath))
	}

	options := []standard.ImageOption{
		standard.WithHalftone(halftonePath),
		standard.WithQRWidth(21),
	}

	w := standard.NewWithWriter(bufferWriter, options...)

	if err = qrc.Save(w); err != nil {
		middleware.LogError("Could not save QR code image", err)
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
