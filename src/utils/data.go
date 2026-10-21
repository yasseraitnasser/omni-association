package utils

import "time"

type Data struct {
	PageTitle string
}

type DashboardData struct {
	PageTitle      string
	AssociationFee int
	Member         struct {
		Name         string
		Role         string
		CreatedAt    time.Time
		TotalPaid    int
		ProjectCount int
	}
	History []struct {
		Year       int
		AmountPaid int
		ReceiptURL string
		PaidAt     time.Time
	}
}
