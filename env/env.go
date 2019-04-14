package env

var (
PGHost     string
PGPort     int
PGDB       string
PGUser     string
PGPassword string

DateTo   string
DateFrom string

DBConnMaxLifetime int
DBMaxIdleConns    int
DBMaxOpenConns    int
)
