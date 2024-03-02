package sites

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SiteController struct {
	siteService *SiteService
}

func NewSiteController(siteService *SiteService) *SiteController {
	return &SiteController{siteService: siteService}
}

// GetSites returns all sites.
func (sc *SiteController) GetSites(c *gin.Context) {
	sites, err := sc.siteService.GetSites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sites)
}

// GetAllWorkspaces returns all workspaces.
func (sc *SiteController) GetWorkspaces(c *gin.Context) {
	workspaces, err := sc.siteService.GetWorkspaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workspaces)
}

// CreateSite creates a new site.
func (sc *SiteController) CreateSite(c *gin.Context) {
	var site Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if site already exists
	if _, err := sc.siteService.GetSiteByName(site.Name); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Site already exists"})
		return
	}
	if err := sc.siteService.CreateSite(site); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, site)
}

// CreateWorkspace creates a new workspace.
func (sc *SiteController) CreateWorkspace(c *gin.Context) {
	var workspace Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if workspace already exists
	if _, err := sc.siteService.GetWorkspaceByName(workspace.Name); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workspace already exists"})
		return
	}
	if err := sc.siteService.CreateWorkspace(workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, workspace)
}

// DeleteSite deletes a site.
func (sc *SiteController) DeleteSite(c *gin.Context) {
	siteId := c.Param("id")
	if siteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site id"})
		return
	}
	if err := sc.siteService.DeleteSite(siteId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Site deleted"})
}

// DeleteWorkspace deletes a workspace.
func (sc *SiteController) DeleteWorkspace(c *gin.Context) {
	workspaceId := c.Param("id")
	if err := sc.siteService.DeleteWorkspace(workspaceId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workspace deleted"})
}

// UpdateSite updates a site.
func (sc *SiteController) UpdateSite(c *gin.Context) {
	var site Site
	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if site.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site id"})
		return
	}
	if err := sc.siteService.UpdateSite(site); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, site)
}

// UpdateWorkspace updates a workspace.
func (sc *SiteController) UpdateWorkspace(c *gin.Context) {
	var workspace Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if workspace.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace id"})
		return
	}
	if err := sc.siteService.UpdateWorkspace(workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workspace)
}
