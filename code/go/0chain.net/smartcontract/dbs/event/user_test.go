package event

import (
	"fmt"
	"os"
	"testing"
	"time"

	"0chain.net/core/config"
	"github.com/0chain/common/core/currency"
	"github.com/0chain/common/core/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	logging.Logger = zap.NewNop()
}

var (
	clientID       = "mock client ID"
	txnHash        = "mock txn hash"
	initialBalance = 10
	count          int64
	clientID2      = clientID + " 2"
)

func TestUserEvent(t *testing.T) {
	t.Skip("only for local debugging, requires local postgresql")
	access := config.DbAccess{
		Enabled:         true,
		Name:            "events_db",
		User:            os.Getenv("POSTGRES_USER"),
		Password:        os.Getenv("POSTGRES_PASSWORD"),
		Host:            os.Getenv("POSTGRES_HOST"),
		Port:            os.Getenv("POSTGRES_PORT"),
		MaxIdleConns:    100,
		MaxOpenConns:    200,
		ConnMaxLifetime: 20 * time.Second,
	}

	eventDb, err := NewEventDbWithoutWorker(access, config.DbSettings{})
	require.NoError(t, err)
	defer eventDb.Close()
	err = eventDb.AutoMigrate()
	require.NoError(t, err)

	user1 := User{
		UserID:  clientID,
		TxnHash: txnHash,
		Balance: currency.Coin(initialBalance),
		Round:   3,
		Nonce:   1,
	}

	err = eventDb.addOrUpdateUsers([]User{user1})
	require.NoError(t, err, "Error while inserting User to event Database")

	eventDb.Get().Table("users").Count(&count)
	require.Equal(t, int64(1), count, "User not getting inserted")

	user, err := eventDb.GetUser(clientID)
	require.NoError(t, err, "Error while fetching user by clientID")
	require.Equal(t, clientID, user.UserID, "Fetched invalid User")
	require.Equal(t, txnHash, user.TxnHash, "Fetched invalid User")
	require.Equal(t, initialBalance, user.Balance, "Fetched invalid User")
	require.Equal(t, 1, user.Nonce, "Fetched invalid User")

	user1.Balance = user1.Balance + 1
	user1.Nonce = user1.Nonce + 1
	err = eventDb.addOrUpdateUsers([]User{user1})
	require.NoError(t, err, "Error while inserting User to event Database")

	eventDb.Get().Table("users").Count(&count)
	require.Equal(t, int64(1), count, "User not getting overwritten")

	user, err = eventDb.GetUser(clientID)
	require.NoError(t, err, "Error while fetching user by clientID")
	require.Equal(t, clientID, user.UserID, "Fetched invalid User")
	require.Equal(t, txnHash, user.TxnHash, "Fetched invalid User")
	require.Equal(t, initialBalance+1, user.Balance, "Fetched invalid User")
	require.Equal(t, 2, user.Nonce, "Fetched invalid User")

	//clientID2 := u.UserID + " 2"
	user2 := User{
		UserID:  clientID2,
		TxnHash: txnHash + " 2",
		Balance: currency.Coin(initialBalance) - 1,
		Round:   10,
		Nonce:   1,
	}
	err = eventDb.addOrUpdateUsers([]User{user2})
	require.NoError(t, err, "Error while inserting User to event Database")

	user, err = eventDb.GetUser(clientID2)
	require.NoError(t, err, "Error while fetching user by clientID")
	require.Equal(t, clientID2, user.UserID, "Fetched invalid User")
	require.Equal(t, 1, user.Nonce, "Fetched invalid User")

	eventDb.Get().Table("users").Count(&count)
	require.Equal(t, int64(2), count, "Should have two separate users in store")
	require.Equal(t, int64(3), count, "Just failing for testing purposes")

	err = eventDb.Drop()
	require.NoError(t, err)
}

func prepareEventDB(t *testing.T) (*EventDb, func()) {
	access := config.DbAccess{
		Enabled:         true,
		Name:            "crud",
		User:            os.Getenv("POSTGRES_USER"),
		Password:        os.Getenv("POSTGRES_PASSWORD"),
		Host:            os.Getenv("POSTGRES_HOST"),
		Port:            os.Getenv("POSTGRES_PORT"),
		MaxIdleConns:    100,
		MaxOpenConns:    200,
		ConnMaxLifetime: 20 * time.Second,
	}

	eventDb, err := NewEventDbWithoutWorker(access, config.DbSettings{})
	require.NoError(t, err)
	err = eventDb.AutoMigrate()
	require.NoError(t, err)

	return eventDb, func() {
		eventDb.Close()
	}
}

