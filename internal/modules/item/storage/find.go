package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ulibaysya/casenya/internal/modules/item"
	"github.com/ulibaysya/casenya/internal/modules/itemtag"
	"github.com/ulibaysya/casenya/pkg/db"
)

func (s *Storage) FindItem(ctx context.Context, id int64) (item.Item, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return item.Item{}, fmt.Errorf("beginning tx: %w", err)
	}
	defer db.RollbackTX(ctx, tx, s.logger)

	const metaSelectSQL = `
		WITH tags AS (
			SELECT
				r.item_id,
				k.name AS key,
				k.is_public,
				v.name AS value
			FROM item_tag_relations AS r
			JOIN item_tag_values AS v ON r.item_tag_value_id = v.id
			JOIN item_tag_keys AS k ON v.key_id = k.id
			ORDER BY k.id, v.id
		)
		SELECT
			i.name,
			i.image_url,
			t.key,
			t.is_public,
			t.value
		FROM items AS i
		LEFT JOIN tags AS t ON t.item_id = i.id
		WHERE i.id = $1;
	`
	metaSelectRows, err := tx.Query(ctx, metaSelectSQL, id)
	if err != nil {
		return item.Item{}, fmt.Errorf("querying item metadata: %w", err)
	}
	metaSelectResults, err := pgx.CollectRows(metaSelectRows, func(row pgx.CollectableRow) (metaSelectResult, error) {
		res := metaSelectResult{}
		err := row.Scan(&res.name, &res.imageURL, &res.key, &res.isPublic, &res.value)
		return res, err
	})
	if err != nil {
		return item.Item{}, fmt.Errorf("collecting item metadata: %w", err)
	}
	if len(metaSelectResults) == 0 {
		return item.Item{}, item.ErrNotFound
	}

	const pricesSelectSQL = `
		SELECT
			pp.code AS item_price_provider,
			p.price
		FROM items_prices AS p
		JOIN items_price_providers AS pp ON pp.id = p.provider_id
		WHERE p.item_id = $1;
	`
	priceSelectRows, err := tx.Query(ctx, pricesSelectSQL, id)
	if err != nil {
		return item.Item{}, fmt.Errorf("querying item prices: %w", err)
	}
	priceSelectResults, err := pgx.CollectRows(priceSelectRows, func(row pgx.CollectableRow) (priceSelectResult, error) {
		res := priceSelectResult{}
		err := row.Scan(&res.itemPriceProvider, &res.price)
		return res, err
	})
	if err != nil {
		return item.Item{}, fmt.Errorf("collecting item prices: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return item.Item{}, fmt.Errorf("commiting tx: %w", err)
	}

	return makeItem(id, metaSelectResults, priceSelectResults), nil
}

type metaSelectResult struct {
	name     string
	imageURL pgtype.Text
	key      pgtype.Text
	isPublic pgtype.Bool
	value    pgtype.Text
}

type priceSelectResult struct {
	itemPriceProvider string
	price             int64
}

func makeItem(id int64, metaResults []metaSelectResult, priceResults []priceSelectResult) item.Item {
	var imageURL string
	if metaResults[0].imageURL.Valid {
		imageURL = metaResults[0].imageURL.String
	} else {
		imageURL = ""
	}

	prices := make([]item.Price, 0, len(priceResults))
	for _, priceRes := range priceResults {
		prices = append(prices, item.Price{
			Provider: item.PriceProvider(priceRes.itemPriceProvider),
			Value:    priceRes.price,
		})
	}

	if !metaResults[0].key.Valid {
		return item.Item{
			ID:       id,
			Name:     metaResults[0].name,
			ImageURL: imageURL,
			Tags:     []itemtag.Tag{},
			Prices:   prices,
		}
	}

	tags := []itemtag.Tag{}
	tag := itemtag.Tag{
		Key:      metaResults[0].key.String,
		IsPublic: metaResults[0].isPublic.Bool,
		Values:   []string{metaResults[0].value.String},
	}
	for _, metaRes := range metaResults[1:] {
		if tag.Key != metaRes.key.String {
			tags = append(tags, tag)
			tag = itemtag.Tag{
				Key:      metaRes.key.String,
				IsPublic: metaRes.isPublic.Bool,
			}
		}
		tag.Values = append(tag.Values, metaRes.value.String)
	}
	tags = append(tags, tag)

	return item.Item{
		ID:       id,
		Name:     metaResults[0].name,
		ImageURL: imageURL,
		Tags:     tags,
		Prices:   prices,
	}
}
