package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type Transaction struct {
	ID                  int    `json:"id"`
	Amount              int    `json:"amount"`
	MessageType         string `json:"conversation_type"`
	CreatedAt           string `json:"created_at"`
	TransactionID       int    `json:"transaction_id"`
	PAN                 int64  `json:"pan"` // Changed to int64 for full PANs
	TransactionCategory string `json:"transaction_category"`
	PostedTimeStamp     string `json:"posted_timestamp"`
	TransactionType     string `json:"transaction_type"`
	SendingAccount      int    `json:"sending_account"`
	ReceivingAccount    int    `json:"receiving_account"`
	TransactionNote     string `json:"transaction_note"`
}

var transactions []Transaction

func main() {
	source := flag.String("transactions", "", "path to JSON file or HTTP URL")
	flag.Parse()

	var err error

	transactions, err = LoadTransactions(*source)
	if err != nil {
		fmt.Printf("Error loading transactions: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/transactions", GetTransactions)
	http.HandleFunc("/transactions/sorted", GetSortedTransactions)
	fmt.Println("Listening on :8000")
	http.ListenAndServe(":8000", nil)
}

func LoadTransactions(source string) ([]Transaction, error) {
	var data []byte
	var err error

	if source[:4] == "http" {
		resp, err := http.Get(source)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
	} else {
		data, err = os.ReadFile(source)
	}

	if err != nil {
		return nil, err
	}

	var txns []Transaction
	err = json.Unmarshal(data, &txns)
	return txns, err
}

func maskPAN(pan int64) string {
	s := strconv.FormatInt(pan, 10)
	if len(s) < 4 {
		return "****"
	}
	return fmt.Sprintf("**** **** **** %s", s[len(s)-4:])
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	out := make([]map[string]interface{}, len(transactions))
	for i, t := range transactions {
		txnMap := map[string]interface{}{
			"id":                   t.ID,
			"amount":               t.Amount,
			"conversation_type":    t.MessageType,
			"created_at":           t.CreatedAt,
			"transaction_id":       t.TransactionID,
			"pan":                  maskPAN(t.PAN),
			"transaction_category": t.TransactionCategory,
			"posted_timestamp":     t.PostedTimeStamp,
			"transaction_type":     t.TransactionType,
			"sending_account":      t.SendingAccount,
			"receiving_account":    t.ReceivingAccount,
			"transaction_note":     t.TransactionNote,
		}
		out[i] = txnMap
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func GetSortedTransactions(w http.ResponseWriter, r *http.Request) {
	sorted := SortTransactionsDesc(transactions)

	out := make([]map[string]interface{}, len(sorted))
	for i, t := range sorted {
		txnMap := map[string]interface{}{
			"id":                   t.ID,
			"amount":               t.Amount,
			"conversation_type":    t.MessageType,
			"created_at":           t.CreatedAt,
			"transaction_id":       t.TransactionID,
			"pan":                  maskPAN(t.PAN),
			"transaction_category": t.TransactionCategory,
			"posted_timestamp":     t.PostedTimeStamp,
			"transaction_type":     t.TransactionType,
			"sending_account":      t.SendingAccount,
			"receiving_account":    t.ReceivingAccount,
			"transaction_note":     t.TransactionNote,
		}
		out[i] = txnMap
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func SortTransactionsDesc(txns []Transaction) []Transaction {
	sorted := make([]Transaction, len(txns))
	copy(sorted, txns)
	sort.Slice(sorted, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, sorted[i].PostedTimeStamp)
		t2, _ := time.Parse(time.RFC3339, sorted[j].PostedTimeStamp)
		return t1.After(t2)
	})
	return sorted
}

func DefaultMockTransactions() []Transaction {
	return []Transaction{
		{
			ID:                  1,
			Amount:              2750,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-01T09:21:15+00:00",
			TransactionID:       110,
			PAN:                 4080230386144446,
			TransactionCategory: "Grocery",
			PostedTimeStamp:     "2025-06-01T09:21:15+00:00",
			TransactionType:     "POS",
			SendingAccount:      50012,
			ReceivingAccount:    110902,
			TransactionNote:     "Whole Foods Market",
		},
		{
			ID:                  2,
			Amount:              8499,
			MessageType:         "Credit",
			CreatedAt:           "2025-06-02T14:17:42+00:00",
			TransactionID:       112,
			PAN:                 5166697943434128,
			TransactionCategory: "Travel",
			PostedTimeStamp:     "2025-06-02T14:17:42+00:00",
			TransactionType:     "Refund",
			SendingAccount:      59400,
			ReceivingAccount:    1234500,
			TransactionNote:     "Air Canada Refund",
		},
		{
			ID:                  3,
			Amount:              50000,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-05T16:48:10+00:00",
			TransactionID:       117,
			PAN:                 5488452462266852,
			TransactionCategory: "ATM",
			PostedTimeStamp:     "2025-06-05T16:48:10+00:00",
			TransactionType:     "Withdrawal",
			SendingAccount:      77302,
			ReceivingAccount:    21809,
			TransactionNote:     "ATM #48211 Granville Street",
		},
		{
			ID:                  4,
			Amount:              12059,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-07T11:15:05+00:00",
			TransactionID:       118,
			PAN:                 4954335252282726,
			TransactionCategory: "Automotive",
			PostedTimeStamp:     "2025-06-07T11:15:05+00:00",
			TransactionType:     "POS",
			SendingAccount:      93839,
			ReceivingAccount:    9233020,
			TransactionNote:     "Quick Lube Service",
		},
		{
			ID:                  5,
			Amount:              7999,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-10T13:20:00+00:00",
			TransactionID:       121,
			PAN:                 4844085301308048,
			TransactionCategory: "Household",
			PostedTimeStamp:     "2025-06-10T13:20:00+00:00",
			TransactionType:     "POS",
			SendingAccount:      10018,
			ReceivingAccount:    9222020,
			TransactionNote:     "IKEA Furniture",
		},
		{
			ID:                  6,
			Amount:              14200,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-12T18:05:30+00:00",
			TransactionID:       133,
			PAN:                 4090070794938361,
			TransactionCategory: "Electronics",
			PostedTimeStamp:     "2025-06-12T18:05:30+00:00",
			TransactionType:     "POS",
			SendingAccount:      93339,
			ReceivingAccount:    9233021,
			TransactionNote:     "Best Buy Canada",
		},
		{
			ID:                  7,
			Amount:              3200,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-03T10:45:24+00:00",
			TransactionID:       124,
			PAN:                 4807678678904632,
			TransactionCategory: "Cryptocurrency",
			PostedTimeStamp:     "2025-06-03T10:45:24+00:00",
			TransactionType:     "Purchase",
			SendingAccount:      83839,
			ReceivingAccount:    9233020,
			TransactionNote:     "Coinbase Buy ETH",
		},
		{
			ID:                  8,
			Amount:              2899,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-11T20:30:00+00:00",
			TransactionID:       144,
			PAN:                 4673062314928753,
			TransactionCategory: "Food and Beverage",
			PostedTimeStamp:     "2025-06-11T20:30:00+00:00",
			TransactionType:     "POS",
			SendingAccount:      8569,
			ReceivingAccount:    533020,
			TransactionNote:     "The Keg Steakhouse",
		},
		{
			ID:                  9,
			Amount:              1999,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-08T08:00:00+00:00",
			TransactionID:       190,
			PAN:                 5109473381765575,
			TransactionCategory: "Internet Services",
			PostedTimeStamp:     "2025-06-08T08:00:00+00:00",
			TransactionType:     "Subscription",
			SendingAccount:      63639,
			ReceivingAccount:    4233010,
			TransactionNote:     "Netflix Monthly Fee",
		},
		{
			ID:                  10,
			Amount:              6500,
			MessageType:         "Debit",
			CreatedAt:           "2025-06-15T12:10:30+00:00",
			TransactionID:       209,
			PAN:                 5158563621617519,
			TransactionCategory: "Health Services",
			PostedTimeStamp:     "2025-06-15T12:10:30+00:00",
			TransactionType:     "POS",
			SendingAccount:      13839,
			ReceivingAccount:    244020,
			TransactionNote:     "Green Leaf Pharmacy",
		},
	}
}
