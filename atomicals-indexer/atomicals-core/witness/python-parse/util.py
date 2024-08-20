from enumeration import Enumeration
from struct import Struct
import sys

struct_le_i = Struct('<i')
struct_le_q = Struct('<q')
struct_le_H = Struct('<H')
struct_le_I = Struct('<I')
struct_le_Q = Struct('<Q')
struct_be_Q = Struct('>Q')
struct_be_H = Struct('>H')
struct_be_I = Struct('>I')
structB = Struct('B')

unpack_le_int32_from = struct_le_i.unpack_from
unpack_le_int64_from = struct_le_q.unpack_from
unpack_le_uint16_from = struct_le_H.unpack_from
unpack_le_uint32_from = struct_le_I.unpack_from
unpack_le_uint64_from = struct_le_Q.unpack_from
unpack_be_uint16_from = struct_be_H.unpack_from
unpack_be_uint32_from = struct_be_I.unpack_from

unpack_le_uint32 = struct_le_I.unpack
unpack_le_uint64 = struct_le_Q.unpack
unpack_be_uint64 = struct_be_Q.unpack
unpack_be_uint32 = struct_be_I.unpack

pack_le_int32 = struct_le_i.pack
pack_le_int64 = struct_le_q.pack
pack_le_uint16 = struct_le_H.pack
pack_le_uint32 = struct_le_I.pack
pack_le_uint64 = struct_le_Q.pack
pack_be_uint64 = struct_be_Q.pack
pack_be_uint16 = struct_be_H.pack
pack_be_uint32 = struct_be_I.pack
pack_byte = structB.pack

hex_to_bytes = bytes.fromhex

# Parses and detects valid Atomicals protocol operations in a witness script
# Stops when it finds the first operation in the first input
def parse_protocols_operations_from_witness_for_input(txinwitness):
    '''Detect and parse all operations across the witness input arrays from a tx'''
    atomical_operation_type_map = {}
    for script in txinwitness:
        n = 0
        script_entry_len = len(script)
        if script_entry_len < 39 or script[0] != 0x20:
            continue
        found_operation_definition = False
        while n < script_entry_len - 5:
            op = script[n]
            n += 1
            # Match the pubkeyhash
            if op == 0x20 and n + 32 <= script_entry_len:
                n = n + 32
                while n < script_entry_len - 5:
                    op = script[n]
                    n += 1 
                    # Get the next if statement    
                    if op == OpCodes.OP_IF:
                        if ATOMICALS_ENVELOPE_MARKER_BYTES == script[n : n + 5].hex():
                            found_operation_definition = True
                            # Parse to ensure it is in the right format
                            operation_type, payload = parse_operation_from_script(script, n + 5)
                            if operation_type != None:
                                return operation_type, payload
                            break
                if found_operation_definition:
                    break
            else:
                break
    return None, None


# Parses the valid operations in an Atomicals script
def parse_operation_from_script(script, n):
    '''Parse an operation'''
    # Check for each protocol operation
    script_len = len(script)
    atom_op_decoded = None
    one_letter_op_len = 2
    two_letter_op_len = 3
    three_letter_op_len = 4

    # check the 3 letter protocol operations
    if n + three_letter_op_len < script_len:
        atom_op = script[n : n + three_letter_op_len].hex()
        if atom_op == "036e6674":
            atom_op_decoded = 'nft'  # nft - Mint non-fungible token
        elif atom_op == "03646674":  
            atom_op_decoded = 'dft'  # dft - Deploy distributed mint fungible token starting point
        elif atom_op == "036d6f64":  
            atom_op_decoded = 'mod'  # mod - Modify general state
        elif atom_op == "03657674": 
            atom_op_decoded = 'evt'  # evt - Message response/reply
        elif atom_op == "03646d74": 
            atom_op_decoded = 'dmt'  # dmt - Mint tokens of distributed mint type (dft)
        elif atom_op == "03646174": 
            atom_op_decoded = 'dat'  # dat - Store data on a transaction (dat)
        if atom_op_decoded:
            return atom_op_decoded, parse_atomicals_data_definition_operation(script, n + three_letter_op_len)
    
    # check the 2 letter protocol operations
    if n + two_letter_op_len < script_len:
        atom_op = script[n : n + two_letter_op_len].hex()
        if atom_op == "026674":
            atom_op_decoded = 'ft'  # ft - Mint fungible token with direct fixed supply
        elif atom_op == "02736c":  
            atom_op_decoded = 'sl'  # sl - Seal an NFT and lock it from further changes forever
        
        if atom_op_decoded:
            return atom_op_decoded, parse_atomicals_data_definition_operation(script, n + two_letter_op_len)
    
    # check the 1 letter
    if n + one_letter_op_len < script_len:
        atom_op = script[n : n + one_letter_op_len].hex()
        # Extract operation (for NFTs only)
        if atom_op == "0178":
            atom_op_decoded = 'x'  # extract - move atomical to 0'th output
        elif atom_op == "0179":
            atom_op_decoded = 'y'  # split - 

        if atom_op_decoded:
            return atom_op_decoded, parse_atomicals_data_definition_operation(script, n + one_letter_op_len)
    
    print(f'Invalid Atomicals Operation Code. Skipping... "{script[n : n + 4].hex()}"')
    return None, None


