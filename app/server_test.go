package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubCreditCardStore struct {
	validation map[string]bool
}

func (s *StubCreditCardStore) GetCardValidation(creditCardNumber string) bool {
	isValid := s.validation[creditCardNumber]
	return isValid
}

func TestCreditCardValidator(t *testing.T) {
	store := StubCreditCardStore{
		map[string]bool{
			"3379 5135 6110 8795": true,
			"3379 5135 6110 8794": false,
		},
	}
	server := &CreditCardValidatorServer{&store}
	// 3379 5135 6110 8795
	// 2769 1483 0405 9987
	t.Run("validating valid numbers for Luhn algorithm", func(t *testing.T) {
		jsonPayload := []byte(`{"CreditCardNumber": "3379 5135 6110 8795"}`)
		request := newGetValidationRequest(jsonPayload)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "true")
	})
	// 3379 5135 6110 8794
	// 2769 1483 0405 9986
	t.Run("returns 400 for invalid numbers for Luhn algorithm", func(t *testing.T) {
		jsonPayload := []byte(`{"CreditCardNumber": "3379 5135 6110 8794"}`)
		request := newGetValidationRequest(jsonPayload)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertResponseBody(t, response.Body.String(), "false")
	})
	t.Run("returns 400 bad request on malformed input", func(t *testing.T) {
		jsonPayload := []byte(`{"CreditCardNumber": "abcdef"}`)
		request := newGetValidationRequest(jsonPayload)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetValidationRequest(creditCardPayload []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/credit_card_number/", bytes.NewBuffer(creditCardPayload))
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
