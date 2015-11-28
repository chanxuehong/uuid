package v1

import (
	"sync"

	"github.com/chanxuehong/internal"
	"github.com/chanxuehong/rand"
)

const sequenceMask uint32 = 0x3FFF // 14bits

var (
	node = internal.MAC[:]

	mutex         sync.Mutex
	sequenceStart uint32 = rand.Uint32() & sequenceMask
	lastTimestamp int64  = -1
	lastSequence  uint32 = sequenceStart
)

func New() (uuid [16]byte) {
	var (
		timestamp = uuidTimestamp()
		sequence  = sequenceStart
	)

	mutex.Lock() // Lock
	switch {
	case timestamp > lastTimestamp:
		lastTimestamp = timestamp
		lastSequence = sequence
		mutex.Unlock() // Unlock
	case timestamp == lastTimestamp:
		sequence = (lastSequence + 1) & sequenceMask
		if sequence == sequenceStart {
			timestamp = tillNext100nano(timestamp)
			lastTimestamp = timestamp
		}
		lastSequence = sequence
		mutex.Unlock() // Unlock
	default: // timestamp < lastTimestamp
		sequenceStart = rand.Uint32() & sequenceMask // NOTE
		sequence = sequenceStart
		lastTimestamp = timestamp
		lastSequence = sequence
		mutex.Unlock() // Unlock
	}

	// time_low
	uuid[0] = byte(timestamp >> 24)
	uuid[1] = byte(timestamp >> 16)
	uuid[2] = byte(timestamp >> 8)
	uuid[3] = byte(timestamp)

	// time_mid
	uuid[4] = byte(timestamp >> 40)
	uuid[5] = byte(timestamp >> 32)

	// time_hi_and_version
	uuid[6] = byte(timestamp>>56) & 0x0F
	uuid[7] = byte(timestamp >> 48)
	uuid[6] |= 0x10 // version, 4bits

	// clk_seq_hi_res
	uuid[8] = byte(sequence>>8) & 0x3F
	uuid[8] |= 0x80 // variant

	// clk_seq_low
	uuid[9] = byte(sequence)

	// node
	copy(uuid[10:], node)
	return
}
