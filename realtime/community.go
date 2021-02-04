package realtime

import (
	"phynix/enum"
	"phynix/models"
	"pipeline"
	"raiki"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

type Community struct {
	sync.RWMutex
	ID       uint64
	Users    []uint64
	Waitlist []uint64
	Media    *CommunityMedia
	Timer    *time.Timer
	Pipe     *pipeline.Pipeline
	Server   *raiki.Server
	Room     string
	Skipped  bool
}

type CommunityMedia struct {
	UserID    uint64              `json:"userID"`
	Elapsed   float64             `json:"elapsed,omitempty"`
	Item      models.PlaylistItem `json:"item"`
	Timestamp time.Time           `json:"timestamp"`
	Votes     CommunityVotes      `json:"votes"`
}

type CommunityVotes struct {
	Woots []uint64 `json:"woots"`
	Saves []uint64 `json:"saves"`
	Mehs  []uint64 `json:"mehs"`
}

type CommunityState struct {
	Chat     []models.Chat    `json:"chat"`
	Users    []models.User    `json:"users"`
	Waitlist []uint64         `json:"waitlist"`
	Media    *CommunityMedia  `json:"media"`
	Meta     models.Community `json:"meta"`
}

type CommunityAdvance struct {
	History *models.History `json:"history,omitempty"`
	Media   *CommunityMedia `json:"media,omitempty"`
}

type CommunitySlice []*Community

func NewCommunitySlice() (payload CommunitySlice) {
	for _, c := range Communities {
		payload = append(payload, c)
	}
	return
}

func (c CommunitySlice) Len() int {
	return len(c)
}

func (c CommunitySlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CommunitySlice) Less(i, j int) bool {
	return len(c[i].Users) > len(c[j].Users)
}

var Communities = make(map[uint64]*Community)

func NewCommunity(id uint64, server *raiki.Server) *Community {
	if c, ok := Communities[id]; ok {
		return c
	}

	c := &Community{
		ID:       id,
		Users:    []uint64{},
		Waitlist: []uint64{},
		Timer:    time.NewTimer(0),
		Pipe:     pipeline.New(5),
		Server:   server,
		Room:     "community-" + string(id),
	}

	Communities[id] = c

	return c
}

func GetCommunity(id uint64) *Community {
	return Communities[id]
}

func (c *Community) State() (CommunityState, error) {
	c.RLock()
	var (
		waitlist = c.Waitlist
		media    = c.Media
		users    = c.Users
	)
	c.RUnlock()

	db := models.DB

	state := CommunityState{
		Waitlist: waitlist,
		Media:    media,
	}

	if err := db.
		Where("created_at > ? and community_id = ?", time.Now().Add(-10*time.Minute), c.ID).
		Order("created_at desc").
		Limit(255).
		Preload("User").
		Find(&state.Chat).
		Error; err != nil {
		return CommunityState{}, err
	}

	for i, j := 0, len(state.Chat)-1; i < j; i, j = i+1, j-1 {
		state.Chat[i], state.Chat[j] = state.Chat[j], state.Chat[i]
	}

	if err := db.
		Where(users).
		Find(&state.Users).Error; err != nil {
		return CommunityState{}, err
	}

	if err := db.
		Preload("Staff").
		Preload("Staff.User").
		Preload("Mutes").
		Preload("User").
		First(&state.Meta, c.ID).Error; err != nil {
		return CommunityState{}, err
	}

	if err := db.
		Where("community_id = ?", c.ID).
		Order("updated_at desc").
		Limit(50).
		Preload("User").
		Find(&state.Meta.History).Error; err != nil {
		return CommunityState{}, err
	}

	if media != nil {
		state.Media.Elapsed = time.Since(media.Timestamp).Seconds()
	}

	return state, nil
}

func (c *Community) Advance() error {
	c.Lock()
	defer c.Unlock()

	c.Timer.Stop()

	var payload CommunityAdvance

	tx := models.DB.Begin()

	var community models.Community
	if err := tx.First(&community, c.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if c.Media != nil {
		if err := tx.Model(models.Media{}).Where("id = ?", c.Media.Item.MediaID).Updates(models.G{
			"plays": gorm.Expr("plays + ?", 1),
			"woots": gorm.Expr("woots + ?", len(c.Media.Votes.Woots)),
			"mehs":  gorm.Expr("mehs + ?", len(c.Media.Votes.Mehs)),
			"saves": gorm.Expr("saves + ?", len(c.Media.Votes.Saves)),
		}).Error; err != nil {
			tx.Rollback()
			return err
		}

		payload.History = &models.History{
			CommunityID: c.ID,
			UserID:      c.Media.UserID,
			MediaID:     c.Media.Item.MediaID,
			Timestamp:   c.Media.Timestamp,
			Type:        c.Media.Item.Type,
			ContentID:   c.Media.Item.ContentID,
			Image:       c.Media.Item.Image,
			Duration:    c.Media.Item.Duration,
			Title:       c.Media.Item.Title,
			Artist:      c.Media.Item.Artist,
			Woots:       len(c.Media.Votes.Woots),
			Mehs:        len(c.Media.Votes.Mehs),
			Saves:       len(c.Media.Votes.Saves),
			Population:  len(c.Users),
		}

		if err := tx.Save(payload.History).Error; err != nil {
			tx.Rollback()
			return err
		}

		if community.DjCycle && !c.Skipped {
			c.Waitlist = append(c.Waitlist, c.Media.UserID)
		} else if c.Skipped {
			c.Skipped = false
		}
	}

	c.Media = nil

	if len(c.Waitlist) == 0 {
		go c.Server.Emit(c.Room, raiki.Event{"waitlist.update", c.Waitlist})
		go c.Server.Emit(c.Room, raiki.Event{"advance", payload})
		tx.Commit()
		return nil
	}

	uid := c.Waitlist[0]

	var playlist models.Playlist
	if err := tx.First(&playlist, "user_id = ? and active = ?", uid, true).Error; err != nil {
		tx.Rollback()
		return err
	}

	var playlistItem models.PlaylistItem
	if err := playlist.ItemFunc(func(items []models.PlaylistItem) []models.PlaylistItem {
		// Should panic if there are no items
		playlistItem = items[0]
		return append(items[1:], playlistItem)
	}); err != nil {
		tx.Rollback()
		return err
	}

	var media models.Media
	if err := tx.First(&media, playlistItem.MediaID).Error; err != nil {
		tx.Rollback()
		return err
	}

	c.Media = &CommunityMedia{
		UserID:    uid,
		Item:      playlistItem,
		Timestamp: time.Now(),
		Votes: CommunityVotes{
			Woots: []uint64{},
			Saves: []uint64{},
			Mehs:  []uint64{},
		},
	}

	payload.Media = c.Media

	c.Waitlist = c.Waitlist[1:]
	c.Timer = time.AfterFunc(time.Duration(media.Duration)*time.Second, func() {
		if err := c.Advance(); err != nil {
			panic(err)
		}
	})

	c.Server.Emit(c.Room, raiki.Event{"waitlist.update", c.Waitlist})
	c.Server.Emit(c.Room, raiki.Event{"advance", payload})

	tx.Commit()
	return nil
}

// Join adds a user to the community
// Returns true on success
// Returns false if the user is already in the community
func (c *Community) Join(id uint64) bool {
	c.Lock()
	defer c.Unlock()

	for _, user := range c.Users {
		if id == user {
			return false
		}
	}

	go func() {
		c.Pipe.Populate()
		defer c.Pipe.Free()
		var user models.User
		models.DB.Find(&user, id)
		c.Server.Emit(c.Room, raiki.Event{"user.join", user})
	}()

	c.Users = append(c.Users, id)
	return true
}

// Leave removes a user from the community
// Returns true on success
// Returns false if the user isn't in the community
func (c *Community) Leave(id uint64) bool {
	c.Lock()
	defer c.Unlock()

	for i, user := range c.Users {
		if id == user {
			go c.LeaveWaitlist(id)
			go c.Vote(id, -1)
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			c.Server.Emit(c.Room, raiki.Event{"user.leave", id})
			return true
		}
	}

	return false
}

func (c *Community) JoinWaitlist(id uint64) bool {
	c.Lock()
	defer c.Unlock()

	if c.Media != nil && c.Media.UserID == id {
		return false
	}

	for _, user := range c.Waitlist {
		if id == user {
			return false
		}
	}

	if c.Media == nil {
		go c.Advance()
	}

	c.Waitlist = append(c.Waitlist, id)
	c.Server.Emit(c.Room, raiki.Event{"waitlist.update", c.Waitlist})
	return true
}

func (c *Community) LeaveWaitlist(id uint64) bool {
	c.Lock()
	defer c.Unlock()

	if c.Media != nil && c.Media.UserID == id {
		c.Skipped = true
		go c.Advance()
	}

	for i, user := range c.Waitlist {
		if id == user {
			c.Waitlist = append(c.Waitlist[:i], c.Waitlist[i+1:]...)
			c.Server.Emit(c.Room, raiki.Event{"waitlist.update", c.Waitlist})
			return true
		}
	}

	return false
}

func (c *Community) Vote(id uint64, voteType enum.VoteType) bool {
	c.Lock()
	defer c.Unlock()

	if c.Media == nil || c.Media.UserID == id {
		return false
	}

	switch voteType {
	case -1:
		for i, user := range c.Media.Votes.Woots {
			if user == id {
				c.Media.Votes.Woots = append(c.Media.Votes.Woots[:i], c.Media.Votes.Woots[i+1:]...)
			}
		}

		for i, user := range c.Media.Votes.Mehs {
			if user == id {
				c.Media.Votes.Mehs = append(c.Media.Votes.Mehs[:i], c.Media.Votes.Mehs[i+1:]...)
			}
		}
	case enum.VoteTypeWoot:
		for i, user := range c.Media.Votes.Mehs {
			if user == id {
				c.Media.Votes.Mehs = append(c.Media.Votes.Mehs[:i], c.Media.Votes.Mehs[i+1:]...)
			}
		}
		for _, user := range c.Media.Votes.Woots {
			if user == id {
				return false
			}
		}
		c.Media.Votes.Woots = append(c.Media.Votes.Woots, id)
		c.Server.Emit(c.Room, raiki.Event{"community.vote.woot", id})
	case enum.VoteTypeSave:
		for _, user := range c.Media.Votes.Saves {
			if user == id {
				return true
			}
		}
		c.Media.Votes.Saves = append(c.Media.Votes.Saves, id)
		c.Server.Emit(c.Room, raiki.Event{"community.vote.save", id})
		c.Unlock()
		c.Vote(id, enum.VoteTypeWoot)
		c.Lock()
	case enum.VoteTypeMeh:
		for i, user := range c.Media.Votes.Woots {
			if user == id {
				c.Media.Votes.Woots = append(c.Media.Votes.Woots[:i], c.Media.Votes.Woots[i+1:]...)
			}
		}
		for _, user := range c.Media.Votes.Mehs {
			if user == id {
				return false
			}
		}
		c.Media.Votes.Mehs = append(c.Media.Votes.Mehs, id)
		c.Server.Emit(c.Room, raiki.Event{"community.vote.meh", id})
	}

	return true
}
