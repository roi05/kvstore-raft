package main

import (
	"fmt"
	"net/http"

	"github.com/roi-05/kvstore-raft/internal/apiserver"
	"github.com/roi-05/kvstore-raft/internal/store"
)

func main() {
	kv := store.NewStore()
	handler := apiserver.NewHandler(kv)

	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", handler)
}