# Parses all of the push datas in a script and then concats/accumulates the bytes together
# It allows the encoding of a multi-push binary data across many pushes
def parse_atomicals_data_definition_operation(script, n):
    '''Extract the payload definitions'''
    accumulated_encoded_bytes = b''
    try:
        script_entry_len = len(script)
        while n < script_entry_len:
            op = script[n]
            n += 1
            # define the next instruction type
            if op == OpCodes.OP_ENDIF:
                break
            elif op <= OpCodes.OP_PUSHDATA4:
                data, n, dlen = parse_push_data(op, n, script)
                accumulated_encoded_bytes = accumulated_encoded_bytes + data
        return accumulated_encoded_bytes
    except Exception as e:
        raise ScriptError(f'parse_atomicals_data_definition_operation script error {e}') from None

# Parses the push datas from a bitcoin script byte sequence
def parse_push_data(op, n, script):
    data = b''
    if op <= OpCodes.OP_PUSHDATA4:
        # Raw bytes follow
        if op < OpCodes.OP_PUSHDATA1:
            dlen = op
        elif op == OpCodes.OP_PUSHDATA1:
            dlen = script[n]
            n += 1
        elif op == OpCodes.OP_PUSHDATA2:
            dlen, = unpack_le_uint16_from(script[n: n + 2])
            n += 2
        elif op == OpCodes.OP_PUSHDATA4:
            dlen, = unpack_le_uint32_from(script[n: n + 4])
            n += 4
        if n + dlen > len(script):
            raise IndexError
        data = script[n : n + dlen]
    return data, n + dlen, dlen


OpCodes = Enumeration("Opcodes", [
    ("OP_0", 0), ("OP_PUSHDATA1", 76),
    "OP_PUSHDATA2", "OP_PUSHDATA4", "OP_1NEGATE",
    "OP_RESERVED",
    "OP_1", "OP_2", "OP_3", "OP_4", "OP_5", "OP_6", "OP_7", "OP_8",
    "OP_9", "OP_10", "OP_11", "OP_12", "OP_13", "OP_14", "OP_15", "OP_16",
    "OP_NOP", "OP_VER", "OP_IF", "OP_NOTIF", "OP_VERIF", "OP_VERNOTIF",
    "OP_ELSE", "OP_ENDIF", "OP_VERIFY", "OP_RETURN",
    "OP_TOALTSTACK", "OP_FROMALTSTACK", "OP_2DROP", "OP_2DUP", "OP_3DUP",
    "OP_2OVER", "OP_2ROT", "OP_2SWAP", "OP_IFDUP", "OP_DEPTH", "OP_DROP",
    "OP_DUP", "OP_NIP", "OP_OVER", "OP_PICK", "OP_ROLL", "OP_ROT",
    "OP_SWAP", "OP_TUCK",
    "OP_CAT", "OP_SUBSTR", "OP_LEFT", "OP_RIGHT", "OP_SIZE",
    "OP_INVERT", "OP_AND", "OP_OR", "OP_XOR", "OP_EQUAL", "OP_EQUALVERIFY",
    "OP_RESERVED1", "OP_RESERVED2",
    "OP_1ADD", "OP_1SUB", "OP_2MUL", "OP_2DIV", "OP_NEGATE", "OP_ABS",
    "OP_NOT", "OP_0NOTEQUAL", "OP_ADD", "OP_SUB", "OP_MUL", "OP_DIV", "OP_MOD",
    "OP_LSHIFT", "OP_RSHIFT", "OP_BOOLAND", "OP_BOOLOR", "OP_NUMEQUAL",
    "OP_NUMEQUALVERIFY", "OP_NUMNOTEQUAL", "OP_LESSTHAN", "OP_GREATERTHAN",
    "OP_LESSTHANOREQUAL", "OP_GREATERTHANOREQUAL", "OP_MIN", "OP_MAX",
    "OP_WITHIN",
    "OP_RIPEMD160", "OP_SHA1", "OP_SHA256", "OP_HASH160", "OP_HASH256",
    "OP_CODESEPARATOR", "OP_CHECKSIG", "OP_CHECKSIGVERIFY", "OP_CHECKMULTISIG",
    "OP_CHECKMULTISIGVERIFY",
    "OP_NOP1",
    "OP_CHECKLOCKTIMEVERIFY", "OP_CHECKSEQUENCEVERIFY"
])

