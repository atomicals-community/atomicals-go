# script.py

from cbor2 import loads
from util import parse_protocols_operations_from_witness_for_input,is_sanitized_dict_whitelist_only,is_density_activated
import sys

def parse(witness_script, height):
    op_name, payload = parse_protocols_operations_from_witness_for_input([bytes.fromhex(witness_script)])
    if not op_name:
        return
    decoded_object = {}
    if payload:
        # Ensure that the payload is cbor encoded dictionary or empty
        try:
            decoded_object = loads(payload)
            if not isinstance(decoded_object, dict):
                print("false")
                return
        except Exception as e:
            print("false")
            return
        # Also enforce that if there are meta, args, or ctx fields that they must be dicts
        # This is done to ensure that these fields are always easily parseable and do not contain unexpected data which could cause parsing problems later
        # Ensure that they are not allowed to contain bytes like objects
        if (
            not is_sanitized_dict_whitelist_only(decoded_object.get("meta", {}))
            or not is_sanitized_dict_whitelist_only(decoded_object.get("args", {}), is_density_activated(height))
            or not is_sanitized_dict_whitelist_only(decoded_object.get("ctx", {}))
            or not is_sanitized_dict_whitelist_only(decoded_object.get("init", {}), True)
        ):
            print("false")
            return
    else:  
        print("false")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <witness_script>")
        sys.exit(1)
    witness_script = sys.argv[1]
    height = int(sys.argv[2])  # Convert height to integer
    parse(witness_script,height)
