package event

import (
	"testing"

	"0chain.net/core/config"
	"0chain.net/core/viper"
	"github.com/0chain/common/core/logging"
	"github.com/stretchr/testify/require"
)

func init() {
	viper.Set("logging.console", true)
	viper.Set("logging.level", "debug")
}

const req = `SELECT
    child.relname       AS child
FROM pg_inherits
    JOIN pg_class parent            ON pg_inherits.inhparent = parent.oid
    JOIN pg_class child             ON pg_inherits.inhrelid   = child.oid
    JOIN pg_namespace nmsp_parent   ON nmsp_parent.oid  = parent.relnamespace
    JOIN pg_namespace nmsp_child    ON nmsp_child.oid   = child.relnamespace
WHERE parent.relname='events'`

func TestPartitionCreate(t *testing.T) {
	logging.InitLogging("development", "")

	config.Configuration().ChainConfig = &TestConfig{conf: &TestConfigData{DbsSettings: config.DbSettings{
		AggregatePeriod:                10,
		PartitionKeepCount:             10,
		PartitionChangePeriod:          100,
		PermanentPartitionChangePeriod: 10,
		PermanentPartitionKeepCount:    1,
	}}}

	db, f := GetTestEventDB(t)
	defer f()
	err := db.addPartition(11, "events")
	require.NoError(t, err)
	err = db.addPartition(101, "events")
	require.NoError(t, err)
	err = db.addPartition(201, "events")
	require.NoError(t, err)

	var partitions []string
	db.Store.Get().Raw(req).Scan(&partitions)
	require.Equal(t, 13, len(partitions))

	err = db.dropPartition(201, "events")
	require.NoError(t, err)

	db.Store.Get().Raw(req).Scan(&partitions)
	require.Equal(t, 13, len(partitions))
}
