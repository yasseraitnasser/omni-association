package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Data struct {
	PageTitle string
}

type Member struct {
	Name             string
	Role             string
	CreatedAt        time.Time
	MembershipStatus string
	TotalPaid        int
}

type History struct {
	Year       int
	AmountPaid int
	ReceiptURL string
	PaidAt     time.Time
}

type DashboardData struct {
	PageTitle      string
	AssociationFee int
	ProjectCount   int
	Member         Member
	History        []History
}

type Claims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
