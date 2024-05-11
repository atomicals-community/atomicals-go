package sha

type digest struct {
	h     [8]uint32
	x     [chunk]byte
	nx    int
	len   uint64
	is224 bool // mark if this digest is SHA-224
}

// SHA256 computes the SHA256 hash of the input data
func SHA256(data []byte) []byte {
	d := &digest{}
	d.Reset()

	d.Write(data)

	da := d.checkSum()

	return da[:]
}

func (d *digest) Reset() {
	d.h[0] = init0
	d.h[1] = init1
	d.h[2] = init2
	d.h[3] = init3
	d.h[4] = init4
	d.h[5] = init5
	d.h[6] = init6
	d.h[7] = init7

	d.nx = 0
	d.len = 0
}

func (d *digest) Write(p []byte) {
	nn := len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == chunk {
			d.block(d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		d.block(p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
}

func (dig *digest) block(p []byte) {
	var w [64]uint32
	h0, h1, h2, h3, h4, h5, h6, h7 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7]
	for len(p) >= chunk {
		// Can interlace the computation of w with the
		// rounds below if needed for speed.
		for i := 0; i < 16; i++ {
			j := i * 4
			w[i] = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])
		}
		for i := 16; i < 64; i++ {
			v1 := w[i-2]
			t1 := (RotateLeft32(v1, -17)) ^ (RotateLeft32(v1, -19)) ^ (v1 >> 10)
			v2 := w[i-15]
			t2 := (RotateLeft32(v2, -7)) ^ (RotateLeft32(v2, -18)) ^ (v2 >> 3)
			w[i] = t1 + w[i-7] + t2 + w[i-16]
		}

		a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7

		for i := 0; i < 64; i++ {
			t1 := h + ((RotateLeft32(e, -6)) ^ (RotateLeft32(e, -11)) ^ (RotateLeft32(e, -25))) + ((e & f) ^ (^e & g)) + _K[i] + w[i]
			t2 := ((RotateLeft32(a, -2)) ^ (RotateLeft32(a, -13)) ^ (RotateLeft32(a, -22))) + ((a & b) ^ (a & c) ^ (b & c))

			h = g
			g = f
			f = e
			e = d + t1
			d = c
			c = b
			b = a
			a = t1 + t2
		}

		h0 += a
		h1 += b
		h2 += c
		h3 += d
		h4 += e
		h5 += f
		h6 += g
		h7 += h

		p = p[chunk:]
	}

	dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7] = h0, h1, h2, h3, h4, h5, h6, h7
}
