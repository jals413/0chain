package storagesc

import (
	// "context"
	"encoding/json"
	// "net/url"
	"testing"
	"time"

	chainState "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/tokenpool"
	"0chain.net/chaincore/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// test extension
//
func (rp *readPool) total(now int64) state.Balance {
	return rp.Pools.total(now)
}

func (rp *readPool) allocTotal(allocID string,
	now int64) state.Balance {

	return rp.Pools.allocTotal(allocID, now)
}

func (rp *readPool) allocBlobberTotal(allocID, blobberID string,
	now int64) state.Balance {

	return rp.Pools.allocBlobberTotal(allocID, blobberID, now)
}

func mustEncode(t *testing.T, val interface{}) []byte {
	var err error
	b, err := json.Marshal(val)
	require.NoError(t, err)
	return b
}

func requireErrMsg(t *testing.T, err error, msg string) {
	t.Helper()
	require.Error(t, err, "missing error")
	require.Equal(t, msg, err.Error(), "unexpected error")
}

func Test_lockRequest_decode(t *testing.T) {
	var lre, lrd lockRequest
	lre.Duration = time.Second * 60
	lre.AllocationID = "alloc_hex"
	lre.BlobberID = "blobber_hex"
	require.NoError(t, lrd.decode(mustEncode(t, &lre)))
	assert.EqualValues(t, lre, lrd)
}

func Test_unlockRequest_decode(t *testing.T) {
	var ure, urd unlockRequest
	ure.PoolID = "pool_hex"
	require.NoError(t, urd.decode(mustEncode(t, ure)))
	assert.EqualValues(t, ure, urd)
}

func Test_readPool_Encode_Decode(t *testing.T) {
	var rpe, rpd readPool
	rpe.Pools.add(&allocationPool{
		ZcnPool: tokenpool.ZcnPool{
			TokenPool: tokenpool.TokenPool{
				ID: "IDENTIFIER", Balance: 100500,
			},
		},
		AllocationID: "ALLOCATION ID",
		Blobbers: blobberPools{
			&blobberPool{
				BlobberID: "BLOBBER ID",
				Balance:   10300,
			},
		},
		ExpireAt: 90210,
	})
	require.NoError(t, json.Unmarshal(mustEncode(t, rpe), &rpd))
	assert.EqualValues(t, rpe, rpd)
}

func Test_readPoolKey(t *testing.T) {
	assert.NotZero(t, readPoolKey("scKey", "clientID"))
}

func Test_readPools_moveToBlobber(t *testing.T) {

	const (
		sscID   = ADDRESS
		allocID = "alloc_hex"
		blobID  = "blob_hex"
		errMsg  = "not enough tokens in read pool"
	)

	// TODO (sfdx): REMADE THE TEST CASE, MOTHERFUCKER

}

func TestStorageSmartContract_getReadPoolBytes(t *testing.T) {
	const (
		clientID = "client_id"
		errMsg1  = "value not present"
	)

	var (
		ssc      = newTestStorageSC()
		balances = newTestBalances()

		rp *readPool

		b, err = ssc.getReadPoolBytes(clientID, balances)
	)

	requireErrMsg(t, err, errMsg1)
	rp = new(readPool)
	require.NoError(t, rp.save(ssc.ID, clientID, balances))
	b, err = ssc.getReadPoolBytes(clientID, balances)
	require.NoError(t, err)
	assert.EqualValues(t, rp.Encode(), b)
}

func TestStorageSmartContract_getReadPool(t *testing.T) {
	const (
		clientID = "client_id"
		errMsg1  = "value not present"
	)

	var (
		ssc      = newTestStorageSC()
		balances = newTestBalances()
		rps, err = ssc.getReadPool(clientID, balances)
		nrps     = new(readPool)
	)

	requireErrMsg(t, err, errMsg1)
	nrps = new(readPool)
	require.NoError(t, nrps.save(ssc.ID, clientID, balances))
	rps, err = ssc.getReadPool(clientID, balances)
	require.NoError(t, err)
	require.EqualValues(t, nrps, rps)
}

func TestStorageSmartContract_newReadPool(t *testing.T) {
	const (
		clientID, txHash = "client_id", "tx_hash"
		errMsg           = "new_read_pool_failed: already exist"
	)

	var (
		ssc      = newTestStorageSC()
		balances = newTestBalances()
		tx       = transaction.Transaction{
			ClientID:   clientID,
			ToClientID: ssc.ID,
			Value:      0,
		}
		resp string
		err  error
	)

	balances.txn = &tx
	tx.Hash = txHash

	resp, err = ssc.newReadPool(&tx, nil, balances)
	require.NoError(t, err)
	var nrp = new(readPool)
	assert.Equal(t, string(nrp.Encode()), resp)

	_, err = ssc.newReadPool(&tx, nil, balances)
	requireErrMsg(t, err, errMsg)
}

func testSetReadPoolConfig(t *testing.T, rpc *readPoolConfig,
	balances chainState.StateContextI, sscID string) {

	var (
		conf scConfig
		err  error
	)
	conf.ReadPool = rpc
	_, err = balances.InsertTrieNode(scConfigKey(sscID), &conf)
	require.NoError(t, err)
}

func TestStorageSmartContract_readPoolLock(t *testing.T) {
	const (
		allocID, txHash = "alloc_hex", "tx_hash"

		errMsg1 = "read_pool_lock_failed: value not present"
		errMsg2 = "read_pool_lock_failed: " +
			"invalid character '}' looking for beginning of value"
		errMsg3 = "read_pool_lock_failed: no tokens to lock"
		errMsg4 = "read_pool_lock_failed: insufficient amount to lock"
		errMsg5 = "read_pool_lock_failed: " +
			"duration (5s) is shorter than min lock period (10s)"
		errMsg6 = "read_pool_lock_failed: " +
			"duration (2m30s) is longer than max lock period (1m40s)"
		errMsg7 = "read_pool_lock_failed: user already has this read pool"
	)

	var (
		ssc      = newTestStorageSC()
		balances = newTestBalances()
		client   = newClient(0, balances)
		tx       = transaction.Transaction{
			ClientID:   client.id,
			ToClientID: ssc.ID,
			Value:      0,
		}
		lr   lockRequest
		resp string
		err  error
	)

	// setup transaction

	balances.txn = &tx
	tx.Hash = txHash

	// setup config

	testSetReadPoolConfig(t, &readPoolConfig{
		MinLock:       10,
		MinLockPeriod: 10 * time.Second,
		MaxLockPeriod: 100 * time.Second,
	}, balances, ssc.ID)

	// 1. no pool
	_, err = ssc.readPoolLock(&tx, nil, balances)
	requireErrMsg(t, err, errMsg1)

	tx.Hash = "new_read_pool_tx_hash"
	_, err = ssc.newReadPool(&tx, nil, balances)
	require.NoError(t, err)
	tx.Hash = txHash
	// 2. malformed request
	_, err = ssc.readPoolLock(&tx, []byte("} malformed {"), balances)
	requireErrMsg(t, err, errMsg2)
	// 3. min lock
	tx.Value = 5
	lr.Duration = 5 * time.Second
	lr.AllocationID = allocID
	_, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	requireErrMsg(t, err, errMsg4)
	// // 4. min lock
	// balances.balances[clientID] = 5
	// _, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	// requireErrMsg(t, err, errMsg4)
	tx.Value = 15
	balances.balances[client.id] = 15
	// 5. min lock period
	_, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	requireErrMsg(t, err, errMsg5)
	// 6. max lock period
	lr.Duration = 150 * time.Second
	_, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	requireErrMsg(t, err, errMsg6)
	// 7. no such allocation
	lr.Duration = 15 * time.Second
	resp, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	require.Error(t, err)

	balances.balances[client.id] = 200e10
	var aid, _ = addAllocation(t, ssc, client, 10, int64(toSeconds(time.Hour)),
		balances)
	// lock
	lr.AllocationID = aid
	resp, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	require.NoError(t, err)
	assert.NotZero(t, resp)
}

/*
func TestStorageSmartContract_readPoolUnlock(t *testing.T) {
	const (
		allocID                      = "alloc_hex"
		clientID, txHash, readPoolID = "client_id", "tx_hash", "pool_id"

		errMsg1 = "read_pool_unlock_failed: value not present"
		errMsg2 = "read_pool_unlock_failed: " +
			"invalid character '}' looking for beginning of value"
		errMsg3 = "read_pool_unlock_failed: pool not found"
		errMsg4 = "read_pool_unlock_failed: " +
			"emptying pool failed: pool is still locked"
	)

	var (
		ssc      = newTestStorageSC()
		balances = newTestBalances()
		tx       = transaction.Transaction{
			ClientID:   clientID,
			ToClientID: ssc.ID,
			Value:      0,
		}
		lr   lockRequest
		ur   unlockRequest
		resp string
		err  error
	)

	balances.txn = &tx
	tx.Hash = txHash

	// 1. no read pools
	_, err = ssc.readPoolUnlock(&tx, nil, balances)
	requireErrMsg(t, err, errMsg1)

	// create read pool
	tx.Hash = "create_read_pool_tx"
	_, err = ssc.newReadPool(&tx, nil, balances)
	require.NoError(t, err)
	tx.Hash = txHash

	// 2. malformed request
	_, err = ssc.readPoolUnlock(&tx, []byte("} malformed {"), balances)
	requireErrMsg(t, err, errMsg2)

	// 3. no read pool
	_, err = ssc.readPoolUnlock(&tx, mustEncode(t, &ur), balances)
	requireErrMsg(t, err, errMsg3)

	// lock tokens
	testSetReadPoolConfig(t, &readPoolConfig{
		MinLock:       10,
		MinLockPeriod: 10 * time.Second,
		MaxLockPeriod: 100 * time.Second,
	}, balances, ssc.ID)
	tx.Hash = readPoolID
	lr.Duration = 15 * time.Second
	lr.AllocationID = allocID
	balances.balances[clientID] = 150
	tx.Value = 150
	_, err = ssc.readPoolLock(&tx, mustEncode(t, &lr), balances)
	require.NoError(t, err)

	delete(balances.balances, clientID)
	tx.Value = 0

	// 4. not expired
	ur.PoolID = readPoolID
	ur.AllocationID = allocID
	_, err = ssc.readPoolUnlock(&tx, mustEncode(t, &ur), balances)
	requireErrMsg(t, err, errMsg4)

	tx.CreationDate = common.Timestamp(20 * time.Second)

	// 5. unlock (ok)
	resp, err = ssc.readPoolUnlock(&tx, mustEncode(t, &ur), balances)
	require.NoError(t, err)
	assert.NotZero(t, resp)

}

func TestStorageSmartContract_getReadPoolStatsHandler(t *testing.T) {

	const (
		allocID  = "alloc_hex"
		clientID = "client_id"
		errMsg1  = "value not present"
	)

	var (
		ssc      = newTestStorageSC()
		balances = newTestBalances()
		ctx      = context.Background()
		params   = url.Values{
			"client_id": []string{clientID},
		}
		resp, err = ssc.getReadPoolStatsHandler(ctx, params, balances)
	)

	requireErrMsg(t, err, errMsg1)

	var (
		rps = newReadPools()
		rp  = newReadPool()
	)

	rp.ID = "pool_id"
	rp.TokenLockInterface = &tokenLock{
		StartTime: 150,
		Duration:  10 * time.Second,
		Owner:     "owner_id",
	}
	rp.Balance = 150

	require.NoError(t, rps.addPool(allocID, rp))
	require.NoError(t, rps.save(ssc.ID, clientID, balances))

	resp, err = ssc.getReadPoolStatsHandler(ctx, params, balances)
	require.NoError(t, err)
	assert.NotZero(t, resp)
}
*/
