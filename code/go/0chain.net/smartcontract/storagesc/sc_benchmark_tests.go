package storagesc

import (
	"encoding/json"
	"strconv"
	"time"

	sc "0chain.net/smartcontract/benchmark"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/core/encryption"

	"github.com/spf13/viper"

	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/datastore"
)

type BenchTest struct {
	name     string
	endpoint func(
		*transaction.Transaction,
		[]byte,
		cstate.StateContextI,
	) (string, error)
	txn   transaction.Transaction
	input []byte
}

func (bt BenchTest) Name() string {
	return bt.name
}

func (bt BenchTest) Transaction() transaction.Transaction {
	return bt.txn
}

func (bt BenchTest) Run(balances cstate.StateContextI) {
	bt.endpoint(&bt.txn, bt.input, balances)
}

func BenchmarkTests(
	vi *viper.Viper,
	clients []string,
	keys []string,
	blobbers []string,
	allocations []string,
) []BenchTest {
	var now = common.Timestamp(vi.GetInt64(sc.Now))
	var ssc = StorageSmartContract{
		SmartContract: sci.NewSC(ADDRESS),
	}
	return []BenchTest{
		/*
			{
				name:     "storage_read_redeem",
				endpoint: ssc.commitBlobberRead,
				txn:      transaction.Transaction{},
				input: func() []byte {
					bytes := (&ReadConnection{
						ReadMarker: &ReadMarker{
							ClientID:        clients[0],
							ClientPublicKey: keys[0],
							BlobberID:       blobbers[0],
							AllocationID:    allocations[0],
							OwnerID:         clients[0],
							Timestamp:       now,
							ReadCounter:     1,
							Signature: "", // todo work out how to sign
							PayerID:         clients[0],
						},
					}).Encode()
					return bytes
				}(),
			},
		*/
		// allocations
		{
			name:     "storage_new_allocation_request",
			endpoint: ssc.newAllocationRequest,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				ClientID:     clients[0],
				CreationDate: now,
				Value:        100 * vi.GetInt64(sc.StorageMinAllocSize),
			},
			input: func() []byte {
				bytes, _ := (&newAllocationRequest{
					DataShards:                 4,
					ParityShards:               4,
					Size:                       100 * vi.GetInt64(sc.StorageMinAllocSize),
					Expiration:                 common.Timestamp(vi.GetDuration(sc.StorageMinAllocDuration).Seconds()) + now,
					Owner:                      clients[0],
					OwnerPublicKey:             keys[0],
					PreferredBlobbers:          []string{},
					ReadPriceRange:             PriceRange{0, state.Balance(vi.GetInt64(sc.StorageMaxReadPrice) * 1e10)},
					WritePriceRange:            PriceRange{0, state.Balance(vi.GetInt64(sc.StorageMaxWritePrice) * 1e10)},
					MaxChallengeCompletionTime: vi.GetDuration(sc.StorageMaxChallengeCompletionTime),
					DiversifyBlobbers:          false,
				}).encode()
				return bytes
			}(),
		},
		{
			name:     "storage_update_allocation_request",
			endpoint: ssc.updateAllocationRequest,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				ClientID: clients[0],
				Value:    100 * vi.GetInt64(sc.StorageMinAllocSize),
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&updateAllocationRequest{
					ID:           allocations[0],
					OwnerID:      clients[0],
					Size:         100 * vi.GetInt64(sc.StorageMinAllocSize),
					Expiration:   common.Timestamp(vi.GetDuration(sc.StorageMinAllocDuration).Seconds()),
					SetImmutable: true,
				})
				return bytes
			}(),
		},
		{
			name:     "storage_finalize_allocation",
			endpoint: ssc.finalizeAllocation,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				CreationDate: common.Timestamp((time.Hour * 1000).Seconds()) + now,
				ClientID:     clients[0],
				ToClientID:   ADDRESS,
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&lockRequest{
					AllocationID: allocations[0],
				})
				return bytes
			}(),
		},
		{
			name:     "storage_cancel_allocation",
			endpoint: ssc.cancelAllocationRequest,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				CreationDate: now + 1,
				ClientID:     clients[0],
				ToClientID:   ADDRESS,
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&lockRequest{
					AllocationID: allocations[0],
				})
				return bytes
			}(),
		},
		// blobbers
		{
			name:     "storage_add_blobber",
			endpoint: ssc.addBlobber,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				CreationDate: now + 1,
				ClientID:     clients[0],
				ToClientID:   ADDRESS,
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&StorageNode{
					ID:                encryption.Hash("my_new_blobber"),
					BaseURL:           "my_new_blobber.com",
					Terms:             getMockBlobberTerms(vi),
					Capacity:          vi.GetInt64(sc.StorageMinBlobberCapacity) * 1000,
					StakePoolSettings: getStakePoolSettings(vi, encryption.Hash("my_new_blobber")),
				})
				return bytes
			}(),
		},
		{
			name:     "storage_add_validator",
			endpoint: ssc.addValidator,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				CreationDate: now + 1,
				ClientID:     clients[0],
				ToClientID:   ADDRESS,
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&ValidationNode{
					ID:                encryption.Hash("my_new_validator"),
					BaseURL:           "my_new_validator.com",
					StakePoolSettings: getStakePoolSettings(vi, encryption.Hash("my_new_validator")),
				})
				return bytes
			}(),
		},
		{
			name:     "storage_blobber_health_check",
			endpoint: ssc.blobberHealthCheck,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				CreationDate: now + 1,
				ClientID:     blobbers[0],
				ToClientID:   ADDRESS,
			},
			input: []byte{},
		},
		{
			name:     "update_blobber_settings",
			endpoint: ssc.updateBlobberSettings,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				CreationDate: now + 1,
				ClientID:     blobbers[0],
				ToClientID:   ADDRESS,
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&StorageNode{
					ID:                blobbers[0],
					Terms:             getMockBlobberTerms(vi),
					Capacity:          vi.GetInt64(sc.StorageMinBlobberCapacity) * 1000,
					StakePoolSettings: getStakePoolSettings(vi, blobbers[0]),
				})
				return bytes
			}(),
		},
		// add_curator
		{
			name:     "storage_add_curator",
			endpoint: ssc.addCurator,
			txn: transaction.Transaction{
				ClientID: clients[0],
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&curatorInput{
					CuratorId:    clients[vi.GetInt(sc.NumCurators)],
					AllocationId: allocations[0],
				})
				return bytes
			}(),
		},
		{
			name:     "storage_remove_curator",
			endpoint: ssc.removeCurator,
			txn: transaction.Transaction{
				ClientID: clients[0],
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&curatorInput{
					CuratorId:    clients[0],
					AllocationId: allocations[0],
				})
				return bytes
			}(),
		},
		// read_pool
		{
			name:     "storage_new_read_pool",
			endpoint: ssc.newReadPool,
			txn:      transaction.Transaction{},
			input:    []byte{},
		}, /*
			{
				name:     "storage_read_pool_unlock",
				endpoint: ssc.readPoolUnlock,
				txn: transaction.Transaction{
					HashIDField: datastore.HashIDField{
						Hash: encryption.Hash("mock transaction hash"),
					},
					Value:      vi.GetInt64(sc.StorageReadPoolMinLock),
					ClientID:   clients[0],
					ToClientID: ADDRESS,
				},
				input: func() []byte {
					bytes, _ := json.Marshal(&unlockRequest{
						PoolID: allocations[0],
					})
					return bytes
				}(),
			},*/
		{
			name:     "storage_read_pool_lock",
			endpoint: ssc.readPoolLock,
			txn: transaction.Transaction{
				HashIDField: datastore.HashIDField{
					Hash: encryption.Hash("mock transaction hash"),
				},
				Value:      vi.GetInt64(sc.StorageReadPoolMinLock),
				ClientID:   clients[0],
				ToClientID: ADDRESS,
			},
			input: func() []byte {
				bytes, _ := json.Marshal(&lockRequest{
					AllocationID: allocations[0],
					TargetId:     getMockReadPoolId(0, 0, 0),
					Duration:     vi.GetDuration(sc.StorageReadPoolMinLockPeriod),
				})
				return bytes
			}(),
		},
		// write pool

		// stake pool
		{
			name:     "storage_stake_pool_pay_interests",
			endpoint: ssc.stakePoolPayInterests,
			txn:      transaction.Transaction{},
			input: func() []byte {
				bytes, _ := json.Marshal(&stakePoolRequest{
					BlobberID: blobbers[0],
					PoolID:    blobbers[0] + "Pool" + strconv.Itoa(0),
				})
				return bytes
			}(),
		},
	}
}
