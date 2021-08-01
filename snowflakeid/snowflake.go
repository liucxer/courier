package snowflakeid

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	InvalidSystemClock = errors.New("invalid system clock")
)

func NewSnowflakeFactory(bitLenWorkerID, bitLenSequence, gapMs uint, startTime time.Time) *SnowflakeFactory {
	return &SnowflakeFactory{
		bitLenWorkerID:  bitLenWorkerID,
		bitLenSequence:  bitLenSequence,
		bitLenTimestamp: 63 - bitLenSequence - bitLenWorkerID,
		startTime:       startTime,
		unit:            time.Duration(gapMs) * time.Millisecond,
	}
}

type SnowflakeFactory struct {
	bitLenWorkerID, bitLenTimestamp, bitLenSequence uint
	startTime                                       time.Time
	unit                                            time.Duration
}

func (f *SnowflakeFactory) MaskSequence(sequence uint32) uint32 {
	return sequence & f.MaxSequence()
}

func (f *SnowflakeFactory) FlakeTimestamp(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(f.unit))
}

func (f *SnowflakeFactory) CurrentElapsedTime() uint64 {
	return f.FlakeTimestamp(time.Now()) - f.FlakeTimestamp(f.startTime)
}

func (f *SnowflakeFactory) SleepTime(overtime time.Duration) time.Duration {
	return overtime*f.unit - time.Duration(time.Now().UnixNano())%f.unit*time.Nanosecond
}

func (f *SnowflakeFactory) BuildID(workerID uint32, elapsedTime uint64, sequence uint32) (uint64, error) {
	if elapsedTime >= 1<<f.bitLenTimestamp {
		return 0, errors.New("over the time limit")
	}
	return elapsedTime<<(f.bitLenSequence+f.bitLenWorkerID) | uint64(sequence)<<f.bitLenWorkerID | uint64(workerID), nil
}

func (f *SnowflakeFactory) MaxSequence() uint32 {
	return 1<<f.bitLenSequence - 1
}

func (f *SnowflakeFactory) MaxWorkerID() uint32 {
	return 1<<f.bitLenWorkerID - 1
}

func (f *SnowflakeFactory) MaxTime() time.Time {
	maxTime := uint64(1<<f.bitLenTimestamp - 1)
	return time.Unix(int64(time.Duration(f.FlakeTimestamp(f.startTime)+maxTime)*f.unit/time.Second), 0)
}

func (f *SnowflakeFactory) NewSnowflake(workerID uint32) (*Snowflake, error) {
	maxWorkerID := f.MaxWorkerID()
	if workerID > maxWorkerID {
		return nil, fmt.Errorf("worker id can't be large than %d", maxWorkerID)
	}
	return &Snowflake{f: f, workerID: workerID, syncMutex: &sync.Mutex{}}, nil
}

func NewSnowflake(workerID uint32) (*Snowflake, error) {
	startTime, _ := time.Parse(time.RFC3339, "2010-11-04T01:42:54.657Z")
	return NewSnowflakeFactory(10, 12, 1, startTime).NewSnowflake(workerID)
}

type Snowflake struct {
	f           *SnowflakeFactory
	workerID    uint32
	elapsedTime uint64
	sequence    uint32
	syncMutex   *sync.Mutex
}

func (sf *Snowflake) WorkerID() uint32 {
	return sf.workerID
}

func (sf *Snowflake) ID() (uint64, error) {
	sf.syncMutex.Lock()
	defer sf.syncMutex.Unlock()

	currentElapsedTime := sf.f.CurrentElapsedTime()

	if sf.elapsedTime < currentElapsedTime {
		sf.elapsedTime = currentElapsedTime
		sf.sequence = generateRandomSequence(9)

		return sf.f.BuildID(sf.workerID, sf.elapsedTime, sf.sequence)
	}

	if sf.elapsedTime > currentElapsedTime {
		currentElapsedTime = sf.f.CurrentElapsedTime()
		if sf.elapsedTime > currentElapsedTime {
			return 0, InvalidSystemClock
		}
	}

	// ==

	sf.sequence = sf.f.MaskSequence(sf.sequence + 1)
	if sf.sequence == 0 {
		sf.elapsedTime = sf.elapsedTime + 1
		time.Sleep(sf.f.SleepTime(time.Duration(sf.elapsedTime - currentElapsedTime)))
	}

	return sf.f.BuildID(sf.workerID, sf.elapsedTime, sf.sequence)
}

func generateRandomSequence(n int32) uint32 {
	return uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(n))
}
