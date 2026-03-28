package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ConversionRequest struct {
	Direction string   `json:"direction"` // "to_numeric" | "to_string"
	Values    []string `json:"values"`
}

type ConversionResult struct {
	Input  string `json:"input"`
	Type   string `json:"type"` // "ipv4" | "ipv4_net" | "ipv6" | "ipv6_net" | "error"
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type ConversionResponse struct {
	Results []ConversionResult `json:"results"`
}

// convertToNumeric converts a string IP/CIDR to its numeric representation.
//
// IPv4 host:    "192.168.1.1"       → output "3232235777"
// IPv4 network: "192.168.1.0/24"    → output "3232235776/24"
// IPv6 host:    "2001:db8::1"       → output "[32,1,13,184,...,1]"
// IPv6 network: "2001:db8::/32"     → output "[32,1,13,184,...,0]/32"
func convertToNumeric(s string) ConversionResult {
	s = strings.TrimSpace(s)
	res := ConversionResult{Input: s}

	if strings.Contains(s, "/") {
		// CIDR notation
		ip, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			res.Type = "error"
			res.Error = fmt.Sprintf("invalid CIDR: %v", err)
			return res
		}

		ones, _ := ipNet.Mask.Size()

		if ip4 := ip.To4(); ip4 != nil {
			// Use the host IP (ip), not the masked network IP (ipNet.IP)
			addr := binary.BigEndian.Uint32(ip4)
			res.Type = "ipv4_net"
			res.Output = fmt.Sprintf("%d/%d", addr, ones)
		} else {
			// IPv6: use host IP bytes
			bytes := ip.To16()
			jsonBytes, _ := json.Marshal([]int(toIntSlice(bytes)))
			res.Type = "ipv6_net"
			res.Output = fmt.Sprintf("%s/%d", string(jsonBytes), ones)
		}
	} else {
		// Host address
		ip := net.ParseIP(s)
		if ip == nil {
			res.Type = "error"
			res.Error = "invalid IP address"
			return res
		}

		if ip4 := ip.To4(); ip4 != nil {
			addr := binary.BigEndian.Uint32(ip4)
			res.Type = "ipv4"
			res.Output = strconv.FormatUint(uint64(addr), 10)
		} else {
			bytes := ip.To16()
			jsonBytes, _ := json.Marshal(toIntSlice(bytes))
			res.Type = "ipv6"
			res.Output = string(jsonBytes)
		}
	}

	return res
}

// convertToString converts a numeric representation back to an IP string.
//
// "3232235777"               → "192.168.1.1"
// "3232235776/24"            → "192.168.1.0/24"
// "[32,1,13,184,...,1]"      → "2001:db8::1"
// "[32,1,13,184,...,0]/32"   → "2001:db8::/32"
func convertToString(s string) ConversionResult {
	s = strings.TrimSpace(s)
	res := ConversionResult{Input: s}

	if strings.HasPrefix(s, "[") {
		// IPv6 byte array
		mask := -1
		byteStr := s

		if idx := strings.LastIndex(s, "]/"); idx != -1 {
			maskStr := s[idx+2:]
			m, err := strconv.Atoi(strings.TrimSpace(maskStr))
			if err != nil {
				res.Type = "error"
				res.Error = "invalid mask: " + maskStr
				return res
			}
			mask = m
			byteStr = s[:idx+1]
		}

		var nums []int
		if err := json.Unmarshal([]byte(byteStr), &nums); err != nil {
			res.Type = "error"
			res.Error = fmt.Sprintf("invalid byte array: %v", err)
			return res
		}

		if len(nums) != 16 {
			res.Type = "error"
			res.Error = fmt.Sprintf("IPv6 requires exactly 16 bytes, got %d", len(nums))
			return res
		}

		bytes := make([]byte, 16)
		for i, n := range nums {
			if n < 0 || n > 255 {
				res.Type = "error"
				res.Error = fmt.Sprintf("byte value out of range at index %d: %d", i, n)
				return res
			}
			bytes[i] = byte(n)
		}

		ip := net.IP(bytes)
		if mask >= 0 {
			// Preserve host IP — just append the mask, don't zero out host bits
			res.Type = "ipv6_net"
			res.Output = fmt.Sprintf("%s/%d", ip.String(), mask)
		} else {
			res.Type = "ipv6"
			res.Output = ip.String()
		}
		return res
	}

	// IPv4 uint32 (with optional /mask)
	mask := -1
	addrStr := s

	if idx := strings.LastIndex(s, "/"); idx != -1 {
		maskStr := s[idx+1:]
		m, err := strconv.Atoi(strings.TrimSpace(maskStr))
		if err == nil {
			mask = m
			addrStr = s[:idx]
		}
	}

	addr, err := strconv.ParseUint(strings.TrimSpace(addrStr), 10, 32)
	if err != nil {
		res.Type = "error"
		res.Error = fmt.Sprintf("invalid value %q: expected uint32 or byte array", addrStr)
		return res
	}

	ip := make(net.IP, 4)
	ip[0] = byte(addr >> 24)
	ip[1] = byte(addr >> 16)
	ip[2] = byte(addr >> 8)
	ip[3] = byte(addr)

	if mask >= 0 {
		// Preserve host IP — just append the mask, don't zero out host bits
		res.Type = "ipv4_net"
		res.Output = fmt.Sprintf("%s/%d", ip.String(), mask)
	} else {
		res.Type = "ipv4"
		res.Output = ip.String()
	}

	return res
}

func toIntSlice(b []byte) []int {
	ints := make([]int, len(b))
	for i, v := range b {
		ints[i] = int(v)
	}
	return ints
}
