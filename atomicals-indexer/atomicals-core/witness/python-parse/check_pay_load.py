# script.py

from cbor2 import loads
from util import is_sanitized_dict_whitelist_only,is_density_activated
import sys

def check_pay_load(payload, height):
    decoded_object = {}
    if payload:
        try:
            decoded_object = loads(payload)
            if not isinstance(decoded_object, dict):
                return
        except Exception as e:
            return
        if (
            not is_sanitized_dict_whitelist_only(decoded_object.get("meta", {}))
            or not is_sanitized_dict_whitelist_only(decoded_object.get("args", {}), is_density_activated(height))
            or not is_sanitized_dict_whitelist_only(decoded_object.get("ctx", {}))
            or not is_sanitized_dict_whitelist_only(decoded_object.get("init", {}), True)
        ):
            return
        print("true")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python script.py <witness_script>")
        sys.exit(1)
    payload = sys.argv[1]
    height = int(sys.argv[2]) 
    check_pay_load(payload,height)
