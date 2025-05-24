package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/rates/v2/internal/model"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the Rates API v2")
}

func Rates(w http.ResponseWriter, r *http.Request) {
	rates := []model.Rate{
		{
			InputSymbol: "ETH",
			ProjectName: "Project Alpha",
			PoolName:    "Alpha Pool 1",
			APY:         4.25,
			ProjectLink: "https://project-alpha.example.com",
		},
		{
			InputSymbol: "USDC",
			ProjectName: "Project Beta",
			PoolName:    "Beta Pool 2",
			APY:         2.10,
			ProjectLink: "https://project-beta.example.com",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
