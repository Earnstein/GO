package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key     string    
	Created time.Time 
}

type RentensionMap map[string]OTP

func NewRentionMap(ctx context.Context, retentionPeriod time.Duration) RentensionMap {
	rm := make(RentensionMap)
	go rm.Retension(ctx, retentionPeriod)
	return rm
}

func (rm RentensionMap) NewOTP() OTP {
	token := OTP{
		Key:     uuid.NewString(),
		Created: time.Now(),
	}

	rm[token.Key] = token

	return token
}

func (rm RentensionMap) verifyOTP(otp string) bool {
	if _, ok := rm[otp]; !ok {
		return false
	}
	delete(rm, otp)
	return true
}

func (rm RentensionMap) Retension(ctx context.Context, retentionPeriod time.Duration) {
	ticker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			for _, token := range rm {
				if token.Created.Add(retentionPeriod).Before(time.Now()) {
					delete(rm, token.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
