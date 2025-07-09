# __init__.py

"""
Krock32 module 0.1.1

Implementation of Crockford's base32 alphabet encoding, decoding,
and pretty printing
"""

# pylama:ignore=W0611
from krock32.decode import Decoder
from krock32.encode import Encoder

__version__ = "0.1.1"
