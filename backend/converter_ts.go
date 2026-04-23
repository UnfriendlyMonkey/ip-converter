package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const humanLayout = "02.01.2006 15:04:05"

// convertToHuman converts a unix timestamp (seconds or ms) to DD.MM.YYYY HH:mm:ss UTC.
// Timestamps > 9_999_999_999 are treated as milliseconds.
func convertToHuman(s string) ConversionResult {
	s = strings.TrimSpace(s)
	res := ConversionResult{Input: s}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		res.Type = "error"
		res.Error = "expected unix timestamp (integer seconds or milliseconds)"
		return res
	}

	var t time.Time
	if n > 9_999_999_999 {
		t = time.Unix(n/1000, (n%1000)*1_000_000).UTC()
	} else {
		t = time.Unix(n, 0).UTC()
	}

	res.Type = "timestamp"
	res.Output = t.Format(humanLayout) + " UTC"
	return res
}

// convertToUnix converts DD.MM.YYYY HH:mm:ss (UTC) to a unix timestamp in seconds.
func convertToUnix(s string) ConversionResult {
	s = strings.TrimSpace(s)
	res := ConversionResult{Input: s}

	clean := strings.TrimSuffix(s, " UTC")
	if clean == s {
		clean = strings.TrimSuffix(s, "UTC")
	}
	clean = strings.TrimSpace(clean)

	t, err := time.ParseInLocation(humanLayout, clean, time.UTC)
	if err != nil {
		res.Type = "error"
		res.Error = fmt.Sprintf("expected DD.MM.YYYY HH:mm:ss, got %q", s)
		return res
	}

	res.Type = "timestamp"
	res.Output = strconv.FormatInt(t.Unix(), 10)
	return res
}

func handleTimestamp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Direction != "to_human" && req.Direction != "to_unix" {
		http.Error(w, `direction must be "to_human" or "to_unix"`, http.StatusBadRequest)
		return
	}

	resp := ConversionResponse{
		Results: make([]ConversionResult, 0, len(req.Values)),
	}

	for _, v := range req.Values {
		var result ConversionResult
		if req.Direction == "to_human" {
			result = convertToHuman(v)
		} else {
			result = convertToUnix(v)
		}
		resp.Results = append(resp.Results, result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
