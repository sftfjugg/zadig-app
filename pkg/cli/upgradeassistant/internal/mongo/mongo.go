package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/koderover/zadig/pkg/config"
	"github.com/koderover/zadig/pkg/microservice/systemconfig/core/codehost/repository/models"
	"github.com/koderover/zadig/pkg/tool/log"
	mongotool "github.com/koderover/zadig/pkg/tool/mongo"
)

func (c *CodehostColl) ChangeType(ID int, sourceType string) error {
	query := bson.M{"id": ID}

	preType := sourceType
	if sourceType == "1" {
		sourceType = "gitlab"
	} else if sourceType == "2" {
		sourceType = "github"
	} else if sourceType == "3" {
		sourceType = "gerrit"
	} else if sourceType == "4" {
		sourceType = "codehub"
	}

	change := bson.M{"$set": bson.M{
		"type": sourceType,
	}}
	_, err := c.Collection.UpdateOne(context.TODO(), query, change)
	if err != nil {
		log.Error("repository update fail")
		return err
	}
	log.Infof("success to change id:%d type:%s to type:%s", ID, preType, sourceType)
	return nil
}

func (c *CodehostColl) RollbackChangeType(ID int, sourceType string) error {
	query := bson.M{"id": ID}

	if sourceType == "gitlab" {
		sourceType = "1"
	} else if sourceType == "github" {
		sourceType = "2"
	} else if sourceType == "gerrit" {
		sourceType = "3"
	} else if sourceType == "codehub" {
		sourceType = "4"
	}

	change := bson.M{"$set": bson.M{
		"type": sourceType,
	}}
	_, err := c.Collection.UpdateOne(context.TODO(), query, change)
	if err != nil {
		log.Error("repository update fail")
		return err
	}
	log.Infof("success to rollback id:%d to type:%s", ID, sourceType)
	return nil
}

type CodehostColl struct {
	*mongo.Collection

	coll string
}

type ListArgs struct {
	Owner   string
	Address string
	Source  string
}

func NewCodehostColl() *CodehostColl {
	name := models.CodeHost{}.TableName()
	coll := &CodehostColl{Collection: mongotool.Database(config.MongoDatabase()).Collection(name), coll: name}

	return coll
}
