package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/joho/godotenv"
	"github.com/l3dlp/logfile"
	"golang.org/x/time/rate"
)

// Define a simple rate limiter
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit
	b   int
}

// Create a new rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

// Get or create a limiter for an IP
func (l *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.ips[ip] = limiter
	}

	return limiter
}

// Create a custom server that uses our rate limiter
func main() {

	// load configuration from .env (non-fatal if the file is absent)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded: %v", err)
	}

	// the recruitment contact address is configuration, not hardcoded
	contactEmail := os.Getenv("HONEYPOT_EMAIL")
	if contactEmail == "" {
		log.Fatal("HONEYPOT_EMAIL is not set (define it in .env)")
	}

	// setting up log file
	logFile := logfile.Use("/var/log/honeypot.log")
	if logFile != nil {
		defer logFile.Close()
	}

	// Create a rate limiter: 3 requests per minute with burst of 5
	limiter := NewIPRateLimiter(rate.Limit(0.05), 5)

	// Handle SSH sessions
	ssh.Handle(func(s ssh.Session) {
		// Get client IP
		remoteAddr := s.RemoteAddr().String()
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			log.Printf("Error getting IP: %v", err)
			return
		}

		// Apply rate limiting
		if !limiter.GetLimiter(ip).Allow() {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			io.WriteString(s, "\n\nToo many connection attempts. Please try again later.\n\n")
			time.Sleep(2 * time.Second) // Delay before closing
			return
		}

		// Log connection
		log.Printf("Connection from: %s", ip)
		
		// Normal response
		io.WriteString(s, fmt.Sprintf("\n\nWant to join a nice remote team?\nSend an e-mail to %s\n\n", contactEmail))
	})

	// Start server on port 22 (requires root privileges)
	log.Println("Starting SSH honeypot server on port 22...")
	log.Fatal(ssh.ListenAndServe(":22", nil))
}

