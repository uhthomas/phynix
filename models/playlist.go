package models

import (
	"aoi"
	"sync"
)

type Playlist struct {
	Model
	Name   string `json:"name"`
	UserID uint64 `json:"userID"`
	Active bool   `json:"active"`
	Cap    uint   `json:"cap"`

	Items []PlaylistItem `json:"items,omitempty"`
}

var playlistLocks = make(map[uint64]*sync.Mutex)

func (p Playlist) ItemFunc(f func([]PlaylistItem) []PlaylistItem) error {
	if _, ok := playlistLocks[p.ID]; !ok {
		playlistLocks[p.ID] = &sync.Mutex{}
	}
	playlistLocks[p.ID].Lock()
	defer playlistLocks[p.ID].Unlock()

	var items []PlaylistItem
	if err := DB.Where("playlist_id = ?", p.ID).Order("`position`").Find(&items).Error; err != nil {
		return err
	}

	items = f(items)
	if items == nil {
		return nil
	}

	tx := DB.Begin()

	// var wg sync.WaitGroup
	// doneChan := make(chan struct{})
	// errChan := make(chan error)
	// var err error

	// wg.Add(len(items))
	// for i := range items {
	// 	i := i
	// 	go func() {
	// 		defer wg.Done()
	// 		item := &items[i]
	// 		item.Order = i
	// 		if err := tx.Where("id = ?", item.ID).Assign(PlaylistItem{Order: i, PlaylistID: p.ID}).FirstOrCreate(item).Error; err != nil {
	// 			errChan <- err
	// 		}
	// 	}()
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(doneChan)
	// }()

	// select {
	// case e := <-errChan:
	// 	close(errChan)
	// 	err = e
	// case <-doneChan:
	// }

	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	l := &sync.Mutex{}
	if err := aoi.Loop(items, func(i int) error {
		item := &items[i]
		item.Position = i + 1
		l.Lock()
		defer l.Unlock()
		e := tx.
			Where("id = ?", item.ID).
			Assign(PlaylistItem{Position: item.Position, PlaylistID: p.ID}).
			FirstOrCreate(item).Error
		return e
	}); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("playlist_id = ?", p.ID).Not("id", func() (out []uint64) {
		for _, item := range items {
			out = append(out, item.ID)
		}
		return
	}()).Delete(&PlaylistItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// if err := tx.Unscoped().Where("playlist_id = ?", p.ID).Not(items).Delete(&PlaylistItem{}).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	tx.Commit()

	return nil
}
