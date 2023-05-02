package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/veraison/ccatoken"
	"github.com/veraison/go-cose"
	"github.com/veraison/psatoken"
)

func convert_to_cbor() error {
	data, err := os.ReadFile("input/base64-token.txt")
	if err != nil {
		return fmt.Errorf("error in reading")
	}
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	_, err = base64.StdEncoding.Decode(dbuf, data)
	if err != nil {
		return fmt.Errorf("error in decoding")
	}
	os.WriteFile("input/cbor-token.cbor", dbuf, (os.ModeAppend | 0x3FF))
	return nil
}

func main() {
	var EvidenceIn ccatoken.Evidence

	err := convert_to_cbor()
	if err != nil {
		fmt.Println("unable to convert :%w", err)
	}
	tokenBytes, err := os.ReadFile("input/cbor-token.cbor")

	err = EvidenceIn.FromCBOR(tokenBytes)
	if err != nil {
		fmt.Println("unable to read CBOR token Bytes %w", err)
	}

	tokenJSON, err := EvidenceIn.MarshalJSON()
	err = os.WriteFile("input/Token.json", tokenJSON, (os.ModeAppend | 0x3FF))
	if buf != nil {
		convert_from_hex_COSE_CBOR()
	}
}

// The input buf is taken from rmm_delegated_attest.c Line 89
func convert_from_hex_COSE_CBOR() error {
	// decode platform
	pSign1 := cose.NewSign1Message()

	err := pSign1.UnmarshalCBOR(buf)
	if err != nil {
		return fmt.Errorf("unable to set the token: %w", err)
	}

	PlatformClaims, err := psatoken.DecodeClaims(pSign1.Payload)

	jd, err := PlatformClaims.ToJSON()
	err = os.WriteFile("input/TokenFromCoSE.json", jd, (os.ModeAppend | 0x3FF))

	return nil
}
