package partition_query_loader

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
var numberPattern = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?$`)

// PartitionDefinition is the supported PARTITION ON COLUMNS metadata from SHOW CREATE TABLE.
type PartitionDefinition struct {
	Columns    []string
	Partitions []Partition
}

// Partition is one partition's conjunction of supported column conditions.
type Partition struct {
	Conditions []PartitionCondition
}

// PartitionCondition compares a partition column with a bound query argument.
type PartitionCondition struct {
	Column   string
	Operator string
	Value    any
}

// Predicate is a parameterized SQL predicate that can safely be appended to a fixed query template.
type Predicate struct {
	SQL  string
	Args []any
}

// ParsePartitionDefinition extracts the supported PARTITION ON COLUMNS definition from SHOW CREATE TABLE output.
func ParsePartitionDefinition(showCreate string) (PartitionDefinition, error) {
	columnsText, partitionsText, err := partitionSections(showCreate)
	if err != nil {
		return PartitionDefinition{}, err
	}
	columns, err := parseColumns(columnsText)
	if err != nil {
		return PartitionDefinition{}, err
	}
	partitionTexts, err := splitTopLevel(partitionsText, ',')
	if err != nil || len(partitionTexts) == 0 {
		return PartitionDefinition{}, fmt.Errorf("unsupported partition definition: missing partition conditions")
	}
	partitions := make([]Partition, 0, len(partitionTexts))
	for _, partitionText := range partitionTexts {
		partition, err := parsePartition(partitionText, columns)
		if err != nil {
			return PartitionDefinition{}, err
		}
		partitions = append(partitions, partition)
	}
	return PartitionDefinition{Columns: columns, Partitions: partitions}, nil
}

// BuildPartitionPredicate returns a parameterized predicate for one parsed partition.
func BuildPartitionPredicate(partition Partition) (Predicate, error) {
	if len(partition.Conditions) == 0 {
		return Predicate{}, fmt.Errorf("unsupported partition definition: empty partition predicate")
	}
	terms := make([]string, 0, len(partition.Conditions))
	args := make([]any, 0, len(partition.Conditions))
	for _, condition := range partition.Conditions {
		if !identifierPattern.MatchString(condition.Column) || !supportedOperator(condition.Operator) {
			return Predicate{}, fmt.Errorf("unsupported partition definition: invalid condition")
		}
		terms = append(terms, quoteIdentifier(condition.Column)+" "+condition.Operator+" ?")
		args = append(args, condition.Value)
	}
	return Predicate{SQL: strings.Join(terms, " AND "), Args: args}, nil
}

func partitionSections(showCreate string) (string, string, error) {
	start := keywordIndex(showCreate, "PARTITION")
	if start < 0 {
		return "", "", fmt.Errorf("unsupported partition definition: PARTITION ON COLUMNS not found")
	}
	rest := strings.TrimSpace(showCreate[start+len("PARTITION"):])
	if !hasPrefixKeyword(rest, "ON") {
		return "", "", fmt.Errorf("unsupported partition definition: expected ON")
	}
	rest = strings.TrimSpace(rest[len("ON"):])
	if !hasPrefixKeyword(rest, "COLUMNS") {
		return "", "", fmt.Errorf("unsupported partition definition: expected COLUMNS")
	}
	rest = strings.TrimSpace(rest[len("COLUMNS"):])
	columns, rest, err := takeParenthesized(rest)
	if err != nil {
		return "", "", fmt.Errorf("unsupported partition definition: invalid partition columns: %w", err)
	}
	partitions, _, err := takeParenthesized(strings.TrimSpace(rest))
	if err != nil {
		return "", "", fmt.Errorf("unsupported partition definition: invalid partition conditions: %w", err)
	}
	return columns, partitions, nil
}

func parseColumns(text string) ([]string, error) {
	parts, err := splitTopLevel(text, ',')
	if err != nil || len(parts) == 0 {
		return nil, fmt.Errorf("unsupported partition definition: invalid partition columns")
	}
	columns := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		column, err := parseIdentifier(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("unsupported partition definition: invalid partition column: %w", err)
		}
		if _, ok := seen[column]; ok {
			return nil, fmt.Errorf("unsupported partition definition: duplicate partition column %q", column)
		}
		seen[column] = struct{}{}
		columns = append(columns, column)
	}
	return columns, nil
}

func parsePartition(text string, columns []string) (Partition, error) {
	text = unwrapCondition(text)
	conditionsText, err := splitConditions(text)
	if err != nil || len(conditionsText) == 0 {
		return Partition{}, fmt.Errorf("unsupported partition definition: invalid partition condition")
	}
	allowed := make(map[string]struct{}, len(columns))
	for _, column := range columns {
		allowed[column] = struct{}{}
	}
	covered := make(map[string]struct{}, len(columns))
	conditions := make([]PartitionCondition, 0, len(conditionsText))
	for _, conditionText := range conditionsText {
		condition, err := parseCondition(conditionText)
		if err != nil {
			return Partition{}, err
		}
		if _, ok := allowed[condition.Column]; !ok {
			return Partition{}, fmt.Errorf("unsupported partition definition: column %q is not partitioned", condition.Column)
		}
		covered[condition.Column] = struct{}{}
		conditions = append(conditions, condition)
	}
	if len(covered) != len(columns) {
		return Partition{}, fmt.Errorf("unsupported partition definition: condition does not cover every partition column")
	}
	return Partition{Conditions: conditions}, nil
}

func parseCondition(text string) (PartitionCondition, error) {
	text = unwrapCondition(text)
	for _, operator := range []string{">=", "<=", "=", "<", ">"} {
		if index := strings.Index(text, operator); index >= 0 {
			column, err := parseIdentifier(strings.TrimSpace(text[:index]))
			if err != nil {
				return PartitionCondition{}, fmt.Errorf("unsupported partition definition: invalid condition column: %w", err)
			}
			value, err := parseLiteral(strings.TrimSpace(text[index+len(operator):]))
			if err != nil {
				return PartitionCondition{}, err
			}
			return PartitionCondition{Column: column, Operator: operator, Value: value}, nil
		}
	}
	return PartitionCondition{}, fmt.Errorf("unsupported partition definition: condition requires a supported comparison")
}

// unwrapCondition accepts grouping parentheses around an otherwise supported condition.
func unwrapCondition(text string) string {
	text = strings.TrimSpace(text)
	for strings.HasPrefix(text, "(") {
		inner, rest, err := takeParenthesized(text)
		if err != nil || strings.TrimSpace(rest) != "" {
			return text
		}
		text = strings.TrimSpace(inner)
	}
	return text
}

func parseLiteral(text string) (any, error) {
	if hasPrefixKeyword(text, "TIMESTAMP") {
		value, err := parseQuoted(strings.TrimSpace(text[len("TIMESTAMP"):]))
		if err != nil {
			return nil, fmt.Errorf("unsupported partition definition: invalid timestamp literal")
		}
		for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999999", "2006-01-02 15:04:05", "2006-01-02"} {
			if timestamp, err := time.Parse(layout, value); err == nil {
				return timestamp, nil
			}
		}
		return nil, fmt.Errorf("unsupported partition definition: invalid timestamp literal")
	}
	if strings.HasPrefix(text, "'") {
		value, err := parseQuoted(text)
		if err != nil {
			return nil, fmt.Errorf("unsupported partition definition: invalid string literal")
		}
		return value, nil
	}
	if !numberPattern.MatchString(text) {
		return nil, fmt.Errorf("unsupported partition definition: literal must be numeric, string, or timestamp")
	}
	// Keep the original lexical representation. Converting through float64 loses
	// precision for large integers and decimal partition bounds.
	return text, nil
}

func parseQuoted(text string) (string, error) {
	if len(text) < 2 || text[0] != '\'' || text[len(text)-1] != '\'' {
		return "", fmt.Errorf("expected quoted literal")
	}
	value := text[1 : len(text)-1]
	for index := 0; index < len(value); index++ {
		if value[index] != '\'' {
			continue
		}
		if index+1 >= len(value) || value[index+1] != '\'' {
			return "", fmt.Errorf("unescaped quote")
		}
		index++
	}
	return strings.ReplaceAll(value, "''", "'"), nil
}

func splitConditions(text string) ([]string, error) {
	parts := make([]string, 0, 2)
	start := 0
	inQuote := false
	for index := 0; index < len(text); index++ {
		if text[index] == '\'' {
			if inQuote && index+1 < len(text) && text[index+1] == '\'' {
				index++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if !inQuote && hasKeywordAt(text, index, "AND") {
			parts = append(parts, strings.TrimSpace(text[start:index]))
			index += len("AND") - 1
			start = index + 1
		}
	}
	if inQuote {
		return nil, fmt.Errorf("unterminated string")
	}
	parts = append(parts, strings.TrimSpace(text[start:]))
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("empty condition")
		}
	}
	return parts, nil
}

func splitTopLevel(text string, separator byte) ([]string, error) {
	parts := make([]string, 0, 1)
	start, depth := 0, 0
	inQuote := false
	for index := 0; index < len(text); index++ {
		switch text[index] {
		case '\'':
			if inQuote && index+1 < len(text) && text[index+1] == '\'' {
				index++
				continue
			}
			inQuote = !inQuote
		case '(':
			if !inQuote {
				depth++
			}
		case ')':
			if !inQuote {
				depth--
				if depth < 0 {
					return nil, fmt.Errorf("unbalanced parentheses")
				}
			}
		default:
			if !inQuote && depth == 0 && text[index] == separator {
				parts = append(parts, strings.TrimSpace(text[start:index]))
				start = index + 1
			}
		}
	}
	if inQuote || depth != 0 {
		return nil, fmt.Errorf("unbalanced partition text")
	}
	parts = append(parts, strings.TrimSpace(text[start:]))
	return parts, nil
}

func takeParenthesized(text string) (string, string, error) {
	if !strings.HasPrefix(text, "(") {
		return "", "", fmt.Errorf("expected opening parenthesis")
	}
	depth, inQuote := 0, false
	for index := 0; index < len(text); index++ {
		if text[index] == '\'' {
			if inQuote && index+1 < len(text) && text[index+1] == '\'' {
				index++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if inQuote {
			continue
		}
		switch text[index] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return text[1:index], text[index+1:], nil
			}
		}
	}
	return "", "", fmt.Errorf("unbalanced parentheses")
}

func parseIdentifier(text string) (string, error) {
	text = strings.TrimSpace(text)
	if len(text) >= 2 && text[0] == '`' && text[len(text)-1] == '`' {
		text = text[1 : len(text)-1]
	}
	if !identifierPattern.MatchString(text) {
		return "", fmt.Errorf("invalid identifier %q", text)
	}
	return text, nil
}

func keywordIndex(text, keyword string) int {
	for index := 0; index <= len(text)-len(keyword); index++ {
		if hasKeywordAt(text, index, keyword) {
			return index
		}
	}
	return -1
}

func hasPrefixKeyword(text, keyword string) bool {
	return len(text) >= len(keyword) && hasKeywordAt(text, 0, keyword)
}

func hasKeywordAt(text string, index int, keyword string) bool {
	end := index + len(keyword)
	if end > len(text) || !strings.EqualFold(text[index:end], keyword) {
		return false
	}
	return (index == 0 || !isIdentifierCharacter(text[index-1])) &&
		(end == len(text) || !isIdentifierCharacter(text[end]))
}

func isIdentifierCharacter(char byte) bool {
	return char == '_' || char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9'
}

func supportedOperator(operator string) bool {
	return operator == "=" || operator == "<" || operator == "<=" || operator == ">" || operator == ">="
}

func quoteIdentifier(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}
