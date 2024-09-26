package msmysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"github.com/jakobsym/BudgetFi/api/pkg/model"
	"github.com/joho/godotenv"
)

// Use ExecContext() for inserts/deletes/updates
// Use QueryContext() for select statements

type Repository struct {
	sync.RWMutex
}

func New() *Repository {
	return &Repository{}
}

// Checks DB via `google_id` to determine if a new user or not
// returns users 'uuid' as a string if they are
// otherwise returns an empty string indicating a new user
func (r *Repository) PrevUserCheck(ctx context.Context, user *model.User) (string, error) {
	var err error
	db, err := MsSqlConnection()
	if err != nil {
		return "", fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("login: db is null")
		return "", err
	}
	tsql := `SELECT id FROM [dbo].[Users] where google_id = (@Google_ID)`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var userUUID string
	err = stmt.QueryRowContext(ctx, sql.Named("google_id", user.Google_Id)).Scan(&userUUID)

	// TOOD: iffy on conditional logic here
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return userUUID, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	//r.Lock()
	//defer r.Unlock()
	var err error
	db, err := MsSqlConnection()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}

	if db == nil {
		err = errors.New("CreateUser: db is null")
		return err
	}
	tsql := `INSERT INTO [dbo].[Users] (name, email, id, google_id) VALUES (@Name,@Email,@Id,@Google_ID);`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	uuid := user.UUID[:] // convert to []byte instead of [16]byte
	_, err = stmt.ExecContext(ctx,
		sql.Named("name", user.Name),
		sql.Named("email", user.Email),
		sql.Named("id", uuid),
		sql.Named("google_id", user.Google_Id))

	if err != nil {
		return err
	}
	return nil
}

// TODO: Test me
func (r *Repository) CreateCategory(ctx context.Context, category *model.Catergory, userId string) error {
	db, err := MsSqlConnection()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("CreateUser: db is null")
		return err
	}
	category.Category_Id = int(uuid.New().ID())
	tsql := `INSERT INTO [dbo].[Category] (category_id, id, name) VALUES (@CategoryId, @Id, @Name);`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sql.Named("name", category.Name),
		sql.Named("id", userId),
		sql.Named("category_id", category.Category_Id),
	)
	return nil
}

func (r *Repository) CreateCategoryByName(ctx context.Context, categoryName, userId string) error {
	db, err := MsSqlConnection()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("CreateUser: db is null")
		return err
	}
	categoryId := int(uuid.New().ID())
	tsql := `INSERT INTO [dbo].[Category] (category_id, id, name) VALUES (@CategoryId, @Id, @Name);`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, sql.Named("name", categoryName), sql.Named("userId", userId), sql.Named("category_id", categoryId))
	return nil
}

// Returns the ID of found category
// if no category of `categoryName` in DB, a category gets created
// Refactor in future: Have it so that select statement returns the Id?
func (r *Repository) SearchCategoryByName(ctx context.Context, categoryName, userId string) (string, error) {
	var err error
	var categoryId string

	db, err := MsSqlConnection()
	if err != nil {
		return "", fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("SearchCategoryByName: db is null")
		return "", err
	}

	tsql := `SELECT [category_id] FROM [dbo].[Category] WHERE categoryName = @name AND userId = @id;`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()
	err = stmt.QueryRowContext(ctx, sql.Named("name", categoryName), sql.Named("id", userId)).Scan(&categoryId)

	// category_id == null
	if err == sql.ErrNoRows {
		err = r.CreateCategoryByName(ctx, categoryName, userId)
		if err != nil {
			return "", fmt.Errorf("error fetching category_id by name: %v", err)
		}
		// rerun select stmt
		err = stmt.QueryRowContext(ctx, sql.Named("name", categoryName), sql.Named("id", userId)).Scan(&categoryId)
		if err != nil {
			return "", fmt.Errorf("error fetching new category_id by name: %v", err)
		}
	} else if err != nil {
		return "", err
	}
	return categoryId, nil
}

// TODO: Test me
func (r *Repository) DeleteCategoryByName(ctx context.Context, categoryName, userId string) error {
	db, err := MsSqlConnection()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("CreateUser: db is null")
		return err
	}
	tsql := `DELETE FROM [dbo].[Category] WHERE [name] = @Name AND userId = @id;`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, sql.Named("name", categoryName), sql.Named("id", userId))
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Deleted %d row\n", row)
	return nil
}

func (r *Repository) UpdateCategoryName(ctx context.Context, category *model.Catergory, newCategoryName string) error {
	return nil
}