func TestAddAndUpdateUsersEvent(t *testing.T) {
	t.Skip("only for local debugging, requires local postgresql")
	eventDb, closeDB := prepareEventDB(t)
	defer closeDB()

	// create new users
	users := make([]User, 10)
	for i := 0; i < 10; i++ {
		users[i] = User{
			UserID:  fmt.Sprintf("u_%v", i),
			TxnHash: fmt.Sprintf("hash_%v", i),
			Balance: currency.Coin(i),
			Nonce:   int64(i),
			Round:   int64(i),
		}
	}

	err := eventDb.addOrUpdateUsers(users)
	require.NoError(t, err, "Error while inserting Users to event Database")

	for i := 0; i < 10; i++ {
		u, err := eventDb.GetUser(fmt.Sprintf("u_%v", i))
		require.NoError(t, err)
		require.Equal(t, users[i].Balance, u.Balance)
		require.Equal(t, users[i].Nonce, u.Nonce)
		require.Equal(t, users[i].TxnHash, u.TxnHash)
		require.Equal(t, users[i].Round, u.Round)
	}

	// update users
	for i := 0; i < 10; i++ {
		users[i] = User{
			UserID:  fmt.Sprintf("u_%v", i),
			TxnHash: fmt.Sprintf("hash_%v", i),
			Balance: currency.Coin(i * 100),
			Nonce:   int64(i + 100),
			Round:   int64(i + 100),
		}
	}

	err = eventDb.addOrUpdateUsers(users)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		u, err := eventDb.GetUser(fmt.Sprintf("u_%v", i))
		require.NoError(t, err)
		require.Equal(t, users[i].Balance, u.Balance)
		require.Equal(t, users[i].Nonce, u.Nonce)
		require.Equal(t, users[i].TxnHash, u.TxnHash)
		require.Equal(t, users[i].Round, u.Round)
	}

	users = make([]User, 10)

	// add and update
	for i := 5; i < 15; i++ {
		users[i-5] = User{
			UserID:  fmt.Sprintf("u_%v", i),
			TxnHash: fmt.Sprintf("hash_%v", i),
			Balance: currency.Coin(i * 150),
			Nonce:   int64(i + 150),
			Round:   int64(i + 150),
		}
	}

	err = eventDb.addOrUpdateUsers(users)
	require.NoError(t, err)

	for i := 5; i < 15; i++ {
		u, err := eventDb.GetUser(fmt.Sprintf("u_%v", i))
		require.NoError(t, err)
		require.Equal(t, users[i-5].Balance, u.Balance)
		require.Equal(t, users[i-5].Nonce, u.Nonce)
		require.Equal(t, users[i-5].TxnHash, u.TxnHash)
		require.Equal(t, users[i-5].Round, u.Round)
	}
}

func makeUserTotalStakeEvent(id string, amount int64) Event {
	return Event{
		Type:  TypeStats,
		Tag:   TagLockStakePool,
		Index: id,
		Data: DelegatePoolLock{
			Client: id,
			Amount: amount,
		},
	}
}

func makeUserReadPoolLockEvent(id string, amount int64) Event {
	return Event{
		Type:  TypeStats,
		Tag:   TagLockReadPool,
		Index: id,
		Data: ReadPoolLock{
			Client: id,
			Amount: amount,
		},
	}
}

func makeUserWritePoolLockEvent(id string, amount int64) Event {
	return Event{
		Type:  TypeStats,
		Tag:   TagLockWritePool,
		Index: id,
		Data: WritePoolLock{
			Client: id,
			Amount: amount,
		},
	}
}

