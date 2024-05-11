package witness

import (
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/utils"
)

type PayLoad struct {
	// ImagePng map[string]string `cbor:"image.png"` // Todo: there are some err when unmarshal ImagePng
	Args                       *Args            `cbor:"args"`
	Meta                       *Meta            `cbor:"meta"`
	TotalAmountToSkipPotential map[string]int64 // key: locationID
	Main                       []byte
}

// parse_atomicals_data_definition_operation
func parseAtomicalsData(script []byte, startIndex int64) ([]byte, error) {
	payloadBytes := []byte{}
	for startIndex < int64(len(script)) {
		op := script[startIndex]
		startIndex++
		// define the next instruction type
		if op == OP_ENDIF {
			break
		}
		if op <= OP_PUSHDATA4 {
			// data, dlen := parsePushData(op, startIndex, script)
			data := []byte{}
			dlen := int64(0)
			if op <= OP_PUSHDATA4 {
				// Raw bytes follow
				if op < OP_PUSHDATA1 {
					dlen = int64(op)
				} else if op == OP_PUSHDATA1 {
					dlen = int64(script[startIndex])
					startIndex++
				} else if op == OP_PUSHDATA2 {
					dlen = int64(utils.UnpackLeUint16From(script[startIndex : startIndex+2]))
					startIndex += 2
				} else if op == OP_PUSHDATA4 {
					dlen = int64(utils.UnpackLeUint32From(script[startIndex : startIndex+4]))
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
		if m.Args.RequestSubRealm != "" {
			request_counter += 1
		}
		if m.Args.RequestContainer != "" {
			request_counter += 1
		}
		if m.Args.RequestTicker != "" {
			request_counter += 1
		}
		if m.Args.RequestDmitem != "" {
			request_counter += 1
		}
		if request_counter > 1 {
			return false
		}
	}
	return true
}
func (m *PayLoad) CheckRequest() bool {
	if m.Args == nil {
		return false
	} else {
		request_counter := 0 // # Ensure that only one of the following may be requested || fail
		if m.Args.RequestRealm != "" {
			request_counter += 1
		}
		if m.Args.RequestSubRealm != "" {
			request_counter += 1
		}
		if m.Args.RequestContainer != "" {
			request_counter += 1
		}
		if m.Args.RequestTicker != "" {
			request_counter += 1
		}
		if m.Args.RequestDmitem != "" {
			request_counter += 1
		}
		if request_counter > 1 {
			return false
		}
	}
	return true
}

type Args struct {
	// utils
	Nonce int64 `cbor:"nonce"`
	Time  int64 `cbor:"time"`

	// optional
	Bitworkc string `cbor:"bitworkc"`
	Bitworkr string `cbor:"bitworkr"`

	I       bool             `cbor:"i"`
	Main    string           `cbor:"main"`
	Proof   []Proof          `cbor:"proof"`
	Parents map[string]int64 `cbor:"parents"` // key: parent_atomical_id, value: , haven't catch this param, used in operation:nft

	// dft & ft
	RequestTicker string `cbor:"request_ticker"`

	// dft
	MintAmount   int64  `cbor:"mint_amount"`
	MintHeight   int64  `cbor:"mint_height"`
	MaxMints     int64  `cbor:"max_mints"`
	MintBitworkc string `cbor:"mint_bitworkc"`
	MintBitworkr string `cbor:"mint_bitworkr"`
	Md           string `cbor:"md"` // emu:"", "0", "1"
	Bv           string `cbor:"bv"`
	Bci          string `cbor:"bci"`
	Bri          string `cbor:"bri"`
	Bcs          int64  `cbor:"bcs"`
	Brs          int64  `cbor:"brs"`
	Maxg         int64  `cbor:"maxg"`

	// dmt
	MintTicker string `cbor:"mint_ticker"` // mint ft name

	// nft: realm
	RequestRealm string `cbor:"request_realm"`

	// nft: subrealm
	RequestSubRealm string               `cbor:"request_subrealm"`
	ClaimType       NftSubrealmClaimType `cbor:"claim_type"`   // enum: "direct" "rule"
	ParentRealm     string               `cbor:"parent_realm"` // ParentRealm atomicalsID

	// nft: container
	RequestContainer string `cbor:"request_container"`

	// nft: dmitem
	RequestDmitem   string `cbor:"request_dmitem"`   // item num in ParentContainer
	ParentContainer string `cbor:"parent_container"` // ParentContainer atomicalsID

}

type NftSubrealmClaimType string

const (
	Direct NftSubrealmClaimType = "direct"
	Rule   NftSubrealmClaimType = "rule"
)

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

type Proof struct {
	D []byte `cbor:"d"`
	P bool   `cbor:"p"`
}

const (
	// Constants for Bitcoin Script opcodes
	OP_0                   = 0
	OP_PUSHDATA1           = 76
	OP_PUSHDATA2           = 77
	OP_PUSHDATA4           = 78
	OP_1NEGATE             = 79
	OP_RESERVED            = 80
	OP_1                   = 81
	OP_2                   = 82
	OP_3                   = 83
	OP_4                   = 84
	OP_5                   = 85
	OP_6                   = 86
	OP_7                   = 87
	OP_8                   = 88
	OP_9                   = 89
	OP_10                  = 90
	OP_11                  = 91
	OP_12                  = 92
	OP_13                  = 93
	OP_14                  = 94
	OP_15                  = 95
	OP_16                  = 96
	OP_NOP                 = 97
	OP_VER                 = 98
	OP_IF                  = 99
	OP_NOTIF               = 100
	OP_VERIF               = 101
	OP_VERNOTIF            = 102
	OP_ELSE                = 103
	OP_ENDIF               = 104
	OP_VERIFY              = 105
	OP_RETURN              = 106
	OP_TOALTSTACK          = 107
	OP_FROMALTSTACK        = 108
	OP_2DROP               = 109
	OP_2DUP                = 110
	OP_3DUP                = 111
	OP_2OVER               = 112
	OP_2ROT                = 113
	OP_2SWAP               = 114
	OP_IFDUP               = 115
	OP_DEPTH               = 116
	OP_DROP                = 117
	OP_DUP                 = 118
	OP_NIP                 = 119
	OP_OVER                = 120
	OP_PICK                = 121
	OP_ROLL                = 122
	OP_ROT                 = 123
	OP_SWAP                = 124
	OP_TUCK                = 125
	OP_CAT                 = 126
	OP_SUBSTR              = 127
	OP_LEFT                = 128
	OP_RIGHT               = 129
	OP_SIZE                = 130
	OP_INVERT              = 131
	OP_AND                 = 132
	OP_OR                  = 133
	OP_XOR                 = 134
	OP_EQUAL               = 135
	OP_EQUALVERIFY         = 136
	OP_RESERVED1           = 137
	OP_RESERVED2           = 138
	OP_1ADD                = 139
	OP_1SUB                = 140
	OP_2MUL                = 141
	OP_2DIV                = 142
	OP_NEGATE              = 143
	OP_ABS                 = 144
	OP_NOT                 = 145
	OP_0NOTEQUAL           = 146
	OP_ADD                 = 147
	OP_SUB                 = 148
	OP_MUL                 = 149
	OP_DIV                 = 150
	OP_MOD                 = 151
	OP_LSHIFT              = 152
	OP_RSHIFT              = 153
	OP_BOOLAND             = 154
	OP_BOOLOR              = 155
	OP_NUMEQUAL            = 156
	OP_NUMEQUALVERIFY      = 157
	OP_NUMNOTEQUAL         = 158
	OP_LESSTHAN            = 159
	OP_GREATERTHAN         = 160
	OP_LESSTHANOREQUAL     = 161
	OP_GREATERTHANOREQUAL  = 162
	OP_MIN                 = 163
	OP_MAX                 = 164
	OP_WITHIN              = 165
	OP_RIPEMD160           = 166
	OP_SHA1                = 167
	OP_SHA256              = 168
	OP_HASH160             = 169
	OP_HASH256             = 170
	OP_CODESEPARATOR       = 171
	OP_CHECKSIG            = 172
	OP_CHECKSIGVERIFY      = 173
	OP_CHECKMULTISIG       = 174
	OP_CHECKMULTISIGVERIFY = 175
	OP_NOP1                = 176
	OP_CHECKLOCKTIMEVERIFY = 177
	OP_CHECKSEQUENCEVERIFY = 178
)
