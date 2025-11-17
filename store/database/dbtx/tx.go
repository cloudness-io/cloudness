package dbtx

import "database/sql"

// TxDefault represents default transaction options.
var TxDefault = &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false}

// TxDefaultReadOnly represents default transaction options for read-only transactions.
var TxDefaultReadOnly = &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: true}

// TxSerializable represents serializable transaction options.
var TxSerializable = &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false}
