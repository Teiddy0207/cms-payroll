package service

import (
	"cal-salary/modules/payroll/entity"
	"strings"
	"testing"
)

func formula(variableName, expression string) entity.PayrollFormula {
	return entity.PayrollFormula{VariableName: variableName, Expression: expression}
}

func orderOf(formulas []entity.PayrollFormula) []string {
	names := make([]string, len(formulas))
	for i, f := range formulas {
		names[i] = f.VariableName
	}
	return names
}

func indexOf(names []string, name string) int {
	for i, n := range names {
		if n == name {
			return i
		}
	}
	return -1
}

func TestSortFormulas_LinearChain(t *testing.T) {
	formulas := []entity.PayrollFormula{
		formula("A", "B + 1"),
		formula("B", "C + 1"),
		formula("C", "1"),
	}

	sorted, err := sortFormulas(formulas)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	order := orderOf(sorted)
	if indexOf(order, "C") > indexOf(order, "B") || indexOf(order, "B") > indexOf(order, "A") {
		t.Fatalf("expected order C, B, A; got %v", order)
	}
}

func TestSortFormulas_Cycle(t *testing.T) {
	formulas := []entity.PayrollFormula{
		formula("A", "B + 1"),
		formula("B", "A + 1"),
	}

	_, err := sortFormulas(formulas)
	if err == nil {
		t.Fatal("expected circular dependency error, got nil")
	}
	if !strings.Contains(err.Error(), "circular") {
		t.Fatalf("expected error to mention 'circular', got: %v", err)
	}
}

func TestSortFormulas_NoDependencies(t *testing.T) {
	formulas := []entity.PayrollFormula{
		formula("BONUS", "100000"),
		formula("ALLOWANCE", "50000"),
	}

	sorted, err := sortFormulas(formulas)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sorted) != 2 {
		t.Fatalf("expected 2 formulas, got %d", len(sorted))
	}
}

func TestSortFormulas_RealisticCascade(t *testing.T) {
	formulas := []entity.PayrollFormula{
		formula("NET_SALARY", "GROSS_SALARY - TAX"),
		formula("TAX", "GROSS_SALARY * 0.1"),
		formula("GROSS_SALARY", "P1 + P2 + P3"),
		formula("P3", "P1 * 0.05"),
	}

	sorted, err := sortFormulas(formulas)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	order := orderOf(sorted)
	if indexOf(order, "P3") > indexOf(order, "GROSS_SALARY") {
		t.Fatalf("expected P3 before GROSS_SALARY; got %v", order)
	}
	if indexOf(order, "GROSS_SALARY") > indexOf(order, "TAX") {
		t.Fatalf("expected GROSS_SALARY before TAX; got %v", order)
	}
	if indexOf(order, "TAX") > indexOf(order, "NET_SALARY") {
		t.Fatalf("expected TAX before NET_SALARY; got %v", order)
	}
}

func TestExtractIdentifiers_Simple(t *testing.T) {
	idents, err := extractIdentifiers("P1 + P2 * 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !idents["P1"] || !idents["P2"] {
		t.Fatalf("expected P1 and P2 in identifiers, got %v", idents)
	}
	if len(idents) != 2 {
		t.Fatalf("expected exactly 2 identifiers, got %v", idents)
	}
}

func TestExtractIdentifiers_FunctionCall(t *testing.T) {
	idents, err := extractIdentifiers("round(GROSS_SALARY, 2)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !idents["GROSS_SALARY"] {
		t.Fatalf("expected GROSS_SALARY in identifiers, got %v", idents)
	}
}
