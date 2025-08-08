package pg

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func PrintPgError(err *pgconn.PgError) {

	fmt.Printf("--------------------------------------------------------------------\n")
	fmt.Printf("PGERROR: Severity: %v\n", err.Severity)
	fmt.Printf("PGERROR: SeverityUnlocalized: %v\n", err.SeverityUnlocalized)
	fmt.Printf("PGERROR: Code: %v\n", err.Code)
	fmt.Printf("PGERROR: Message: %v\n", err.Message)
	fmt.Printf("PGERROR: Detail: %v\n", err.Detail)
	fmt.Printf("PGERROR: Hint: %v\n", err.Hint)
	fmt.Printf("PGERROR: Position: %v\n", err.Position)
	fmt.Printf("PGERROR: InternalPosition: %v\n", err.InternalPosition)
	fmt.Printf("PGERROR: InternalQuery: %v\n", err.InternalQuery)
	fmt.Printf("PGERROR: Where: %v\n", err.Where)
	fmt.Printf("PGERROR: SchemaName: %v\n", err.SchemaName)
	fmt.Printf("PGERROR: TableName: %v\n", err.TableName)
	fmt.Printf("PGERROR: ColumnName: %v\n", err.ColumnName)
	fmt.Printf("PGERROR: DataTypeName: %v\n", err.DataTypeName)
	fmt.Printf("PGERROR: ConstraintName: %v\n", err.ConstraintName)
	fmt.Printf("PGERROR: File: %v\n", err.File)
	fmt.Printf("PGERROR: Line: %v\n", err.Line)
	fmt.Printf("PGERROR: Routine: %v\n", err.Routine)
	fmt.Printf("--------------------------------------------------------------------\n")
}
