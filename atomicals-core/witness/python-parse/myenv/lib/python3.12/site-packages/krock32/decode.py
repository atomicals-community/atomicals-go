# decode.py

"""
Krock32 Decoder

Class definition for the Krock32 decoder, taking a Crockford Base32 string
and turning it into bytes-data.
"""

from collections import namedtuple


class DecoderAlreadyFinalizedException(Exception):
    pass


class DecoderInvalidStringLengthException(Exception):
    pass


class DecoderNonZeroCarryException(Exception):
    pass


class DecoderChecksumException(Exception):
    pass


class Decoder:
    def __init__(
        self,
        strict: bool = False,
        ignore_non_alphabet: bool = True,
        checksum: bool = False,
    ):
        self._string_buffer: str = ""
        self._bytearray: bytearray = bytearray()
        self._strict: bool = strict
        self._ignore_non_alphabet: bool = ignore_non_alphabet
        self._alphabet = self._make_alphabet(
            "0123456789ABCDEFGHJKMNPQRSTVWXYZ*~$=U", strict=self._strict
        )
        self._is_finished: bool = False
        self._p_byte = namedtuple("ProcessedByte", ["byte", "carry"])
        self._do_checksum: bool = checksum
        self._checksum: int = 0

    def _make_alphabet(self, alphabet_string: str, strict: bool) -> dict:
        alphabet = {}
        for i, x in enumerate(alphabet_string):
            alphabet[x.upper()] = i
            if not strict:
                alphabet[x.lower()] = i
        if not strict:
            alphabet["O"], alphabet["o"] = 0, 0
            (alphabet["I"], alphabet["i"], alphabet["L"], alphabet["l"]) = 1, 1, 1, 1
        return alphabet

    def _update_checksum(self, byte: int):
        if not self._do_checksum:
            return
        self._checksum = ((self._checksum << 8) + byte) % 37

    def _decode_first_byte(self, symbols: str) -> tuple:
        byte = self._alphabet.get(symbols[0]) << 3
        second_sym = self._alphabet.get(symbols[1])
        byte += second_sym >> 2
        carry = second_sym & 0b11
        self._update_checksum(byte)
        return self._p_byte(byte=byte, carry=carry)

    def _decode_second_byte(self, symbols: str, carry: int) -> tuple:
        byte = (self._alphabet.get(symbols[0]) << 1) + (carry << 6)
        second_sym = self._alphabet.get(symbols[1])
        byte += second_sym >> 4
        carry = second_sym & 0b1111
        self._update_checksum(byte)
        return self._p_byte(byte=byte, carry=carry)

    def _decode_third_byte(self, symbols: str, carry: int) -> tuple:
        sym = self._alphabet.get(symbols)
        byte = (sym >> 1) + (carry << 4)
        carry = sym & 1
        self._update_checksum(byte)
        return self._p_byte(byte=byte, carry=carry)

    def _decode_fourth_byte(self, symbols: str, carry: int) -> tuple:
        byte = (self._alphabet.get(symbols[0]) << 2) + (carry << 7)
        second_sym = self._alphabet.get(symbols[1])
        byte += second_sym >> 3
        carry = second_sym & 0b111
        self._update_checksum(byte)
        return self._p_byte(byte=byte, carry=carry)

    def _decode_fifth_byte(self, symbols: str, carry: int) -> tuple:
        byte = (carry << 5) + self._alphabet.get(symbols)
        self._update_checksum(byte)
        return self._p_byte(byte=byte, carry=0)

    def _return_quantum(
        self, quantum: str, p_byte: tuple, array: bytearray
    ) -> bytearray:
        if p_byte.carry == 0:
            return array
        else:
            raise DecoderNonZeroCarryException(
                "Quantum %s decoded with non-zero carry %i" % (quantum, p_byte.carry)
            )

    def _decode_quantum(self, quantum: str) -> bytearray:
        if not len(quantum) in [2, 4, 5, 7, 8]:
            raise DecoderInvalidStringLengthException
        buffer = bytearray()
        p_byte = self._decode_first_byte(quantum[0:2])
        buffer.append(p_byte.byte)
        if len(quantum) == 2:
            return self._return_quantum(quantum, p_byte, buffer)
        p_byte = self._decode_second_byte(quantum[2:4], p_byte.carry)
        buffer.append(p_byte.byte)
        if len(quantum) == 4:
            return self._return_quantum(quantum, p_byte, buffer)
        p_byte = self._decode_third_byte(quantum[4], p_byte.carry)
        buffer.append(p_byte.byte)
        if len(quantum) == 5:
            return self._return_quantum(quantum, p_byte, buffer)
        p_byte = self._decode_fourth_byte(quantum[5:7], p_byte.carry)
        buffer.append(p_byte.byte)
        if len(quantum) == 7:
            return self._return_quantum(quantum, p_byte, buffer)
        p_byte = self._decode_fifth_byte(quantum[7], p_byte.carry)
        buffer.append(p_byte.byte)
        return buffer

    def _consume(self):
        tail = 0
        for head in range(8, len(self._string_buffer), 8):
            quantum: str = self._string_buffer[tail:head]
            self._bytearray.extend(self._decode_quantum(quantum))
            tail = head
        self._string_buffer = self._string_buffer[tail:]

    def update(self, string: str):
        if self._is_finished:
            raise DecoderAlreadyFinalizedException
        self._string_buffer += string
        self._consume()

    def _check_checksum(self, checksymbol: str) -> bytes:
        expected = self._alphabet.get(checksymbol)
        if self._checksum == expected:
            return bytes(self._bytearray)
        else:
            raise DecoderChecksumException(
                "Calculated checksum %i, expected %i" % (self._checksum, expected)
            )

    def finalize(self) -> bytes:
        if self._is_finished:
            raise DecoderAlreadyFinalizedException
        self._is_finished = True
        if self._do_checksum:
            checksymbol = self._string_buffer[-1]
            self._string_buffer = self._string_buffer[:-1]
        self._bytearray.extend(
            self._decode_quantum(self._string_buffer)
            if len(self._string_buffer) > 0
            else []
        )
        if self._do_checksum:
            return self._check_checksum(checksymbol)
        else:
            return bytes(self._bytearray)
