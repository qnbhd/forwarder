package main

type Account struct {
	VkLogin string `json:"login"`
	VkPass  string `json:"pass"`
	TgId    int64  `json:"tg_id"`
}

func saveAccounts(accounts []Account) error {
	return nil
}
func loadAccounts() ([]Account, error) {
	return []Account{}, nil
}
