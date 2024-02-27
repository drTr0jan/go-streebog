package streebog

import (
	"fmt"
	"io"
	"testing"
)

type streebogTest struct {
	out string
	in  string
}

var golden256 = []streebogTest{
	{"9d151eefd8590b89daa6ba6cb74af9275dd051026bb149a452fd84e5e57b5500",
		"012345678901234567890123456789012345678901234567890123456789012"},
	{"9dd2fe4e90409e5da87f53976d7405b0c0cac628fc669a741d50063c557e8f50",
		"\xd1\xe5\x20\xe2\xe5\xf2\xf0\xe8\x2c\x20\xd1\xf2\xf0\xe8\xe1\xee\xe6\xe8\x20\xe2\xed\xf3\xf6\xe8\x2c\x20\xe2\xe5\xfe\xf2\xfa\x20\xf1\x20\xec\xee\xf0\xff\x20\xf1\xf2\xf0\xe5\xeb\xe0\xec\xe8\x20\xed\xe0\x20\xf5\xf0\xe0\xe1\xf0\xfb\xff\x20\xef\xeb\xfa\xea\xfb\x20\xc8\xe3\xee\xf0\xe5\xe2\xfb"},
	{"3f539a213e97c802cc229d474c6aa32a825a360b2a933a949fd925208d9ce1bb",
		""},
	{"df1fda9ce83191390537358031db2ecaa6aa54cd0eda241dc107105e13636b95",
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"},
}

func TestGolden(t *testing.T) {
	for i := 0; i < len(golden256); i++ {
		g := golden256[i]
		s := fmt.Sprintf("%x", Sum256([]byte(g.in)))
		if s != g.out {
			t.Fatalf("Sum function: streebog(%s) = %s want %s", g.in, s, g.out)
		}
		c := New()
		buf := make([]byte, len(g.in)+4)
		for j := 0; j < 3+4; j++ {
			if j < 2 {
				io.WriteString(c, g.in)
			} else if j == 2 {
				io.WriteString(c, g.in[:len(g.in)/2])
				c.Sum(nil)
				io.WriteString(c, g.in[len(g.in)/2:])
			} else if j > 2 {
				// test unaligned write
				buf = buf[1:]
				copy(buf, g.in)
				c.Write(buf[:len(g.in)])
			}
			s := fmt.Sprintf("%x", c.Sum(nil))
			if s != g.out {
				t.Fatalf("streebog[%d](%s) = %s want %s", j, g.in, s, g.out)
			}
			c.Reset()
		}
	}
}
