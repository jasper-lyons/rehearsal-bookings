package sql

import (
	"testing"
	"strings"
	"fmt"
	"github.com/andreyvit/diff"
)

type Node struct {
	Type string
	Children []Node
	Value string
}

func Map[A any, R any](array []A, handler func (A) R) []R {
	var result []R

	for _, element := range array {
		result = append(result, handler(element))
	}

	return result 
}

func Statement(nodes ...Node) Node {
	return Node {
		Type: "statement",
		Children: nodes,
	}
}

func Select(nodes ...Node) Node {
	return Node {
		Type: "select",
		Children: nodes,
	}
}

func All() Node {
	return Node { Type: "all", Children: nil }
}

func From(nodes ...Node) Node {
	return Node {
		Type: "from",
		Children: nodes,
	}
}

func Table(name string) Node {
	return Node {
		Type: "table",
		Value: name,
	}
}

func CreateTable(name string, children ...Node) Node {
	return Node {
		Type: "create",
		Value: name,
		Children: children,
	}
}

func ColumnDefinition(name string, children ...Node) Node {
	return Node {
		Type: "column_definition",
		Value: name,
		Children: children,
	}
}

func Integer() Node { return Node { Type: "type", Value: "Integer" } }

func Annotation(value string) Node {
	return Node {
		Type: "annotation",
		Value: value,
	}
}

func Datetime() Node {
	return Node {
		Type: "type",
		Value: "Datetime",
	}
}

func Text() Node {
	return Node {
		Type: "type",
		Value: "Text",
	}
}

func Constraint(value string) Node {
	return Node {
		Type: "constraint",
		Value: value,
	}
}

func Default(node Node) Node {
	return Node {
		Type: "default",
		Children: []Node{node},
	}
}

func Now() Node {
	return Node {
		Type: "now",
	}
}

func Check(node Node) Node {
	return Node {
		Type: "check",
		Children: []Node{node},
	}
}

func GT(left Node, right Node) Node {
	return Node {
		Type: "gt",
		Children: []Node{left, right},
	}
}

func Column(name string) Node {
	return Node {
		Type: "column",
		Value: name,
	}
}

func Trigger(name string, timing string, event string, table string, statements ...Node) Node {
	return Node {
		Type: "trigger",
		Value: name,
		Children: append([]Node{
			Node{Type: "trigger_timing", Value: timing},
			Node{Type: "trigger_event", Value: event},
			Node{Type: "trigger_table", Value: table},
		}, statements...),
	}
}

func Update(table string, assignments ...Node) Node {
	return Node {
		Type: "update",
		Value: table,
		Children: assignments,
	}
}

func Set(assignments ...Node) Node {
	return Node {
		Type: "set",
		Children: assignments,
	}
}

func Assignment(column string, value Node) Node {
	return Node {
		Type: "assignment",
		Value: column,
		Children: []Node{value},
	}
}

func Where(condition Node) Node {
	return Node {
		Type: "where",
		Children: []Node{condition},
	}
}

func Eq(left Node, right Node) Node {
	return Node {
		Type: "eq",
		Children: []Node{left, right},
	}
}

func DatetimeFunc(value string) Node {
	return Node {
		Type: "datetime_func",
		Value: value,
	}
}

func New(column string) Node {
	return Node {
		Type: "new",
		Value: column,
	}
}

func Begin(statements ...Node) Node {
	return Node {
		Type: "begin",
		Children: statements,
	}
}

func End() Node {
	return Node {
		Type: "end",
	}
}

func CreateIndex(name string, table string, columns ...Node) Node {
	return Node {
		Type: "create_index",
		Value: name,
		Children: append([]Node{Table(table)}, columns...),
	}
}

