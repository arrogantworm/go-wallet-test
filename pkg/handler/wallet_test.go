package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"wallet-app/pkg/service"
	"wallet-app/pkg/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetWallet(t *testing.T) {
	repo := testutils.SetupTestPG(t)
	s := service.NewService(repo)
	h := NewHandler(s)
	ctx := context.Background()
	walletID := uuid.New()

	t.Run("wallet not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
		rr := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/api/v1/wallets/{id}", h.getWalletInfo)

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("valid get", func(t *testing.T) {
		walletBalance := 100

		err := repo.NewWallet(ctx, walletID, walletBalance)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
		rr := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/api/v1/wallets/{id}", h.getWalletInfo)

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp WalletResp
		err = json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)

		balanceStr := resp.Balance
		balanceStr = strings.Replace(balanceStr, ".", "", 1)
		balance, _ := strconv.Atoi(balanceStr)

		assert.Equal(t, walletBalance, balance)

	})

}

func TestDepositAndWithdraw(t *testing.T) {
	repo := testutils.SetupTestPG(t)
	s := service.NewService(repo)
	h := NewHandler(s)
	ctx := context.Background()
	walletID := uuid.New()

	_ = repo.NewWallet(ctx, walletID, 100)

	router := chi.NewRouter()
	router.Post("/api/v1/wallet", h.updateWalletBalance)

	t.Run("valid deposit", func(t *testing.T) {
		body := `{"walletId":"` + walletID.String() + `", "operationType": "deposit", "amount": "50.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("invalid deposit - wrong amount", func(t *testing.T) {
		body := `{"walletId":"` + walletID.String() + `", "operationType": "deposit", "amount": "5000"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("negative withdraw", func(t *testing.T) {
		body := `{"walletId":"` + walletID.String() + `", "operationType": "withdraw", "amount": "-30.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("valid withdraw", func(t *testing.T) {
		body := `{"walletId":"` + walletID.String() + `", "operationType": "withdraw", "amount": "30.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("withdraw with not enough funds", func(t *testing.T) {
		body := `{"walletId":"` + walletID.String() + `", "operationType": "withdraw", "amount": "1000.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid operation type", func(t *testing.T) {
		body := `{"walletId":"` + walletID.String() + `", "operationType": "test", "amount": "10.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		body := `{"walletId":"test", "operationType": "deposit", "amount": "10.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("create wallet on first deposit", func(t *testing.T) {
		newWalletID := uuid.New()
		body := `{"walletId":"` + newWalletID.String() + `", "operationType": "deposit", "amount": "25.00"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		balance, err := repo.GetBalance(ctx, newWalletID)
		assert.NoError(t, err)
		assert.Equal(t, 25*100, balance)
	})

}

func TestConcurrentDeposits(t *testing.T) {
	repo := testutils.SetupTestPG(t)
	s := service.NewService(repo)
	h := NewHandler(s)
	ctx := context.Background()
	walletID := uuid.New()

	_ = repo.NewWallet(ctx, walletID, 0)

	router := chi.NewRouter()
	router.Post("/api/v1/wallet", h.updateWalletBalance)

	server := httptest.NewServer(router)
	defer server.Close()

	var wg sync.WaitGroup
	const totalRequests = 1000
	const concurrency = 1000

	sem := make(chan struct{}, concurrency)

	var successCount int64
	var failCount int64

	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{Transport: transport}

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			body := fmt.Sprintf(`{"walletId":"%s", "operationType": "deposit", "amount": "1.00"}`, walletID)
			req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/wallet", strings.NewReader(body))
			if err != nil {
				atomic.AddInt64(&failCount, 1)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				atomic.AddInt64(&failCount, 1)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}
		}()
	}

	wg.Wait()

	finalBalance, err := repo.GetBalance(ctx, walletID)
	// log.Println(finalBalance)
	assert.NoError(t, err)

	assert.Equal(t, int(totalRequests)*100, finalBalance)
	assert.Equal(t, int64(totalRequests), successCount)
	assert.Equal(t, int64(0), failCount)
}
