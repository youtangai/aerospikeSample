package utari

type Block struct {
	Id         string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Version    int32    `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Prehash    string   `protobuf:"bytes,3,opt,name=prehash,proto3" json:"prehash,omitempty"`
	Merkleroot string   `protobuf:"bytes,4,opt,name=merkleroot,proto3" json:"merkleroot,omitempty"`
	Timestamp  string   `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Level      string   `protobuf:"bytes,6,opt,name=level,proto3" json:"level,omitempty"`
	Nonce      uint32   `protobuf:"varint,7,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Size       int64    `protobuf:"varint,8,opt,name=size,proto3" json:"size,omitempty"`
	Txcount    int64    `protobuf:"varint,9,opt,name=txcount,proto3" json:"txcount,omitempty"`
	TxidList   []string `protobuf:"bytes,10,rep,name=txid_list,json=txidList,proto3" json:"txid_list,omitempty"`
}
