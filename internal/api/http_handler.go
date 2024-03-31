package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"transaction-parser/internal/data"
	"transaction-parser/internal/service"
)

type HTTPHandler struct {
	parser service.Parser
}

func NewHTTPHandler(parser service.Parser) *HTTPHandler {
	return &HTTPHandler{parser: parser}
}

func (h *HTTPHandler) HandleGetCurrentBlock(w http.ResponseWriter, _ *http.Request) {
	block := h.parser.GetCurrentBlock()
	responseJSON(w, http.StatusOK, map[string]int{"block": block})
}

func (h *HTTPHandler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 20) // 1 MB
	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to parse the form"})
		return

	}
	address := r.FormValue("address")
	if address == "" {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "address is required"})
		return
	}
	var ok bool
	address, ok = validateAndFormatEthereumAddress(address)
	if !ok {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid address"})
		return
	}
	err = h.parser.Subscribe(address)
	if err != nil {
		if errors.Is(err, data.ErrAddressAlreadySubscribed) {
			responseJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, data.ErrCapacityExceeded) {
			responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		log.Printf("failed to subscribe: %v", err)
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})

	}
	responseJSON(w, http.StatusOK, map[string]bool{"subscribed": true})
}

func (h *HTTPHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "address is required"})
		return
	}
	var ok bool
	address, ok = validateAndFormatEthereumAddress(address)
	if !ok {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid address"})
		return
	}
	transactions, err := h.parser.GetTransactions(address)
	if err != nil {
		if errors.Is(err, data.ErrAddressNotSubscribed) {
			responseJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		log.Printf("failed to get transactions: %v", err)
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	responseJSON(w, http.StatusOK, transactions)
}

func responseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to serialize response: %v", err)
		return
	}
	if _, err = w.Write(body); err != nil {
		log.Printf("failed to write the body: %v", err)
	}
}
