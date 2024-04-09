package api

import "google.golang.org/protobuf/proto"

func (x *ApplyStandaloneConfigCommand) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *ApplyStandaloneConfigCommand) Unmarshal(cmdMarshalled []byte) error {
	return proto.Unmarshal(cmdMarshalled, x)
}

func (x *ApplyConfigGroupCommand) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *ApplyConfigGroupCommand) Unmarshal(cmdMarshalled []byte) error {
	return proto.Unmarshal(cmdMarshalled, x)
}
