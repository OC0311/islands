package block

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutOuts struct {
	UTXOS []*UTXO
}

func (t *TxOutOuts) Serialize() []byte {

	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(t)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func UnSerializeTxOutOuts(bt []byte) *TxOutOuts {
	var tx TxOutOuts
	decoder := gob.NewDecoder(bytes.NewReader(bt))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}
