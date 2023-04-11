package sqlstoragetesting

import (
	"testing"

	"github.com/formancehq/ledger/pkg/storage/sqlstorage"
	ledgerstore "github.com/formancehq/ledger/pkg/storage/sqlstorage/ledger"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/schema"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/utils"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"
	"github.com/stretchr/testify/require"
)

func StorageDriver(t pgtesting.TestingT) *sqlstorage.Driver {
	pgServer := pgtesting.NewPostgresDatabase(t)

	db, err := utils.OpenSQLDB(pgServer.ConnString(), testing.Verbose())
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return sqlstorage.NewDriver("postgres", schema.NewPostgresDB(db), ledgerstore.DefaultStoreConfig)
}