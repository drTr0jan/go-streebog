package streebog

func xlps(x uint512, y uint512, data *uint512) {
	var r uint512

	for i := 0; i < 8; i++ {
		r[i] = x[i] ^ y[i]
	}

	for i := 0; i < 8; i++ {
		(*data)[i] = ax[0][(r[0]>>(i<<3))&0xff]
		(*data)[i] ^= ax[1][(r[1]>>(i<<3))&0xff]
		(*data)[i] ^= ax[2][(r[2]>>(i<<3))&0xff]
		(*data)[i] ^= ax[3][(r[3]>>(i<<3))&0xff]
		(*data)[i] ^= ax[4][(r[4]>>(i<<3))&0xff]
		(*data)[i] ^= ax[5][(r[5]>>(i<<3))&0xff]
		(*data)[i] ^= ax[6][(r[6]>>(i<<3))&0xff]
		(*data)[i] ^= ax[7][(r[7]>>(i<<3))&0xff]
	}
}

func round(i int, ki *uint512, data *uint512) {
	xlps(*ki, c[i], ki)
	xlps(*ki, *data, data)
}

func x(x uint512, y uint512, z *uint512) {
	for i := 0; i < 8; i++ {
		(*z)[i] = x[i] ^ y[i]
	}
}

func g(h *uint512, n uint512, m uint512) {
	var data, ki uint512

	xlps(*h, n, &data)

	/* Starting E() */
	ki = data
	xlps(ki, m, &data)

	for i := 0; i < 11; i++ {
		round(i, &ki, &data)
	}

	xlps(ki, c[11], &ki)
	x(ki, data, &data)
	/* E() done */

	x(data, *h, &data)
	x(data, m, h)
}

func add512(x uint512, y uint512, r *uint512) {
	var cf int

	for i := 0; i < 8; i++ {
		left := x[i]
		sum := left + y[i] + uint64(cf)
		if sum < left {
			cf = 1
		} else if sum > left {
			cf = 0
		}
		(*r)[i] = sum
	}
}
