package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Data struct {
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
}

type Member struct {
	ID               int
	Name             string
	Email            string
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
	PageTitle           string
	IsAuthenticated     bool
	AuthenticatedUserID int
	AssociationFee      int
	ProjectCount        int
	Member              Member
	History             []History
}

type Claims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type Project struct {
	ID                        int
	Name                      string
	Description               string
	Status                    string
	Budget                    int
	TotalIncome               int
	TotalExpenses             int
	RemainingBalance          int
	FundingProgressPercentage float64
	CreatedAt                 time.Time
}
