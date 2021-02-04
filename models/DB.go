package models

import (
	"crypto/sha512"
	"fmt"
	"pipeline"
	"time"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
)

var (
	algoliaClient = algoliasearch.NewClient("IL54TIA6D9", "e54d89a32483b34a827e7b57f3ce38b1")
	DB            *gorm.DB
)

func init() {
	db, err := gorm.Open("mysql", "phynix:phynixdb@/phynix?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	// db.LogMode(true)
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(
		&Ban{},
		&Chat{},
		&Community{},
		&GlobalBan{},
		&History{},
		&Media{},
		&Mute{},
		&Playlist{},
		&PlaylistItem{},
		&Session{},
		&Staff{},
		&User{},
		&Verification{},
	)

	DB = db

	// var communities []Community
	// if err := db.Preload("User").Find(&communities).Error; err != nil {
	// 	panic(err)
	// }

	// objects := make([]interface{}, len(communities))
	// for i, c := range communities {
	// 	objects[i] = c
	// }

	// index := algoliaClient.InitIndex("community")
	// if _, err := index.AddObjects(objects); err != nil {
	// 	panic(err)
	// }
}

type Model struct {
	ID        uint64     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type G map[string]interface{}

func Tokenize(s string) string {
	b := []byte(s)
	r := securecookie.GenerateRandomKey(64)
	return fmt.Sprintf("%x", sha512.Sum512(append(b, r...)))
}

var algoliaPipeline = pipeline.New(10)

func UploadToAlgolia(indexName string, obj interface{}) error {
	algoliaPipeline.Populate()
	defer algoliaPipeline.Free()

	index := algoliaClient.InitIndex(indexName)
	_, err := index.AddObjects([]interface{}{obj})
	return err
}
