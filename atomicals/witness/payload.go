package witness

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/errors"
)

type PayLoad struct {
	// ImagePng map[string]string `cbor:"image.png"` // Todo: there are some err when unmarshal ImagePng
	Args *Args `cbor:"args"`
	Meta *Meta `cbor:"meta"`
}

// parse_atomicals_data_definition_operation
func parseAtomicalsData(script []byte, startIndex int64) ([]byte, error) {
	payloadBytes := []byte{}
	for startIndex < int64(len(script)) {
		op := script[startIndex]
		startIndex++
		// define the next instruction type
		if op == common.OP_ENDIF {
			break
		}
		if op <= common.OP_PUSHDATA4 {
			// data, dlen := parsePushData(op, startIndex, script)
			data := []byte{}
			dlen := int64(0)
			if op <= common.OP_PUSHDATA4 {
				// Raw bytes follow
				if op < common.OP_PUSHDATA1 {
					dlen = int64(op)
				} else if op == common.OP_PUSHDATA1 {
					dlen = int64(script[startIndex])
					startIndex++
				} else if op == common.OP_PUSHDATA2 {
					dlen = int64(common.UnpackLeUint16From(script[startIndex : startIndex+2]))
					startIndex += 2
				} else if op == common.OP_PUSHDATA4 {
					dlen = int64(common.UnpackLeUint32From(script[startIndex : startIndex+4]))
					startIndex += 4
				}
				if int64(startIndex+dlen) > int64(len(script)) {
					return nil, errors.ErrInvalidAtomicalsData
				}
				data = script[startIndex : startIndex+dlen]
				startIndex = startIndex + dlen
			}
			payloadBytes = append(payloadBytes, data...)
		}
	}
	return payloadBytes, nil
}

func (m *PayLoad) check() bool {
	if m.Args == nil {
		return false
	} else {
		request_counter := 0 // # Ensure that only one of the following may be requested || fail
		if m.Args.RequestRealm != "" {
			request_counter += 1
		}
		if m.Args.RequestContainer != "" {
			request_counter += 1
		}
		if request_counter > 1 {
			return false
		}
	}
	return true
}

type Args struct {
	// common
	Nonce    int64  `cbor:"nonce"`
	Time     int64  `cbor:"time"`
	Bitworkc string `cbor:"bitworkc"`

	// dft
	MintAmount    float64 `cbor:"mint_amount"`
	MintHeight    int64   `cbor:"mint_height"`
	MaxMints      int64   `cbor:"max_mints"`
	MintBitworkc  string  `cbor:"mint_bitworkc"`
	RequestTicker string  `cbor:"request_ticker"`

	// nft
	RequestRealm     string `cbor:"request_realm"`
	RequestSubRealm  string `cbor:"request_subrealm"`
	RequestContainer string `cbor:"request_container"`
	RequestDmitem    string `cbor:"request_dmitem"`

	// dmt
	MintTicker string `cbor:"mint_ticker"` // mint ft name
}

type Meta struct {
	Name        string `cbor:"name"`
	Description string `cbor:"description"`
	Legal       *Legal `cbor:"legal"`
}

type Legal struct {
	Terms string `cbor:"terms"`
}

type ImagePng struct {
	CT string `cbor:"$ct"`
	B  string `cbor:"$b"`
}
