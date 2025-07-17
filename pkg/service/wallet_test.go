package service

import (
	"fmt"
	"sync"
	"testing"
)

type fakeRepo struct {
	mu      sync.Mutex
	balance map[int]float64
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{balance: make(map[int]float64)}
}

func (f *fakeRepo) Deposit(walletID int, amount float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.balance[walletID] += amount
	return nil
}

func (f *fakeRepo) Withdraw(walletID int, amount float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.balance[walletID] < amount {
		return fmt.Errorf("insufficient funds")
	}
	f.balance[walletID] -= amount
	return nil
}

func (f *fakeRepo) GetBalance(walletID int) float64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.balance[walletID]
}

func TestConcurrentDeposits(t *testing.T) {
	repo := newFakeRepo()
	const walletID = 1
	const depositAmount = 1.0
	const goroutines = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = repo.Deposit(walletID, depositAmount)
		}()
	}

	wg.Wait()

	finalBalance := repo.GetBalance(walletID)
	expected := depositAmount * goroutines

	if finalBalance != expected {
		t.Errorf("expected balance %v, got %v", expected, finalBalance)
	}
}
