package db

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

type InsertBuilderWithDuplicateKeyUpdate struct {
	*sqlbuilder.InsertBuilder

	args        *sqlbuilder.Args
	assignments []string
}

func NewInsertBuilderWithDuplicateKeyUpdate() *InsertBuilderWithDuplicateKeyUpdate {
	args := &sqlbuilder.Args{}
	return &InsertBuilderWithDuplicateKeyUpdate{
		InsertBuilder: sqlbuilder.NewInsertBuilder(),
		args:          args,
	}
}

func (ib *InsertBuilderWithDuplicateKeyUpdate) Assign(field string, value interface{}) string {
	return fmt.Sprintf("%v = %v", sqlbuilder.Escape(field), ib.args.Add(value))
}

func (ib *InsertBuilderWithDuplicateKeyUpdate) OnDuplicateKeyUpdate(assignment ...string) *InsertBuilderWithDuplicateKeyUpdate {
	ib.assignments = append(ib.assignments, assignment...)
	return ib
}

func (ib *InsertBuilderWithDuplicateKeyUpdate) Build() (sql string, args []interface{}) {
	sql1, args1 := ib.InsertBuilder.Build()
	sql2, args2 := ib.args.Compile(strings.Join(ib.assignments, ", "))

	sql = fmt.Sprintf("%v ON DUPLICATE KEY UPDATE %v", sql1, sql2)
	args = make([]interface{}, 0, len(args1)+len(args2))
	args = append(args, args1...)
	args = append(args, args2...)
	return
}
