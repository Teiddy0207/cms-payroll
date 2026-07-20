package service

import (
	"cal-salary/modules/payroll/entity"
	"fmt"

	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
)

// identifierVisitor collects every identifier name referenced in an expression.
type identifierVisitor struct {
	identifiers map[string]bool
}

func (v *identifierVisitor) Visit(node *ast.Node) {
	if id, ok := (*node).(*ast.IdentifierNode); ok {
		v.identifiers[id.Value] = true
	}
}

// extractIdentifiers returns the set of bare identifier names referenced in expression.
// It includes function-call names (e.g. "round" in "round(x)") — harmless, since a
// dependency edge is only created when a name matches another formula's VariableName.
func extractIdentifiers(expression string) (map[string]bool, error) {
	tree, err := parser.Parse(expression)
	if err != nil {
		return nil, err
	}

	visitor := &identifierVisitor{identifiers: map[string]bool{}}
	ast.Walk(&tree.Node, visitor)
	return visitor.identifiers, nil
}

// sortFormulas orders formulas so that every formula runs after the formulas it
// depends on (identified by referencing their VariableName in its Expression),
// using Kahn's algorithm. Returns an error if a circular dependency is found.
func sortFormulas(formulas []entity.PayrollFormula) ([]entity.PayrollFormula, error) {
	varToIdx := make(map[string]int, len(formulas))
	for i, f := range formulas {
		varToIdx[f.VariableName] = i
	}

	indegree := make([]int, len(formulas))
	dependents := make([][]int, len(formulas)) // dependents[j] = formulas that depend on formula j

	for i, f := range formulas {
		idents, err := extractIdentifiers(f.Expression)
		if err != nil {
			return nil, fmt.Errorf("formula %s: invalid expression: %w", f.VariableName, err)
		}
		for name := range idents {
			j, ok := varToIdx[name]
			if !ok || j == i {
				continue
			}
			dependents[j] = append(dependents[j], i)
			indegree[i]++
		}
	}

	queue := make([]int, 0, len(formulas))
	for i := range formulas {
		if indegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	order := make([]int, 0, len(formulas))
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)
		for _, dep := range dependents[n] {
			indegree[dep]--
			if indegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	if len(order) != len(formulas) {
		return nil, fmt.Errorf("circular dependency detected among payroll formulas")
	}

	result := make([]entity.PayrollFormula, len(formulas))
	for pos, idx := range order {
		result[pos] = formulas[idx]
	}
	return result, nil
}
