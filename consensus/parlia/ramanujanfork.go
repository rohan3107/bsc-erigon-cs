package parlia

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ledgerwatch/erigon/consensus"
	"github.com/ledgerwatch/erigon/core/types"
)

const (
	wiggleTimeBeforeFork       = 500 * time.Millisecond // Random delay (per signer) to allow concurrent signers
	fixedBackOffTimeBeforeFork = 200 * time.Millisecond
)

func (p *Parlia) delayForRamanujanFork(snap *Snapshot, header *types.Header) time.Duration {
	delay := time.Until(time.Unix(int64(header.Time), 0)) // nolint: gosimple
	if p.chainConfig.IsRamanujan(header.Number.Uint64()) {
		return delay
	}
	if header.Difficulty.Cmp(diffNoTurn) == 0 {
		// It's not our turn explicitly to sign, delay it a bit
		wiggle := time.Duration(len(snap.Validators)/2+1) * wiggleTimeBeforeFork
		delay += fixedBackOffTimeBeforeFork + cryptoRandDuration(wiggle)
	}
	return delay
}

func cryptoRandDuration(max time.Duration) time.Duration {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		// In case of an error with the crypto/rand, fallback to a default delay
		return fixedBackOffTimeBeforeFork
	}
	return time.Duration(nBig.Int64())
}

func (p *Parlia) blockTimeForRamanujanFork(snap *Snapshot, header, parent *types.Header) uint64 {
	blockTime := parent.Time + p.config.Period
	if p.chainConfig.IsRamanujan(header.Number.Uint64()) {
		blockTime = blockTime + backOffTime(snap, header, p.val, p.chainConfig)
	}
	return blockTime
}

func (p *Parlia) blockTimeVerifyForRamanujanFork(snap *Snapshot, header, parent *types.Header) error {
	if p.chainConfig.IsRamanujan(header.Number.Uint64()) {
		if header.Time < parent.Time+p.config.Period+backOffTime(snap, header, header.Coinbase, p.chainConfig) {
			return fmt.Errorf("header %d, time %d, now %d, period: %d, backof: %d, %w", header.Number.Uint64(), header.Time, time.Now().Unix(), p.config.Period, backOffTime(snap, header, header.Coinbase, p.chainConfig), consensus.ErrFutureBlock)
		}
	}
	return nil
}
