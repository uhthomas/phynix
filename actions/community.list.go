package actions

import (
	"math"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
	"sort"
)

func CommunityList(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	communities := realtime.NewCommunitySlice()

	sort.Sort(communities)

	communities = communities[:int(math.Min(float64(len(communities)), 50))]

	var results []models.Community

	if err := models.DB.Where(func() (payload []uint64) {
		for _, c := range communities {
			payload = append(payload, c.ID)
		}
		return
	}()).Preload("User").Find(&results).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityList
	}

	var payload []models.Community

	for _, c := range communities {
		for i := range results {
			r := &results[i]
			if c.ID == r.ID {
				r.Population = len(c.Users)
				r.Waitlist = len(c.Waitlist)
				if m := c.Media; m != nil {
					r.Media = m.Item
				}

				payload = append(payload, *r)
				break
			}
		}
	}

	// for i := range results {
	// 	r := &results[i]
	// 	c := communities[i]
	// 	r.Population = len(c.Users)
	// 	fmt.Println(len(c.Users))
	// 	r.Waitlist = len(c.Waitlist)
	// 	if m := c.Media; m != nil {
	// 		r.Media = &m.Item
	// 	}
	// }

	return int(enum.ResponseCodeOk), payload, nil
}
