package storagesc

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// MarshalMsg implements msgp.Marshaler
func (z *storageNodeV1) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 12
	// string "Provider"
	o = append(o, 0x8c, 0xa8, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72)
	o, err = z.Provider.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Provider")
		return
	}
	// string "BaseURL"
	o = append(o, 0xa7, 0x42, 0x61, 0x73, 0x65, 0x55, 0x52, 0x4c)
	o = msgp.AppendString(o, z.BaseURL)
	// string "Terms"
	o = append(o, 0xa5, 0x54, 0x65, 0x72, 0x6d, 0x73)
	o, err = z.Terms.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Terms")
		return
	}
	// string "Capacity"
	o = append(o, 0xa8, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79)
	o = msgp.AppendInt64(o, z.Capacity)
	// string "Allocated"
	o = append(o, 0xa9, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.Allocated)
	// string "PublicKey"
	o = append(o, 0xa9, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79)
	o = msgp.AppendString(o, z.PublicKey)
	// string "SavedData"
	o = append(o, 0xa9, 0x53, 0x61, 0x76, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendInt64(o, z.SavedData)
	// string "DataReadLastRewardRound"
	o = append(o, 0xb7, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x61, 0x64, 0x4c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendFloat64(o, z.DataReadLastRewardRound)
	// string "LastRewardDataReadRound"
	o = append(o, 0xb7, 0x4c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x61, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendInt64(o, z.LastRewardDataReadRound)
	// string "StakePoolSettings"
	o = append(o, 0xb1, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73)
	o, err = z.StakePoolSettings.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "StakePoolSettings")
		return
	}
	// string "RewardRound"
	o = append(o, 0xab, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o, err = z.RewardRound.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "RewardRound")
		return
	}
	// string "NotAvailable"
	o = append(o, 0xac, 0x4e, 0x6f, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65)
	o = msgp.AppendBool(o, z.NotAvailable)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *storageNodeV1) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Provider":
			bts, err = z.Provider.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Provider")
				return
			}
		case "BaseURL":
			z.BaseURL, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BaseURL")
				return
			}
		case "Terms":
			bts, err = z.Terms.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Terms")
				return
			}
		case "Capacity":
			z.Capacity, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Capacity")
				return
			}
		case "Allocated":
			z.Allocated, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Allocated")
				return
			}
		case "PublicKey":
			z.PublicKey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PublicKey")
				return
			}
		case "SavedData":
			z.SavedData, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "SavedData")
				return
			}
		case "DataReadLastRewardRound":
			z.DataReadLastRewardRound, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DataReadLastRewardRound")
				return
			}
		case "LastRewardDataReadRound":
			z.LastRewardDataReadRound, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LastRewardDataReadRound")
				return
			}
		case "StakePoolSettings":
			bts, err = z.StakePoolSettings.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "StakePoolSettings")
				return
			}
		case "RewardRound":
			bts, err = z.RewardRound.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "RewardRound")
				return
			}
		case "NotAvailable":
			z.NotAvailable, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NotAvailable")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *storageNodeV1) Msgsize() (s int) {
	s = 1 + 9 + z.Provider.Msgsize() + 8 + msgp.StringPrefixSize + len(z.BaseURL) + 6 + z.Terms.Msgsize() + 9 + msgp.Int64Size + 10 + msgp.Int64Size + 10 + msgp.StringPrefixSize + len(z.PublicKey) + 10 + msgp.Int64Size + 24 + msgp.Float64Size + 24 + msgp.Int64Size + 18 + z.StakePoolSettings.Msgsize() + 12 + z.RewardRound.Msgsize() + 13 + msgp.BoolSize
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *storageNodeV2) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 14
	// string "Provider"
	o = append(o, 0x8e, 0xa8, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72)
	o, err = z.Provider.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Provider")
		return
	}
	// string "version"
	o = append(o, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Version)
	// string "BaseURL"
	o = append(o, 0xa7, 0x42, 0x61, 0x73, 0x65, 0x55, 0x52, 0x4c)
	o = msgp.AppendString(o, z.BaseURL)
	// string "Terms"
	o = append(o, 0xa5, 0x54, 0x65, 0x72, 0x6d, 0x73)
	o, err = z.Terms.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Terms")
		return
	}
	// string "Capacity"
	o = append(o, 0xa8, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79)
	o = msgp.AppendInt64(o, z.Capacity)
	// string "Allocated"
	o = append(o, 0xa9, 0x41, 0x6c, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendInt64(o, z.Allocated)
	// string "PublicKey"
	o = append(o, 0xa9, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79)
	o = msgp.AppendString(o, z.PublicKey)
	// string "SavedData"
	o = append(o, 0xa9, 0x53, 0x61, 0x76, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61)
	o = msgp.AppendInt64(o, z.SavedData)
	// string "DataReadLastRewardRound"
	o = append(o, 0xb7, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x61, 0x64, 0x4c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendFloat64(o, z.DataReadLastRewardRound)
	// string "LastRewardDataReadRound"
	o = append(o, 0xb7, 0x4c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x61, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o = msgp.AppendInt64(o, z.LastRewardDataReadRound)
	// string "StakePoolSettings"
	o = append(o, 0xb1, 0x53, 0x74, 0x61, 0x6b, 0x65, 0x50, 0x6f, 0x6f, 0x6c, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73)
	o, err = z.StakePoolSettings.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "StakePoolSettings")
		return
	}
	// string "RewardRound"
	o = append(o, 0xab, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64)
	o, err = z.RewardRound.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "RewardRound")
		return
	}
	// string "NotAvailable"
	o = append(o, 0xac, 0x4e, 0x6f, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65)
	o = msgp.AppendBool(o, z.NotAvailable)
	// string "IsRestricted"
	o = append(o, 0xac, 0x49, 0x73, 0x52, 0x65, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x65, 0x64)
	if z.IsRestricted == nil {
		o = msgp.AppendNil(o)
	} else {
		o = msgp.AppendBool(o, *z.IsRestricted)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *storageNodeV2) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Provider":
			bts, err = z.Provider.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Provider")
				return
			}
		case "version":
			z.Version, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		case "BaseURL":
			z.BaseURL, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "BaseURL")
				return
			}
		case "Terms":
			bts, err = z.Terms.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Terms")
				return
			}
		case "Capacity":
			z.Capacity, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Capacity")
				return
			}
		case "Allocated":
			z.Allocated, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Allocated")
				return
			}
		case "PublicKey":
			z.PublicKey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "PublicKey")
				return
			}
		case "SavedData":
			z.SavedData, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "SavedData")
				return
			}
		case "DataReadLastRewardRound":
			z.DataReadLastRewardRound, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "DataReadLastRewardRound")
				return
			}
		case "LastRewardDataReadRound":
			z.LastRewardDataReadRound, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "LastRewardDataReadRound")
				return
			}
		case "StakePoolSettings":
			bts, err = z.StakePoolSettings.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "StakePoolSettings")
				return
			}
		case "RewardRound":
			bts, err = z.RewardRound.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "RewardRound")
				return
			}
		case "NotAvailable":
			z.NotAvailable, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "NotAvailable")
				return
			}
		case "IsRestricted":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.IsRestricted = nil
			} else {
				if z.IsRestricted == nil {
					z.IsRestricted = new(bool)
				}
				*z.IsRestricted, bts, err = msgp.ReadBoolBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "IsRestricted")
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *storageNodeV2) Msgsize() (s int) {
	s = 1 + 9 + z.Provider.Msgsize() + 8 + msgp.StringPrefixSize + len(z.Version) + 8 + msgp.StringPrefixSize + len(z.BaseURL) + 6 + z.Terms.Msgsize() + 9 + msgp.Int64Size + 10 + msgp.Int64Size + 10 + msgp.StringPrefixSize + len(z.PublicKey) + 10 + msgp.Int64Size + 24 + msgp.Float64Size + 24 + msgp.Int64Size + 18 + z.StakePoolSettings.Msgsize() + 12 + z.RewardRound.Msgsize() + 13 + msgp.BoolSize + 13
	if z.IsRestricted == nil {
		s += msgp.NilSize
	} else {
		s += msgp.BoolSize
	}
	return
}
