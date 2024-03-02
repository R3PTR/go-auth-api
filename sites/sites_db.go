package sites

import (
	"context"
	"log"

	"github.com/R3PTR/go-auth-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SitesDbService struct {
	mongoClient *database.MongoDBClient
}

func NewSitesDbService(mongoClient *database.MongoDBClient) *SitesDbService {
	return &SitesDbService{mongoClient: mongoClient}
}

// GetSites returns all sites.
func (s *SitesDbService) GetSites() ([]Site, error) {
	var sites []Site

	cursor, err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var site Site
		err := cursor.Decode(&site)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}

	return sites, nil
}

// GetAllWorkspaces returns all workspaces.
func (s *SitesDbService) GetWorkspaces() ([]Workspace, error) {
	var workspaces []Workspace

	cursor, err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var workspace Workspace
		err := cursor.Decode(&workspace)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

// GetSiteByName returns a site by name.
func (s *SitesDbService) GetSiteByName(name string) (Site, error) {
	var site Site
	err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).FindOne(context.Background(), bson.M{"name": name}).Decode(&site)
	return site, err
}

// GetSiteById returns a site by id.
func (s *SitesDbService) GetSiteById(siteId string) (Site, error) {
	var site Site
	objectId, err := primitive.ObjectIDFromHex(siteId)
	if err != nil {
		log.Println("Invalid id")
	}
	err = s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&site)
	return site, err
}

// GetWorkspaceByName returns a workspace by name.
func (s *SitesDbService) GetWorkspaceByName(name string) (Workspace, error) {
	var workspace Workspace
	err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).FindOne(context.Background(), bson.M{"name": name}).Decode(&workspace)
	return workspace, err
}

// GetWorkspaceById returns a workspace by id.
func (s *SitesDbService) GetWorkspaceById(workspaceId string) (Workspace, error) {
	var workspace Workspace
	objectId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		log.Println("Invalid id")
	}
	err = s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&workspace)
	return workspace, err
}

// CreateSite creates a new site.
func (s *SitesDbService) CreateSite(site Site) error {
	_, err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).InsertOne(context.Background(), site)
	return err
}

// CreateWorkspace creates a new workspace.
func (s *SitesDbService) CreateWorkspace(workspace Workspace) error {
	_, err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).InsertOne(context.Background(), workspace)
	return err
}

// DeleteSite deletes a site.
func (s *SitesDbService) DeleteSite(siteId string) error {
	objectId, err := primitive.ObjectIDFromHex(siteId)
	if err != nil {
		log.Println("Invalid id")
	}
	_, err = s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).DeleteOne(context.Background(), bson.M{"_id": objectId})
	return err
}

// DeleteWorkspace deletes a workspace.
func (s *SitesDbService) DeleteWorkspace(workspaceId string) error {
	objectId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		log.Println("Invalid id")
	}
	_, err = s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).DeleteOne(context.Background(), bson.M{"_id": objectId})
	return err
}

// UpdateSite updates a site.
func (s *SitesDbService) UpdateSite(site Site) error {
	// Delete Id from stuct, to prevent overwriting
	siteId := site.Id
	site.Id = ""
	objectId, err := primitive.ObjectIDFromHex(siteId)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": objectId}
	result, err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).ReplaceOne(context.Background(), filter, site)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 0 {
		return nil
	}
	result, err = s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.SiteCollection).UpdateOne(context.Background(), filter, bson.M{"$set": site})
	if err != nil {
		return err
	}
	if result.ModifiedCount != 0 {
		return nil
	}
	return err
}

// UpdateWorkspace updates a workspace.
func (s *SitesDbService) UpdateWorkspace(workspace Workspace) error {
	// Delete Id from stuct, to prevent overwriting
	workspaceId := workspace.Id
	workspace.Id = ""
	objectId, err := primitive.ObjectIDFromHex(workspaceId)
	if err != nil {
		log.Println("Invalid id")
	}
	filter := bson.M{"_id": objectId}
	result, err := s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).ReplaceOne(context.Background(), filter, workspace)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 0 {
		return nil
	}
	result, err = s.mongoClient.GetCollection(s.mongoClient.Config.SiteDatabase, s.mongoClient.Config.WorkspaceCollection).UpdateOne(context.Background(), filter, bson.M{"$set": workspace})
	if err != nil {
		return err
	}
	if result.ModifiedCount != 0 {
		return nil
	}
	return err
}