// `Expense` belongs to a category which
// to create an `Expense` it must belong to some category
func (r *Repository) CreateExpense(ctx context.Context, expense *model.Expense, categoryName, userId string) error {
	db, err := MsSqlConnection()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("CreateExpense: db is null")
		return err
	}
	categoryId, err := r.SearchCategoryByName(ctx, categoryName, userId)
	if err != nil {
		return fmt.Errorf("error searching category by name: %v", err)
	}
	expenseId := int(uuid.New().ID())
	tsql := `INSERT INTO [dbo].[Expense] (expense_id, id, category_id, amount, expense_name) VALUES (@ExpenseId, @Id, @Category_id, @Amount, @Expense_name);`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx,
		sql.Named("expense_id", expenseId),
		sql.Named("id", userId),
		sql.Named("category_id", categoryId),
		sql.Named("amount", expense.Amount),
		sql.Named("expense_name", expense.ExpenseName),
	)
	if err != nil {
		return fmt.Errorf("error inserting into Expense Table: %v", err)
	}
	return nil
}

// TODO: Consider the fact that Expenses only exist within a given Category
// This does not work as intended
func (r *Repository) DeleteExpense(ctx context.Context, expenseName, userId string) error {
	db, err := MsSqlConnection()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("DeleteExpense: db is null")
		return err
	}
	tsql := `DELETE FROM [dbo].[Expense] WHERE expense_name = @expenseName AND id = @userId;`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, sql.Named("expense_name", expenseName), sql.Named("id", userId))
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Deleted %d row\n", row)
	return nil
}

// TODO:
// Add Budget to a User will need to pass a category
// Update Budget
// Delete Budget

// Add Expense to User will need to pass a category
// Update Expense
// Delete Expense

// Util
func MsSqlConnection() (*sql.DB, error) {
	var err error
	var port int
	var db *sql.DB

	env := loadDbEnv()
	port, err = strconv.Atoi(env["DB_PORT"])
	if err != nil {
		return nil, fmt.Errorf("error converting string to int: %v", err)
	}
	// connection string
	cString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		env["DB_SERVER"], env["DB_USER"], env["DB_PW"], port, env["DB"])

	db, err = sql.Open("sqlserver", cString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}
	//fmt.Printf("Connected to ms-sqlDB")
	return db, nil
}

// Util
func loadDbEnv() map[string]string {
	err := godotenv.Load("backend.env")
	if err != nil {
		log.Fatal("error loading .env")
	}

	envMap := map[string]string{
		"DB_SERVER": os.Getenv("DB_SERVER"),
		"DB_PORT":   os.Getenv("DB_PORT"),
		"DB_USER":   os.Getenv("DB_USER"),
		"DB_PW":     os.Getenv("DB_PW"),
		"DB":        os.Getenv("DB"),
	}

	return envMap
}

/*
// Returns the name of the category if it exists, otherwise create new category returning the name
func (r *Repository) DeprecatedSearchCategoryByName(ctx context.Context, categoryName, userId string) (string, error) {
	var err error
	db, err := MsSqlConnection()
	if err != nil {
		return "", fmt.Errorf("error establishing DB connection: %v", err)
	}
	if db == nil {
		err = errors.New("SearchCategoryByName: db is null")
		return "", err
	}
	tsql := `SELECT [name] FROM [dbo].[Category] WHERE categoryName = @name AND userId = @id;`
	stmt, err := db.PrepareContext(ctx, tsql)
	if err != nil {
		return "", nil
	}
	defer stmt.Close()
	var cName string
	err = stmt.QueryRowContext(ctx, sql.Named("name", categoryName), sql.Named("id", userId)).Scan(&cName)

	if err == sql.ErrNoRows {
		err = r.CreateCategoryByName(ctx, categoryName, userId)
		if err != nil {
			return "", fmt.Errorf("error creating category by name: %v", err)
		}
		// re-run select statement fetching newly created category
		err = stmt.QueryRowContext(ctx, sql.Named("name", categoryName), sql.Named("id", userId)).Scan(&cName)
		if err != nil {
			return "", fmt.Errorf("error fetching new category by name: %v", err)
		}
	} else if err != nil {
		return "", err
	}

	return cName, nil
}
*/

/*
// Deprecated
func (r *Repository) Post(_ context.Context, user *model.User) error {
	r.Lock()
	defer r.Unlock()
	db, err := MsSqlConnection()
	defer db.Close()
	if err != nil {
		return fmt.Errorf("error establishing DB connection: %v", err)
	}
	return nil
}
*/