func Sqlite(node Node) string {
	switch node.Type {
	case "statement":
		return strings.Join(Map(node.Children, Sqlite), " ") + ";"
	case "select": 
		return "select " + strings.Join(Map(node.Children, Sqlite), " ,")
	case "all": 
		return "*"
	case "from": 
		return "from " + strings.Join(Map(node.Children, Sqlite), " ,")
	case "table":
		return node.Value
	case "create":
		return fmt.Sprintf("create table %s (\n%s\n)", 
			node.Value,
			strings.Join(Map(node.Children, func(n Node) string {
				return "\t" + Sqlite(n)
			}), ",\n"),
		)
	case "column_definition":
		parts := []string{node.Value}
		for _, child := range node.Children {
			parts = append(parts, Sqlite(child))
		}
		return strings.Join(parts, " ")
	case "type":
		switch node.Value {
		case "Integer":
			return "integer"
		case "Datetime":
			return "datetime"
		case "Text":
			return "text"
		}
	case "annotation":
		if node.Value == "primary key" {
			return "primary key autoincrement"
		}
		return node.Value
	case "constraint":
		return node.Value
	case "default":
		return "default " + Sqlite(node.Children[0])
	case "now":
		return "current_timestamp"
	case "check":
		return "check (" + Sqlite(node.Children[0]) + ")"
	case "gt":
		return Sqlite(node.Children[0]) + " > " + Sqlite(node.Children[1])
	case "column":
		return node.Value
	case "create_index":
		columns := strings.Join(Map(node.Children[1:], Sqlite), ", ")
		return fmt.Sprintf("create index %s on %s(%s)", node.Value, Sqlite(node.Children[0]), columns)
	case "trigger":
		parts := []string{
			fmt.Sprintf("create trigger %s", node.Value),
			fmt.Sprintf("    %s %s on %s", 
				strings.ToLower(Sqlite(node.Children[0])),
				strings.ToLower(Sqlite(node.Children[1])),
				Sqlite(node.Children[2])),
		}
		parts = append(parts, Map(node.Children[3:], Sqlite)...)
		return strings.Join(parts, "\n")
	case "begin":
		stmts := Map(node.Children, func(n Node) string {
			return "        " + Sqlite(n)
		})
		return "    begin\n" + strings.Join(stmts, "\n") + "\n    end"
	case "update":
		parts := []string{"update", node.Value}
		for _, child := range node.Children {
			parts = append(parts, Sqlite(child))
		}
		return strings.Join(parts, " ") + ";"
	case "set":
		return "set " + strings.Join(Map(node.Children, Sqlite), ", ")
	case "assignment":
		return node.Value + " = " + Sqlite(node.Children[0])
	case "where":
		return "where " + Sqlite(node.Children[0])
	case "eq":
		return Sqlite(node.Children[0]) + " = " + Sqlite(node.Children[1])
	case "datetime_func":
		return "datetime('" + node.Value + "')"
	case "new":
		return "new." + node.Value
	case "end":
		return "end"
	}
	return ""
}

func check(t *testing.T, output string, expected string) {
	if output != expected {
		t.Fatalf("%v", diff.LineDiff(expected, output))
	}
}

func TestGenerateSql(t *testing.T) {
	check(
		t,
		Sqlite(Statement(Select(All()), From(Table("bookings")))),
		"select * from bookings;",
	)

	check(
		t,
		Sqlite(Statement(CreateTable("bookings",
			ColumnDefinition("id", Integer(), Annotation("primary key")),
			ColumnDefinition("created_at", Datetime(), Constraint("not null"), Default(Now())),
			ColumnDefinition("updated_at", Datetime(), Constraint("not null"), Default(Now())),
			ColumnDefinition("customer_name", Text(), Constraint("not null")),
			ColumnDefinition("customer_email", Text(), Constraint("not null")),
			ColumnDefinition("customer_phone", Text()),
			ColumnDefinition("room_name", Text(), Constraint("not null")),
			ColumnDefinition("start_time", Datetime(), Constraint("not null")),
			ColumnDefinition("end_time", Datetime(), Constraint("not null")),
			ColumnDefinition("notes", Text()),
			Check(GT(Column("end_time"), Column("start_time"))),
		))),
`create table bookings (
	id integer primary key autoincrement,
	created_at datetime not null default current_timestamp,
	updated_at datetime not null default current_timestamp,
	customer_name text not null,
	customer_email text not null,
	customer_phone text,
	room_name text not null,
	start_time datetime not null,
	end_time datetime not null,
	notes text,
	check (end_time > start_time)
);`,
	)

	check(
			t,
			Sqlite(Statement(CreateIndex("idx_bookings_room_times", "bookings", 
					Column("room_name"), Column("start_time"), Column("end_time")))),
			`create index idx_bookings_room_times on bookings(room_name, start_time, end_time);`,
	)

	check(
			t,
			Sqlite(Statement(
					Trigger("update_bookings_updated_at", "after", "update", "bookings",
							Begin(
									Update("bookings",
											Set(Assignment("updated_at", DatetimeFunc("now"))),
											Where(Eq(Column("id"), New("id"))),
									),
							),
					),
			)),
`create trigger update_bookings_updated_at after update
		on
    begin
        update bookings set updated_at = datetime('now') where id = new.id;
    end;`,
	)
}
