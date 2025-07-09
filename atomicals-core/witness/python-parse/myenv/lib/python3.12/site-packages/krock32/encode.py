# encode.py

"""
Krock32 encoder

Class definition for the krock32 encoder, taking bytes-data and turning it
into Crockford's Base32 string.
"""

from collections import namedtuple
from typing import List


class EncoderAlreadyFinalizedException(Exception):
    pass


class Encoder:
    def __init__(self, checksum: bool = False):
        self._string_list: List[str] = []
        self._byte_buffer: bytearray = bytearray()
        self._alphabet: dict = self._make_alphabet(
            "0123456789ABCDEFGHJKMNPQRSTVWXYZ*~$=U"
        )
        self._is_finished: bool = False
        self._checksum = 0
        self._do_checksum: bool = checksum
        self._p_quin = namedtuple("ProcessedQuin", ["sym", "rem"])

    def _make_alphabet(self, alphabet_string: str) -> dict:
        alphabet = {}
        for i, x in enumerate(alphabet_string):
            alphabet[i] = x
        return alphabet

    def _update_checksum(self, byte: int):
        if not self._do_checksum:
            return
        self._checksum = ((self._checksum << 8) + byte) % 37

    def _encode_first_quin(self, byte) -> tuple:
        self._update_checksum(byte)
        quin = byte >> 3
        rem = (byte & 0b111) << 2
        return self._p_quin(sym=self._alphabet.get(quin), rem=rem)

    def _encode_second_quin(self, byte, remainder) -> tuple:
        self._update_checksum(byte)
        sym = ""
        quin = (byte >> 6) + remainder
        sym += self._alphabet.get(quin)
        quin = (byte >> 1) & 0b11111
        rem = (byte & 0b1) << 4
        sym += self._alphabet.get(quin)
        return self._p_quin(sym=sym, rem=rem)

    def _encode_third_quin(self, byte, remainder) -> tuple:
        self._update_checksum(byte)
        quin = (byte >> 4) + remainder
        rem = (byte & 0b1111) << 1
        return self._p_quin(sym=self._alphabet.get(quin), rem=rem)

    def _encode_fourth_quin(self, byte, remainder) -> tuple:
        self._update_checksum(byte)
        sym = ""
        quin = (byte >> 7) + remainder
        sym += self._alphabet.get(quin)
        quin = (byte >> 2) & 0b11111
        sym += self._alphabet.get(quin)
        rem = (byte & 0b11) << 3
        return self._p_quin(sym=sym, rem=rem)

    def _encode_fifth_quin(self, byte, remainder) -> tuple:
        self._update_checksum(byte)
        sym = ""
        quin = (byte >> 5) + remainder
        sym += self._alphabet.get(quin)
        quin = byte & 0b11111
        sym += self._alphabet.get(quin)
        return self._p_quin(sym=sym, rem=0)

    def _encode_quantum(self, quantum: bytearray) -> List[str]:
        slist: List[str] = []
        p_quin = self._encode_first_quin(quantum[0])
        slist.append(p_quin.sym)
        if len(quantum) == 1:
            slist.append(self._alphabet.get(p_quin.rem))
            return slist
        p_quin = self._encode_second_quin(quantum[1], p_quin.rem)
        slist.append(p_quin.sym)
        if len(quantum) == 2:
            slist.append(self._alphabet.get(p_quin.rem))
            return slist
        p_quin = self._encode_third_quin(quantum[2], p_quin.rem)
        slist.append(p_quin.sym)
        if len(quantum) == 3:
            slist.append(self._alphabet.get(p_quin.rem))
            return slist
        p_quin = self._encode_fourth_quin(quantum[3], p_quin.rem)
        slist.append(p_quin.sym)
        if len(quantum) == 4:
            slist.append(self._alphabet.get(p_quin.rem))
            return slist
        p_quin = self._encode_fifth_quin(quantum[4], p_quin.rem)
        slist.append(p_quin.sym)
        return slist

    def _consume(self):
        tail = 0
        for head in range(5, len(self._byte_buffer), 5):
            quantum = self._byte_buffer[tail:head]
            self._string_list.extend(self._encode_quantum(quantum))
            tail = head
        self._byte_buffer = self._byte_buffer[tail:]

    def update(self, data: bytes):
        if self._is_finished:
            raise EncoderAlreadyFinalizedException
        self._byte_buffer.extend(bytearray(data))
        self._consume()

    def finalize(self) -> str:
        if self._is_finished:
            raise EncoderAlreadyFinalizedException
        self._is_finished = True
        self._string_list.extend(
            self._encode_quantum(self._byte_buffer)
            if len(self._byte_buffer) > 0
            else []
        )
        if self._do_checksum:
            self._string_list.append(self._alphabet.get(self._checksum))
        return "".join(self._string_list)
