package storage

import (
	"slices"
	"strings"
	"testing"

	"github.com/ulibaysya/casenya/internal/modules/item"
)

func formatSQL(sql string) string {
	sql = strings.Join(strings.Fields(strings.TrimSpace(sql)), " ")
	sql = strings.ReplaceAll(sql, " ,", ",")
	sql = strings.ReplaceAll(sql, ", ", ",")
	return sql
}

func TestBuildTagsInsertSQL(t *testing.T) {
	tags := []item.Tag{
		{
			Key:    "1",
			Values: []string{"1.1"},
		},
		{
			Key:    "2",
			Values: []string{"2.1", "2.2"},
		},
		{
			Key:    "3",
			Values: []string{"3.1", "3.2", "3.3"},
		},
		{
			Key:    "4",
			Values: []string{"4.1", "4.2", "4.3", "4.4"},
		},
		{
			Key:    "5",
			Values: []string{"5.1", "5.2", "5.3", "5.4", "5.5"},
		},
	}
	someID := int64(57183)

	tests := []struct {
		name         string
		tags         []item.Tag
		id           int64
		expectedSQL  string
		expectedArgs []any
	}{
		{
			name: "singleKeySingleValue",
			tags: []item.Tag{tags[0]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $3, v.id
				FROM values_ids AS v;
			`,
			id:           someID,
			expectedArgs: []any{tags[0].Key, tags[0].Values[0], someID},
		},
		{
			name: "singleKeyTwoValues",
			tags: []item.Tag{tags[1]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($1, $3)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $4, v.id
				FROM values_ids AS v;
			`,
			id:           someID,
			expectedArgs: []any{tags[1].Key, tags[1].Values[0], tags[1].Values[1], someID},
		},
		{
			name: "singleKeyManyValues",
			tags: []item.Tag{tags[4]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($1, $3),
								($1, $4),
								($1, $5),
								($1, $6)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $7, v.id
				FROM values_ids AS v;
			`,
			id: someID,
			expectedArgs: []any{
				tags[4].Key,
				tags[4].Values[0],
				tags[4].Values[1],
				tags[4].Values[2],
				tags[4].Values[3],
				tags[4].Values[4],
				someID},
		},
		{
			name: "twoKeysOneValue",
			tags: []item.Tag{tags[0], tags[0]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($3, $4)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $5, v.id
				FROM values_ids AS v;
			`,
			id: someID,
			expectedArgs: []any{
				tags[0].Key,
				tags[0].Values[0],
				tags[0].Key,
				tags[0].Values[0],
				someID},
		},
		{
			name: "twoKeysOneMultipleValues",
			tags: []item.Tag{tags[0], tags[2]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($3, $4),
								($3, $5),
								($3, $6)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $7, v.id
				FROM values_ids AS v;
			`,
			id: someID,
			expectedArgs: []any{
				tags[0].Key,
				tags[0].Values[0],
				tags[2].Key,
				tags[2].Values[0],
				tags[2].Values[1],
				tags[2].Values[2],
				someID},
		},
		{
			name: "twoKeysMultipleOneValues",
			tags: []item.Tag{tags[2], tags[0]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($1, $3),
								($1, $4),
								($5, $6)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $7, v.id
				FROM values_ids AS v;
			`,
			id: someID,
			expectedArgs: []any{
				tags[2].Key,
				tags[2].Values[0],
				tags[2].Values[1],
				tags[2].Values[2],
				tags[0].Key,
				tags[0].Values[0],
				someID},
		},
		{
			name: "twoKeysMultipleMultipleValues",
			tags: []item.Tag{tags[2], tags[3]},
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($1, $3),
								($1, $4),
								($5, $6),
								($5, $7),
								($5, $8),
								($5, $9)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $10, v.id
				FROM values_ids AS v;
			`,
			id: someID,
			expectedArgs: []any{
				tags[2].Key,
				tags[2].Values[0],
				tags[2].Values[1],
				tags[2].Values[2],
				tags[3].Key,
				tags[3].Values[0],
				tags[3].Values[1],
				tags[3].Values[2],
				tags[3].Values[3],
				someID,
			},
		},
		{
			name: "multipleKeysMultipleValues",
			tags: tags,
			expectedSQL: `
				WITH values_ids AS (
					SELECT v.id AS id
					FROM item_tag_values AS v
					JOIN item_tag_keys AS k ON k.id = v.key_id
					JOIN (
							VALUES
								($1, $2),
								($3, $4), ($3, $5),
								($6, $7), ($6, $8), ($6, $9),
								($10, $11), ($10, $12), ($10, $13), ($10, $14),
								($15, $16), ($15, $17), ($15, $18), ($15, $19), ($15, $20)
					) AS n(key_name, value_name)
					ON n.key_name = k.name AND n.value_name = v.name
				)
				INSERT INTO item_tag_relations
				SELECT $21, v.id
				FROM values_ids AS v;
			`,
			id: someID,
			expectedArgs: []any{
				tags[0].Key, tags[0].Values[0],
				tags[1].Key, tags[1].Values[0], tags[1].Values[1],
				tags[2].Key, tags[2].Values[0], tags[2].Values[1], tags[2].Values[2],
				tags[3].Key, tags[3].Values[0], tags[3].Values[1], tags[3].Values[2], tags[3].Values[3],
				tags[4].Key, tags[4].Values[0], tags[4].Values[1], tags[4].Values[2], tags[4].Values[3], tags[4].Values[4],
				someID,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sql, args := buildTagsInsertSQL(test.id, test.tags)

			if !slices.Equal(test.expectedArgs, args) {
				t.Fatalf("slices are not equal:\n\t expected %v\n\t got %v", test.expectedArgs, args)
			}

			formattedExpectedSQL := formatSQL(test.expectedSQL)
			formattedSQL := formatSQL(sql)
			if formattedSQL != formattedExpectedSQL {
				t.Fatalf("got unexpected sql:\n\t expected %v\n\t got %v", formattedExpectedSQL, formattedSQL)
			}
		})
	}
}
