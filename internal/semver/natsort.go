package semver

// NatCompare ports PHP's strnatcmp() (ext/standard/strnatcmp.c,
// strnatcmp_ex without case folding) as used by sort(..., SORT_NATURAL).
// Characters past the end read as NUL like the zend string terminator.
func NatCompare(a, b string) int {
	if len(a) == 0 || len(b) == 0 {
		switch {
		case len(a) == len(b):
			return 0
		case len(a) > len(b):
			return 1
		default:
			return -1
		}
	}
	at := func(s string, i int) byte {
		if i < len(s) {
			return s[i]
		}
		return 0
	}
	ap, bp := 0, 0
	leading := true
	for {
		ca, cb := at(a, ap), at(b, bp)
		// skip over leading zeros (first iteration only)
		for leading && ca == '0' && ap+1 < len(a) && isDigitByte(a[ap+1]) {
			ap++
			ca = at(a, ap)
		}
		for leading && cb == '0' && bp+1 < len(b) && isDigitByte(b[bp+1]) {
			bp++
			cb = at(b, bp)
		}
		leading = false
		// skip consecutive whitespace
		for isSpaceByte(ca) {
			ap++
			ca = at(a, ap)
		}
		for isSpaceByte(cb) {
			bp++
			cb = at(b, bp)
		}
		if isDigitByte(ca) && isDigitByte(cb) {
			var result int
			if ca == '0' || cb == '0' {
				result = natCompareLeft(a, &ap, b, &bp)
			} else {
				result = natCompareRight(a, &ap, b, &bp)
			}
			switch {
			case result != 0:
				return result
			case ap == len(a) && bp == len(b):
				return 0
			case ap == len(a):
				return -1
			case bp == len(b):
				return 1
			}
			// keep on comparing from the current point (no whitespace skip)
			ca, cb = at(a, ap), at(b, bp)
		}
		if ca < cb {
			return -1
		}
		if ca > cb {
			return 1
		}
		ap++
		bp++
		switch {
		case ap >= len(a) && bp >= len(b):
			return 0
		case ap >= len(a):
			return -1
		case bp >= len(b):
			return 1
		}
	}
}

func isDigitByte(c byte) bool { return c >= '0' && c <= '9' }

func isSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// natCompareRight compares right-aligned numbers: the longest run wins,
// otherwise the first difference (bias).
func natCompareRight(a string, ap *int, b string, bp *int) int {
	bias := 0
	for ; ; *ap, *bp = *ap+1, *bp+1 {
		ea := *ap >= len(a) || !isDigitByte(a[*ap])
		eb := *bp >= len(b) || !isDigitByte(b[*bp])
		switch {
		case ea && eb:
			return bias
		case ea:
			return -1
		case eb:
			return 1
		case a[*ap] < b[*bp]:
			if bias == 0 {
				bias = -1
			}
		case a[*ap] > b[*bp]:
			if bias == 0 {
				bias = 1
			}
		}
	}
}

// natCompareLeft compares left-aligned (fractional) numbers.
func natCompareLeft(a string, ap *int, b string, bp *int) int {
	for ; ; *ap, *bp = *ap+1, *bp+1 {
		ea := *ap >= len(a) || !isDigitByte(a[*ap])
		eb := *bp >= len(b) || !isDigitByte(b[*bp])
		switch {
		case ea && eb:
			return 0
		case ea:
			return -1
		case eb:
			return 1
		case a[*ap] < b[*bp]:
			return -1
		case a[*ap] > b[*bp]:
			return 1
		}
	}
}
