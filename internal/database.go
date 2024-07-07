package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "expense_tracker.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the connection is available
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Enable foreign key support (SQLite specific)
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Create tables
	createTables()
}

func createTables() {
	createUserTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        payment_methods TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        info TEXT,
        password TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	createBillMateTable := `
	CREATE TABLE IF NOT EXISTS bill_mates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	);`

	createBillGroupTable := `
    CREATE TABLE IF NOT EXISTS bill_groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT
    );`

	createBillMateToGroupTable := `
    CREATE TABLE IF NOT EXISTS bill_mate_to_group (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        bill_mate_id INTEGER NOT NULL,
        bill_group_id INTEGER NOT NULL,
        FOREIGN KEY (bill_mate_id) REFERENCES bill_mates(id) ON DELETE CASCADE,
        FOREIGN KEY (bill_group_id) REFERENCES bill_groups(id) ON DELETE CASCADE
    );`

	createCategoryTable := `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT
    );`

	createBudgetTable := `
    CREATE TABLE IF NOT EXISTS budgets (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER NOT NULL,
        time_period INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        amount INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        UNIQUE (user_id, category_id)
    );`

	createExpenseTable := `
    CREATE TABLE IF NOT EXISTS expenses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        amount INTEGER NOT NULL,
        description TEXT,
        title TEXT NOT NULL,
        date DATETIME DEFAULT CURRENT_TIMESTAMP,
        payment_method TEXT,
        vendor TEXT,
        user_id INTEGER NOT NULL,
        bill_group_id INTEGER NOT NULL,
        category_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (bill_group_id) REFERENCES bill_groups(id) ON DELETE SET NULL,
        FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
    );`

	createSplitTable := `
    CREATE TABLE IF NOT EXISTS splits (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        paid REAL NOT NULL,
        responsible REAL NOT NULL,
        bill_mate_id INTEGER NOT NULL,
        expense_id INTEGER NOT NULL,
        FOREIGN KEY (bill_mate_id) REFERENCES bill_mates(id) ON DELETE CASCADE,
        FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE CASCADE
    );`

	createBlacklistTable := `
    CREATE TABLE IF NOT EXISTS blacklist (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        token TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX IF NOT EXISTS idx_token_created_at ON blacklist(token, created_at);`

	// Execute table creation queries
	execQuery(createUserTable)
	execQuery(createBillMateTable)
	execQuery(createBillGroupTable)
	execQuery(createBillMateToGroupTable)
	execQuery(createCategoryTable)
	execQuery(createBudgetTable)
	execQuery(createExpenseTable)
	execQuery(createSplitTable)
	execQuery(createBlacklistTable)
}

func execQuery(query string) {
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
}

func CloseDatabase() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Failed to close database: %v", err)
	}
}
