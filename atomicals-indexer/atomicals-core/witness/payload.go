package witness

import (
	"encoding/hex"

	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/utils"
	"github.com/fxamacker/cbor/v2"
)

type PayLoad struct {
	Args      *Args  `cbor:"args"`
	A         int    `cbor:"$a"`        // for mod
	Dmint     *Dmint `cbor:"dmint"`     // for mod
	Subrealms *Dmint `cbor:"subrealms"` // for mod
	Meta      *Meta  `cbor:"meta"`
}

type Args struct {
	Nonce int64 `cbor:"nonce"`
	Time  int64 `cbor:"time"`

	// optional
	Bitworkc string `cbor:"bitworkc"`
	Bitworkr string `cbor:"bitworkr"`

	Immutable     bool              `cbor:"$immutable"`
	I             bool              `cbor:"i"`
	Main          string            `cbor:"main"`
	DynamicFields map[string][]byte `cbor:"-"` // use for Main, sometimes, Main is "image.png" or unsure name... so bad a rule
	Proof         []Proof           `cbor:"proof"`
	Parents       map[string]int64  `cbor:"$parents"` // key: parent_atomical_id, value: , haven't catch this param, used in operation:nft

	// for dft & ft
	RequestTicker string `cbor:"request_ticker"`

	// for dft
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

	// for dmt
	MintTicker string `cbor:"mint_ticker"` // mint ft name

	// for nft: realm
	RequestRealm string `cbor:"request_realm"`

	// for nft: subrealm
	RequestSubRealm string               `cbor:"request_subrealm"`
	ClaimType       NftSubrealmClaimType `cbor:"claim_type"`   // enum: "direct" "rule"
	ParentRealm     string               `cbor:"parent_realm"` // ParentRealm atomicalsID

	// for nft: container
	RequestContainer string `cbor:"request_container"`

	// for nft: dmitem
	RequestDmitem   string `cbor:"request_dmitem"`   // item num in ParentContainer
	ParentContainer string `cbor:"parent_container"` // ParentContainer atomicalsID

	// for y
	TotalAmountToSkipPotential map[string]int64 // key: locationID
}

type NftSubrealmClaimType string

const (
	Direct NftSubrealmClaimType = "direct"
	Rule   NftSubrealmClaimType = "rule"
)

type Meta struct {
	Name        string `cbor:"name"`
	Description string `cbor:"description"`
	Legal       struct {
		Terms string `cbor:"terms"`
	} `cbor:"legal"`
}

type Proof struct {
	D string `cbor:"d"`
	P bool   `cbor:"p"`
}

type Dmint struct {
	A          int64       `cbor:"$a"`
	V          string      `cbor:"v"`
	Items      int64       `cbor:"items"`
	Rules      []*RuleInfo `cbor:"rules"`
	Merkle     string      `cbor:"merkle"`
	Immutable  bool        `cbor:"immutable"`
	MintHeight int64       `cbor:"mint_height"`
}

type RuleInfo struct {
	P        string             `cbor:"p"`
	O        map[string]*Output `cbor:"o"` // key:
	Bitworkc string             `cbor:"bitworkc"`
	Bitworkr string             `cbor:"bitworkr"`
}
type Output struct {
	ID string `cbor:"id"`
	V  int64  `cbor:"v"`
}

type Subrealms struct {
	Rules []*RuleInfo `cbor:"rules"`
}

func parseOperationAndPayLoad(script string) (string, *PayLoad, error) {
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		return "", nil, err
	}
	scriptEntryLen := int64(len(scriptBytes))
	if scriptEntryLen < 39 || scriptBytes[0] != 0x20 {
		return "", nil, errors.ErrInvalidWitnessScriptLength
	}
	pkFlag := scriptBytes[0]
	if pkFlag != 0x20 {
		return "", nil, errors.ErrInvalidWitnessScriptPkFlag
	}
	for index := int64(35); index < scriptEntryLen-6; index++ {
		opFlag := scriptBytes[index]
		if opFlag != OP_IF {
			continue
		}
		if hex.EncodeToString(scriptBytes[index+1:index+6]) != utils.ATOMICALS_ENVELOPE_MARKER_BYTES {
			continue
		}
		operation, startIndex := parseAtomicalsOperation(scriptBytes, index+6)
		if operation == "" {
			continue
		}
		payloadBytes, err := parseAtomicalsData(scriptBytes, startIndex)
		if err != nil {
			return "", nil, err
		}
		if payloadBytes == nil {
			continue
		}
		// get DynamicFields[main]
		payload := &PayLoad{
			Args: &Args{
				DynamicFields:              map[string][]byte{},
				TotalAmountToSkipPotential: make(map[string]int64, 0),
			},
		}
		if err := cbor.Unmarshal(payloadBytes, payload); err != nil {
			return "", nil, err
		}
		tempMap := map[string]interface{}{}
		if err := cbor.Unmarshal(payloadBytes, &tempMap); err != nil {
			return "", nil, err
		}
		if _, ok := tempMap[payload.Args.Main]; ok {
			payload.Args.DynamicFields[payload.Args.Main] = tempMap[payload.Args.Main].([]byte)
		}
		// get TotalAmountToSkipPotential
		if operation == "y" {
			if err := cbor.Unmarshal(payloadBytes, &payload.Args.TotalAmountToSkipPotential); err != nil {
				return "", nil, err
			}
		}
		return operation, payload, nil
	}
	return "", nil, errors.ErrOptionNotFound
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
