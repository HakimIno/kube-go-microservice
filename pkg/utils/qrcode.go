package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"

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
		fmt.Printf("could not generate QRCode: %v", err)
		return "", err
	}

	var buf bytes.Buffer
	bufferWriter := &bufferCloser{&buf}

	options := []standard.ImageOption{
		standard.WithCircleShape(),
		standard.WithLogoImageFileJPEG(logoPath),
		standard.WithBorderWidth(10),
	}

	w := standard.NewWithWriter(bufferWriter, options...)

	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
