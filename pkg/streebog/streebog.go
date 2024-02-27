// GOST R 34.11-2012 hash function.
// RFC 6986.
package streebog

import (
	"encoding/binary"
	"hash"
)

type uint512 [8]uint64

const (
	// size512 is the size, in bytes, of a Streebog-512 checksum.
	size512 = 64

	// size256 is the size, in bytes, of a Streebog-256 checksum.
	size256 = 32

	// The blocksize of Streebog in bytes.
	blockSize = 64
)

// digest represents the partial evaluation of a checksum.
type digest struct {
	hash  uint512
	x     [blockSize]byte
	nx    int
	n     uint512
	sigma uint512
	size  int // size256 or size512
}

func (d *digest) Reset() {
	d.hash = uint512{}
	d.x = [blockSize]byte{}
	d.nx = 0
	d.n = uint512{}
	d.sigma = uint512{}

	for i := 0; i < 8; i++ {
		if d.size == size256 {
			d.hash[i] = 0x0101010101010101
		}
	}
}

// New returns a new [hash.Hash] computing the Streebog checksum.
func New() hash.Hash {
	d := &digest{size: size256}
	d.Reset()
	return d
}

func (d *digest) Size() int { return d.size }

func (d *digest) BlockSize() int { return blockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)

	if d.nx > 0 {
		chunkSize := blockSize - d.nx
		if chunkSize > len(p) {
			chunkSize = len(p)
		}
		copy(d.x[d.nx:], p[:chunkSize])

		d.nx += chunkSize
		p = p[chunkSize:]

		if d.nx == blockSize {
			var m uint512
			m[0] = binary.LittleEndian.Uint64(d.x[0:])
			m[1] = binary.LittleEndian.Uint64(d.x[8:])
			m[2] = binary.LittleEndian.Uint64(d.x[16:])
			m[3] = binary.LittleEndian.Uint64(d.x[24:])
			m[4] = binary.LittleEndian.Uint64(d.x[32:])
			m[5] = binary.LittleEndian.Uint64(d.x[40:])
			m[6] = binary.LittleEndian.Uint64(d.x[48:])
			m[7] = binary.LittleEndian.Uint64(d.x[56:])
			d.nx = 0

			g(&d.hash, d.n, m)
			add512(d.n, buffer512, &d.n)
			add512(d.sigma, m, &d.sigma)
		}
	}

	for len(p) >= blockSize {
		var m uint512
		m[0] = binary.LittleEndian.Uint64(p[0:])
		m[1] = binary.LittleEndian.Uint64(p[8:])
		m[2] = binary.LittleEndian.Uint64(p[16:])
		m[3] = binary.LittleEndian.Uint64(p[24:])
		m[4] = binary.LittleEndian.Uint64(p[32:])
		m[5] = binary.LittleEndian.Uint64(p[40:])
		m[6] = binary.LittleEndian.Uint64(p[48:])
		m[7] = binary.LittleEndian.Uint64(p[56:])
		p = p[blockSize:]

		g(&d.hash, d.n, m)
		add512(d.n, buffer512, &d.n)
		add512(d.sigma, m, &d.sigma)
	}

	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}

	return
}

func (d *digest) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	hash := d0.checkSum256()
	return append(in, hash[:]...)
}

func (d *digest) checkSum256() [size256]byte {
	var buf, m uint512

	buf[0] = uint64(d.nx << 3)

	clear(d.x[d.nx:])
	d.x[d.nx] = 1

	m[0] = binary.LittleEndian.Uint64(d.x[0:])
	m[1] = binary.LittleEndian.Uint64(d.x[8:])
	m[2] = binary.LittleEndian.Uint64(d.x[16:])
	m[3] = binary.LittleEndian.Uint64(d.x[24:])
	m[4] = binary.LittleEndian.Uint64(d.x[32:])
	m[5] = binary.LittleEndian.Uint64(d.x[40:])
	m[6] = binary.LittleEndian.Uint64(d.x[48:])
	m[7] = binary.LittleEndian.Uint64(d.x[56:])

	g(&d.hash, d.n, m)
	add512(d.n, buf, &d.n)
	add512(d.sigma, m, &d.sigma)
	g(&d.hash, uint512{}, d.n)
	g(&d.hash, uint512{}, d.sigma)

	var digest [size256]byte
	binary.LittleEndian.PutUint64(digest[0:], d.hash[4])
	binary.LittleEndian.PutUint64(digest[8:], d.hash[5])
	binary.LittleEndian.PutUint64(digest[16:], d.hash[6])
	binary.LittleEndian.PutUint64(digest[24:], d.hash[7])

	return digest
}

func (d *digest) checkSum512() [size512]byte {
	var buf, m uint512

	buf[0] = uint64(d.nx << 3)

	clear(d.x[d.nx:])
	d.x[d.nx] = 1

	m[0] = binary.LittleEndian.Uint64(d.x[0:])
	m[1] = binary.LittleEndian.Uint64(d.x[8:])
	m[2] = binary.LittleEndian.Uint64(d.x[16:])
	m[3] = binary.LittleEndian.Uint64(d.x[24:])
	m[4] = binary.LittleEndian.Uint64(d.x[32:])
	m[5] = binary.LittleEndian.Uint64(d.x[40:])
	m[6] = binary.LittleEndian.Uint64(d.x[48:])
	m[7] = binary.LittleEndian.Uint64(d.x[56:])

	g(&d.hash, d.n, m)
	add512(d.n, buf, &d.n)
	add512(d.sigma, m, &d.sigma)
	g(&d.hash, uint512{}, d.n)
	g(&d.hash, uint512{}, d.sigma)

	var digest [size512]byte
	binary.LittleEndian.PutUint64(digest[0:], d.hash[0])
	binary.LittleEndian.PutUint64(digest[8:], d.hash[1])
	binary.LittleEndian.PutUint64(digest[16:], d.hash[2])
	binary.LittleEndian.PutUint64(digest[24:], d.hash[3])
	binary.LittleEndian.PutUint64(digest[32:], d.hash[4])
	binary.LittleEndian.PutUint64(digest[40:], d.hash[5])
	binary.LittleEndian.PutUint64(digest[48:], d.hash[6])
	binary.LittleEndian.PutUint64(digest[56:], d.hash[7])

	return digest
}

// Sum returns the Streebog256 checksum of the data.
func Sum256(data []byte) [size256]byte {
	d := &digest{size: size256}
	d.Reset()
	d.Write(data)
	return d.checkSum256()
}

// Sum returns the Streebog256 checksum of the data.
func Sum512(data []byte) [size512]byte {
	d := &digest{size: size512}
	d.Reset()
	d.Write(data)
	return d.checkSum512()
}
