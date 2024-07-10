# script.py

from cbor2 import loads
from util import parse_protocols_operations_from_witness_for_input
import sys

def parse(witenss_script):
    # print("python parse witenss_script:", witenss_script)
    op_name, payload = parse_protocols_operations_from_witness_for_input([bytes.fromhex(witenss_script)])
    decoded_object = loads(payload)
    print("python parse op_name:", op_name)
    print("python parse decoded_object:", decoded_object)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python script.py <witness_script>")
        sys.exit(1)
    witenss_script = sys.argv[1]
    parse(witenss_script)