# The maximum height difference between the commit and reveal transactions of any Atomical mint
# This is used to limit the amount of cache we would need in future optimizations.
MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS = 100

# The maximum height difference between the commit and reveal transactions of any of the sub(realm) mints
# This is needed to prevent front-running of realms. 
MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS = 3

MINT_SUBNAME_RULES_BECOME_EFFECTIVE_IN_BLOCKS = 1 # magic number that rules become effective in one block

# The path namespace to look for when determining what price/regex patterns are allowed if any
SUBREALM_MINT_PATH = 'subrealms'

# The path namespace to look for when determining what price/regex patterns are allowed if any
DMINT_PATH = 'dmint'

# The maximum height difference between the reveal transaction of the winning subrealm claim and the blocks to pay the necessary fee to the parent realm
# It is intentionally made longer since it may take some time for the purchaser to get the funds together
MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS = 15 # ~2.5 hours.
# MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS and therefore it gives the user about 1.5 hours to make the payment after they know
# that they won the realm (and no one else can claim/reveal)

# The Envelope is for the reveal script and also the op_return payment markers
# "atom" 
ATOMICALS_ENVELOPE_MARKER_BYTES = '0461746f6d'

# Limit the smallest payment amount allowed for a subrealm
SUBNAME_MIN_PAYMENT_DUST_LIMIT = 0 # It can be possible to do free

# Maximum size of the rules of a subrealm or container dmint rule set array
MAX_SUBNAME_RULE_SIZE_LEN = 100000
# Maximum number of minting rules allowed
MAX_SUBNAME_RULE_ENTRIES = 100

# Minimum amount in satoshis for a DFT mint operation. Set at satoshi dust of 546
DFT_MINT_AMOUNT_MIN = 546

# Maximum amount in satoshis for the DFT mint operation. Set at 1 BTC for the ballers
DFT_MINT_AMOUNT_MAX = 100000000

# The minimum number of DFT max_mints. Set at 1
DFT_MINT_MAX_MIN_COUNT = 1
# The maximum number (legacy) of DFT max_mints. Set at 500,000 mints mainly for efficieny reasons in legacy.
DFT_MINT_MAX_MAX_COUNT_LEGACY = 500000
# The maximum number of DFT max_mints (after legacy 'DENSITY' update). Set at 21,000,000 max mints.
DFT_MINT_MAX_MAX_COUNT_DENSITY = 21000000

# This would never change, but we put it as a constant for clarity
DFT_MINT_HEIGHT_MIN = 0
# This value would never change, it's added in case someone accidentally tries to use a unixtime
DFT_MINT_HEIGHT_MAX = 10000000 # 10 million blocks
  
def is_sanitized_dict_whitelist_only(d: dict, allow_bytes=False):
    if not isinstance(d, dict):
        return False
    for k, v in d.items():
        if isinstance(v, dict):
            return is_sanitized_dict_whitelist_only(v, allow_bytes)
        if not allow_bytes and isinstance(v, bytes):
            print( f"parse {k} {v} ..." )
            return False
        if (
            not isinstance(v, int)
            and not isinstance(v, float)
            and not isinstance(v, list)
            and not isinstance(v, str)
            and not isinstance(v, bytes)
        ):
            # Prohibit anything except int, float, lists, strings and bytes
            return False
    return True

def is_density_activated(height: int):
    ATOMICALS_ACTIVATION_HEIGHT_DENSITY = 828128
    if height >= ATOMICALS_ACTIVATION_HEIGHT_DENSITY:
        return True
    return False