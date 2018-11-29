package utari

type Transaction struct {
	Txid      string  `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	Output    string  `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
	Input     string  `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	Amount    float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Timestamp string  `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Sign      string  `protobuf:"bytes,6,opt,name=sign,proto3" json:"sign,omitempty"`
	Pubkey    string  `protobuf:"bytes,7,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
}