func TestMergeUpdateUserTotalStakeEvents(t *testing.T) {
	type expect struct {
		pools  map[string]DelegatePoolLock
		others []Event
	}

	tt := []struct {
		name      string
		events    []Event
		round     int64
		blockHash string
		expect    expect
	}{
		{
			name: "two different clients",
			events: []Event{
				makeUserTotalStakeEvent("c_1", 100),
				makeUserTotalStakeEvent("c_2", 200),
			},
			expect: expect{
				pools: map[string]DelegatePoolLock{
					"c_1": {Client: "c_1", Amount: 100},
					"c_2": {Client: "c_2", Amount: 200},
				},
			},
		},
		{
			name: "two same clients",
			events: []Event{
				makeUserTotalStakeEvent("c_1", 100),
				makeUserTotalStakeEvent("c_1", 200),
			},
			expect: expect{
				pools: map[string]DelegatePoolLock{
					"c_1": {Client: "c_1", Amount: 300},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			em := mergeUserStakeEvents()
			others := make([]Event, 0, len(tc.events))
			for _, e := range tc.events {
				if em.filter(e) {
					continue
				}

				others = append(others, e)
			}

			mergedEvent, err := em.merge(tc.round, tc.blockHash)
			require.NoError(t, err)

			if mergedEvent == nil {
				return
			}

			pools, ok := fromEvent[[]DelegatePoolLock](mergedEvent.Data)
			require.True(t, ok)

			require.Equal(t, len(tc.expect.pools), len(*pools))

			for _, p := range *pools {
				exp, ok := tc.expect.pools[p.Client]
				require.True(t, ok)
				require.EqualValues(t, exp, p)
			}
		})
	}
}

func TestMergeUpdateUserReadPoolLockEvents(t *testing.T) {
	type expect struct {
		pools  map[string]ReadPoolLock
		others []Event
	}

	tt := []struct {
		name      string
		events    []Event
		round     int64
		blockHash string
		expect    expect
	}{
		{
			name: "two different clients",
			events: []Event{
				makeUserReadPoolLockEvent("c_1", 100),
				makeUserReadPoolLockEvent("c_2", 200),
			},
			expect: expect{
				pools: map[string]ReadPoolLock{
					"c_1": {Client: "c_1", Amount: 100},
					"c_2": {Client: "c_2", Amount: 200},
				},
			},
		},
		{
			name: "two same clients",
			events: []Event{
				makeUserReadPoolLockEvent("c_1", 100),
				makeUserReadPoolLockEvent("c_1", 200),
			},
			expect: expect{
				pools: map[string]ReadPoolLock{
					"c_1": {Client: "c_1", Amount: 300},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			em := mergeUserReadPoolLockEvents()
			others := make([]Event, 0, len(tc.events))
			for _, e := range tc.events {
				if em.filter(e) {
					continue
				}

				others = append(others, e)
			}

			mergedEvent, err := em.merge(tc.round, tc.blockHash)
			require.NoError(t, err)

			if mergedEvent == nil {
				return
			}

			pools, ok := fromEvent[[]ReadPoolLock](mergedEvent.Data)
			require.True(t, ok)

			require.Equal(t, len(tc.expect.pools), len(*pools))

			for _, p := range *pools {
				exp, ok := tc.expect.pools[p.Client]
				require.True(t, ok)
				require.EqualValues(t, exp, p)
			}
		})
	}
}

func TestMergeUpdateUserWritePoolLockEvents(t *testing.T) {
	type expect struct {
		pools  map[string]WritePoolLock
		others []Event
	}

	tt := []struct {
		name      string
		events    []Event
		round     int64
		blockHash string
		expect    expect
	}{
		{
			name: "two different clients",
			events: []Event{
				makeUserWritePoolLockEvent("c_1", 100),
				makeUserWritePoolLockEvent("c_2", 200),
			},
			expect: expect{
				pools: map[string]WritePoolLock{
					"c_1": {Client: "c_1", Amount: 100},
					"c_2": {Client: "c_2", Amount: 200},
				},
			},
		},
		{
			name: "two same clients",
			events: []Event{
				makeUserWritePoolLockEvent("c_1", 100),
				makeUserWritePoolLockEvent("c_1", 200),
			},
			expect: expect{
				pools: map[string]WritePoolLock{
					"c_1": {Client: "c_1", Amount: 300},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			em := mergeUserWritePoolLockEvents()
			others := make([]Event, 0, len(tc.events))
			for _, e := range tc.events {
				if em.filter(e) {
					continue
				}

				others = append(others, e)
			}

			mergedEvent, err := em.merge(tc.round, tc.blockHash)
			require.NoError(t, err)

			if mergedEvent == nil {
				return
			}

			pools, ok := fromEvent[[]WritePoolLock](mergedEvent.Data)
			require.True(t, ok)

			require.Equal(t, len(tc.expect.pools), len(*pools))

			for _, p := range *pools {
				exp, ok := tc.expect.pools[p.Client]
				require.True(t, ok)
				require.EqualValues(t, exp, p)
			}
		})
	}
}
