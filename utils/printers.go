package utils

import (
	"csv_extractor/models"
	"fmt"
)

func PrintMap(m map[string]models.Expense) {
	fmt.Printf("\n- Lista de gastos: \n")
	for i, expense := range m {
		fmt.Printf("    %s: %.2f \n", i, expense.Value)
	}
}

func PrintTotal(m map[string]models.Expense) {
	var total = 0.0

	for _, expense := range m {
		total += expense.Value
	}

	fmt.Printf("\nTotal: %.2f \n\n", total)
}

func SumByCategory(m map[string]models.Expense) {
	groups := make(map[string]float64)

	for _, expense := range m {
		groups[expense.Category] += expense.Value
	}

	fmt.Printf("\n- Gastos por Categoria: \n")

	for category, value := range groups {
		fmt.Printf("    %s: %.2f \n", category, value)
	}

	fmt.Printf("\n")
}
