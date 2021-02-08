package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

//DBClient DBClient
type DBClient interface {
	Count(filit bson.M) (int, error)
	FindOneAndUpdate(filit bson.M, setData interface{}, result interface{}) error
	FindOne(filit bson.M, result interface{}) error
	FindOneField(filit bson.M, field bson.M, result interface{}) error
	Create(setData interface{}) error
	FindOneAndRemove(filit bson.M, result interface{}) error
	Find(filit bson.M, result interface{}) error
	FindField(filit bson.M, field bson.M, result interface{}) error
	Remove(filit bson.M) error
}

// Count Count
func (coll Collection) Count(filit bson.M) (int, error) {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	result, err := c.Find(filit).Count()
	return result, err
}

//FindOneAndUpdate FindOneAndUpdate
func (coll Collection) FindOneAndUpdate(filit bson.M, setData interface{}, result interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	var err error
	bIsBson := false
	switch setData.(type) {
	case bson.M:
		bIsBson = true
	}
	if bIsBson {
		err = c.Update(filit, bson.M{"$set": setData})
	} else {
		setDataSerial, _ := bson.Marshal(setData)
		setDataMap := make(bson.M)
		bson.Unmarshal(setDataSerial, &setDataMap)
		err = c.Update(filit, bson.M{"$set": setDataMap})
	}
	err = c.Find(filit).One(result)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

//FindOne FindOne
func (coll Collection) FindOne(filit bson.M, result interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Find(filit).One(result)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

//FindOneField FindOneField
func (coll Collection) FindOneField(filit bson.M, field bson.M, result interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Find(filit).Select(field).One(result)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

//Create Create
func (coll Collection) Create(setData interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Insert(setData)
	return err
}

//FindOneAndRemove FindOneAndRemove
func (coll Collection) FindOneAndRemove(filit bson.M, result interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Find(filit).One(result)
	err = c.Remove(filit)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

//Find Find
func (coll Collection) Find(filit bson.M, result interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Find(filit).All(result)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

//FindField FindField
func (coll Collection) FindField(filit bson.M, field bson.M, result interface{}) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Find(filit).Select(field).All(result)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

//Remove Remove
func (coll Collection) Remove(filit bson.M) error {
	s := coll.Database.Session.Copy()
	defer s.Close()
	c := coll.Collection
	err := c.Remove(filit)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}
