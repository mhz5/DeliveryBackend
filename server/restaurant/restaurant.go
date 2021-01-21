package restaurant

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/type/latlng"

	cpb "fda/proto/common"
	rpb "fda/proto/restaurant"
)


type restaurant struct {
	Id       int64                  `firestore:"Id"`
	Name     string                 `firestore:"Name"`
	Location *latlng.LatLng         `firestore:"Location"`
	Menu     *firestore.DocumentRef `firestore:"Menu"`
}

type menu struct {
	Items []menuItem `firestore:"items"`
}

type menuItem struct {
	Name string   `firestore:"name"`
	Price float32 `firestore:"price"`
}

func parseRestaurant(ctx context.Context, doc *firestore.DocumentSnapshot) (*rpb.Restaurant, error) {
	var data restaurant
	if err := doc.DataTo(&data); err != nil {
		return nil, errors.Wrap(err, "could not parse restaurant")
	}

	restaurant := &rpb.Restaurant{
		Id: data.Id,
		Name: data.Name,
		Location: &cpb.Point{
			Lat:  data.Location.Latitude,
			Long: data.Location.Longitude,
		},
		Menu: &rpb.Menu{
			Items: []*rpb.MenuItem{},
		},
	}

	menuData, err := data.Menu.Get(ctx)
	if err != nil {
		return restaurant, errors.Wrapf(err, "Could not fetch menu for restaurant with Id %d", data.Id)
	}
	menu, err := parseMenu(menuData)
	if err != nil {
		return restaurant, errors.Wrapf(err, "Could not parse menu for restaurant with Id %d", data.Id)
	}
	restaurant.Menu = menu

	return restaurant, nil
}

func parseMenu(m *firestore.DocumentSnapshot) (*rpb.Menu, error) {
	var data menu
	if err := m.DataTo(&data); err != nil {
		return nil, errors.Wrap(err, "Could not parse menu data")
	}
	menu := &rpb.Menu{
		Items: []*rpb.MenuItem{},
	}
	for _, item := range data.Items {
		menu.Items = append(menu.Items, &rpb.MenuItem{
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return menu, nil
}
