package utils

func IsBoardMember(claims Claims) bool {
	isBoardMember := claims.Role == "president" ||
		claims.Role == "vice-president" ||
		claims.Role == "treasurer" ||
		claims.Role == "assistant-treasurer" ||
		claims.Role == "general-secretary" ||
		claims.Role == "assistant-general-secretary" ||
		claims.Role == "advisor"
	if !isBoardMember {
		return false
	}
	return true
}
