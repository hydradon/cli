package results

import (
	"fmt"

	flinkgatewayv1beta1 "github.com/confluentinc/ccloud-sdk-go-v2/flink-gateway/v1beta1"

	"github.com/confluentinc/cli/v3/pkg/flink/types"
)

func convertToInternalField(field any, details flinkgatewayv1beta1.ColumnDetails) types.StatementResultField {
	converter := GetConverterForType(details.GetType())
	if converter != nil {
		return converter(field)
	}

	return types.AtomicStatementResultField{
		Type:  types.Null,
		Value: "NULL",
	}
}

func ConvertToInternalResults(results []any, resultSchema flinkgatewayv1beta1.SqlV1beta1ResultSchema) (*types.StatementResults, error) {
	headers := make([]string, len(resultSchema.GetColumns()))
	for idx, column := range resultSchema.GetColumns() {
		headers[idx] = column.GetName()
	}

	convertedResults := make([]types.StatementResultRow, len(results))
	for rowIdx, result := range results {
		resultItem, ok := result.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("given result item does not match op/row schema")
		}

		items, _ := resultItem["row"].([]any)
		if len(items) != len(resultSchema.GetColumns()) {
			return nil, fmt.Errorf("given result row does not match the provided schema")
		}

		convertedFields := make([]types.StatementResultField, len(items))
		for colIdx, field := range items {
			columnSchema := resultSchema.GetColumns()[colIdx]
			convertedFields[colIdx] = convertToInternalField(field, columnSchema)
		}

		op, _ := resultItem["op"].(float64)
		convertedResults[rowIdx] = types.StatementResultRow{
			Operation: types.StatementResultOperation(op),
			Fields:    convertedFields,
		}
	}
	return &types.StatementResults{
		Headers: headers,
		Rows:    convertedResults,
	}, nil
}
