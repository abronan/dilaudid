package dilaudid

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
const timeChars = 10
const alphabetSize = uint64(len(alphabet))

var unalphabet = map[byte]byte{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
	'G': 16,
	'H': 17,
	'J': 18,
	'K': 19,
	'M': 20,
	'N': 21,
	'P': 22,
	'Q': 23,
	'R': 24,
	'S': 25,
	'T': 26,
	'V': 27,
	'W': 28,
	'X': 29,
	'Y': 30,
	'Z': 31,
}

// A ULID is a 6-byte timestamp (encoded to 10 chars) and a 10-byte
// nonce (encoded to 16 chars). Those are the facts.
type ULID struct {
	// Timestamp
	time time.Time
	// Nonce bytes
	nonce [10]byte
}

var (
	// Nil as empty value to handle errors
	Nil ULID
)

// Unique returns a unique ULID with the current time
//
// We panic in the unlikely event that we are unable to obtain a
// nonce number.
func Unique() ULID {
	var ulid ULID
	_, err := rand.Read(ulid.nonce[:])
	if err != nil {
		panic(fmt.Errorf("dilaudid: unable to obtain a nonce number: %v", err))
	}
	ulid.time = time.Now()
	return ulid
}

// New returns a ULID with the given time and nonce
//
// Time will be standardized to UTC.
func New(t time.Time, nonce [10]byte) ULID {
	return ULID{t.UTC(), nonce}
}

// String encodes the ULID into a human-friendly string, as per the spec
func (ulid ULID) String() string {
	msecs := uint64(ulid.time.UnixNano()) / uint64(time.Millisecond)

	var builder [26]byte

	// MSB-encode milliseconds into 10 characters, each representing
	// five bits
	for i := timeChars - 1; i >= 0; i-- {
		remainder := msecs % alphabetSize
		builder[i] = alphabet[remainder]
		msecs = (msecs - remainder) / alphabetSize
	}

	for chunk := uint(0); chunk < 2; chunk++ {
		// Convert five 8-bit numbers into a 40-bit number
		x := uint64(0)
		for i := uint(0); i < 5; i++ {
			x = (x << 8) | uint64(ulid.nonce[chunk*5+i])
		}

		// Convert the 40-bit number into 8 characters, each
		// representing five bits
		for i := uint(0); i < 8; i++ {
			shift := (7 - i) * 5
			builder[timeChars+chunk*8+i] = alphabet[(x>>shift)&0x1f]
		}
	}

	return string(builder[:])
}

// Decode takes a ULID string and returns a ULID containing its
// constitutent parts (time and nonce)
//
// Returns Nil and an error if the string is malformed or the wrong
// length
func Decode(hay string) (ULID, error) {
	var ulid ULID
	err := decode(hay, &ulid.time, &ulid.nonce)
	if err != nil {
		return Nil, err
	}
	return ulid, nil
}

// MarshalJSON serializes a ULID as its string representation
func (ulid ULID) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('"')
	buf.WriteString(ulid.String())
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

// UnmarshalJSON deserializes a ULID from its string representation
//
// Time will be standardized to UTC.
func (ulid *ULID) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return fmt.Errorf("dilaudid: unable to unmarshal bytes as string: %v", err)
	}
	err = decode(str, &ulid.time, &ulid.nonce)
	if err != nil {
		return fmt.Errorf("dilaudid: unable to decode string as ULID: %v", err)
	}
	return nil
}

func parseByte(hay string, offset uint) (byte, uint, error) {
	if offset >= uint(len(hay)) {
		return 0, offset, fmt.Errorf("at offset %v unexpectedly reached the end", offset)
	}
	if byte, ok := unalphabet[hay[offset]]; ok {
		return byte, offset + 1, nil
	}
	return 0, offset, fmt.Errorf("unable to parse byte %v in %v", offset, hay)
}

func parseBytes(hay string, offset uint, bytes []byte) (uint, error) {
	j := offset
	for i := 0; i < len(bytes); i++ {
		byte, next, err := parseByte(hay, j)
		if err != nil {
			return offset, err
		}
		bytes[i] = byte
		j = next
	}
	return offset + uint(len(bytes)), nil
}

func decode(hay string, timestamp *time.Time, nonce *[10]byte) error {
	var timeBytes [10]byte
	var nonceBytes [16]byte

	hay = strings.ToUpper(hay)

	offset, err := parseBytes(hay, 0, timeBytes[:])
	if err != nil {
		return fmt.Errorf("dilaudid: parsing error at offset %v: %v", offset, err)
	}
	offset, err = parseBytes(hay, offset, nonceBytes[:])
	if err != nil {
		return fmt.Errorf("dilaudid: parsing error at offset %v: %v", offset, err)
	}
	if offset != uint(len(hay)) {
		return fmt.Errorf("dilaudid: at offset %v the string is not yet finished: %v", offset, hay)
	}

	var msecs uint64
	for _, byte := range timeBytes {
		msecs = msecs*alphabetSize + uint64(byte)
	}
	*timestamp = time.Unix(int64(msecs/1000), int64((msecs%1000)*1e6)).UTC()

	for chunk := uint(0); chunk < 2; chunk++ {
		// Convert eight 5-bit numbers into a 40-bit number
		x := uint64(0)
		for i := uint(0); i < 8; i++ {
			x = (x << 5) | uint64(nonceBytes[chunk*8+i])
		}

		// Convert the 40-bit number into five bytes
		for i := uint(0); i < 5; i++ {
			shift := (4 - i) * 8
			nonce[chunk*5+i] = byte((x >> shift) & 0xff)
		}
	}

	return nil
}
