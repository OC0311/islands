package block

type TXOutput struct {
	Value        int64
	ScriptPubKey string //用户名
}

func (out *TXOutput) UnLockWithAddress(address string) bool {
	return out.ScriptPubKey == address
}
