package scraper

import (
	"context"

	"github.com/chromedp/cdproto/dom"
)

func getAttributesFromDom(
	ctx context.Context,
	selector string,
	attributeName string,
) ([]string, error) {
	rootNode, err := dom.GetDocument().Do(ctx)
	if err != nil {
		return nil, err
	}
	nodes, err := dom.QuerySelectorAll(rootNode.NodeID, selector).Do(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(nodes))
	for _, node := range nodes {
		attributesRaw, err := dom.GetAttributes(node).Do(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, attributesToMap(attributesRaw)[attributeName])
	}

	return result, nil
}

func attributesToMap(rawAttributes []string) map[string]string {
	result := make(map[string]string, len(rawAttributes)/2)
	for i := 0; i < len(rawAttributes)/2; i++ {
		result[rawAttributes[2*i]] = rawAttributes[1+2*i]
	}
	return result
}
