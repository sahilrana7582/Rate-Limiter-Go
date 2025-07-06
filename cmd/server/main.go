package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sahilrana7582/rate-limiter-go/internal/limiter"
	"github.com/sahilrana7582/rate-limiter-go/internal/middleware"
)

func main() {
	lim := limiter.NewLimiter(1, 5, 30*time.Second)
	lim.WatchCleaner()
	lim.StartCleaning()

	mw := middleware.NewRateLimiter(lim)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "✅ Request succeeded!")
	})

	http.Handle("/", mw.Limit(testHandler))

	fmt.Println("🚀 Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
