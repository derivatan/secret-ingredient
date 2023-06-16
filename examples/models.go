package main

import (
	"github.com/derivatan/si"
	"github.com/gofrs/uuid"
)

type contact struct {
	si.Model

	Email    string
	Phone    string
	ArtistID uuid.UUID

	artist si.RelationData[artist]
}

func (c contact) GetModel() si.Model {
	return c.Model
}

func (c contact) Artist() *si.Relation[contact, artist] {
	return si.BelongsTo[contact, artist](c, "artist", func(c *contact) *si.RelationData[artist] {
		return &c.artist
	})
}

type artist struct {
	si.Model

	Name string
	Year int

	contact si.RelationData[contact]
	albums  si.RelationData[album]
}

func (a artist) GetModel() si.Model {
	return a.Model
}

func (a artist) Contact() *si.Relation[artist, contact] {
	return si.HasOne[artist, contact](a, "contact", func(a *artist) *si.RelationData[contact] {
		return &a.contact
	})
}

func (a artist) Albums() *si.Relation[artist, album] {
	return si.HasMany[artist, album](a, "albums", func(a *artist) *si.RelationData[album] {
		return &a.albums
	})
}

type album struct {
	si.Model

	Name     string
	Year     int
	ArtistID uuid.UUID

	artist si.RelationData[artist]
}

func (a album) GetModel() si.Model {
	return a.Model
}

func (a album) Artist() *si.Relation[album, artist] {
	return si.BelongsTo[album, artist](a, "artist", func(a *album) *si.RelationData[artist] {
		return &a.artist
	})
}
