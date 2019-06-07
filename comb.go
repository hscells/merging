package merging

import (
	"sort"
)

// CombSUM sums an item's score from all lists where it was present.
type CombSUM struct {
	Normaliser
}

func (c CombSUM) Merge(itemsLists []Items) Items {
	var (
		unique Items
		i      int
	)

	// Sum item scores from all the item lists where it is present.
	seen := make(map[string]Item)
	for _, items := range itemsLists {
		for _, item := range items {
			if _, ok := seen[item.Id]; !ok {
				seen[item.Id] = Item{
					Id:    item.Id,
					Score: item.Score,
				}
			} else {
				s := seen[item.Id].Score
				seen[item.Id] = Item{
					Id:    item.Id,
					Score: item.Score + s,
				}
			}
		}
	}

	// Create a flat slice from the unique items.
	unique = make(Items, len(seen))
	for _, v := range seen {
		unique[i] = v
		i++
	}

	c.Init(unique)
	for i, item := range unique {
		unique[i].Score = c.Normalise(item)
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Score > unique[j].Score
	})

	return unique
}

// CombMNZ additionally multiplies the CombSUM score by the number of lists that contain that item.
type CombMNZ struct {
	Normaliser
}

func (c CombMNZ) Merge(itemsLists []Items) Items {
	// Compute the CombSUM score for each item.
	csum := CombSUM{Normaliser: c.Normaliser}
	its := csum.Merge(itemsLists)

	// Then, record how many times each item appears in each of the lists of items.
	k := make(map[string]float64)
	for _, items := range itemsLists {
		for _, item := range items {
			if _, ok := k[item.Id]; !ok {
				k[item.Id] = 1
			} else {
				k[item.Id]++
			}
		}
	}

	// Finally, multiply each item score by the number of times a list of items contained that item.
	for _, item := range its {
		item.Score *= k[item.Id]
	}

	sort.Slice(its, func(i, j int) bool {
		return its[i].Score > its[j].Score
	})

	return its
}
