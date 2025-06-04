package dpsk

// Model 模型结构体
type Model struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

//  ModelList 模型列表结构体
type ModelList struct {
	Data   []Model `json:"data,omitempty"`
	Object string  `json:"object,omitempty"`
}

// Balance 余额结构体
type Balance struct {
	IsAvailable  bool          `json:"is_available,omitempty"`
	BalanceInfos []BalanceInfo `json:"balance_infos,omitempty"`
}

//  BalanceInfo 余额信息结构体
type BalanceInfo struct {
	Currency        string `json:"currency,omitempty"`          // 货币,人民币或美元
	TotalBalance    string `json:"total_balance,omitempty"`     // 总的可用余额,包括赠金和充值余额
	GrantedBalance  string `json:"granted_balance,omitempty"`   // 未过期的赠金余额
	ToppedUpBalance string `json:"topped_up_balance,omitempty"` // 充值余额
}
