package sites

type SiteService struct {
	sitesDbService *SitesDbService
}

func NewSiteService(sitesDbService *SitesDbService) *SiteService {
	return &SiteService{sitesDbService: sitesDbService}
}

// GetSites returns all sites.
func (s *SiteService) GetSites() ([]Site, error) {
	return s.sitesDbService.GetSites()
}

// GetAllWorkspaces returns all workspaces.
func (s *SiteService) GetWorkspaces() ([]Workspace, error) {
	return s.sitesDbService.GetWorkspaces()
}

// GetSiteByName returns a site by name.
func (s *SiteService) GetSiteByName(name string) (Site, error) {
	return s.sitesDbService.GetSiteByName(name)
}

// GetWorkspaceByName returns a workspace by name.
func (s *SiteService) GetWorkspaceByName(name string) (Workspace, error) {
	return s.sitesDbService.GetWorkspaceByName(name)
}

// CreateSite creates a new site.
func (s *SiteService) CreateSite(site Site) error {
	return s.sitesDbService.CreateSite(site)
}

// CreateWorkspace creates a new workspace.
func (s *SiteService) CreateWorkspace(workspace Workspace) error {
	return s.sitesDbService.CreateWorkspace(workspace)
}

// DeleteSite deletes a site.
func (s *SiteService) DeleteSite(siteId string) error {
	return s.sitesDbService.DeleteSite(siteId)
}

// DeleteWorkspace deletes a workspace.
func (s *SiteService) DeleteWorkspace(workspaceId string) error {
	return s.sitesDbService.DeleteWorkspace(workspaceId)
}

// UpdateSite updates a site.
func (s *SiteService) UpdateSite(site Site) error {
	return s.sitesDbService.UpdateSite(site)
}

// UpdateWorkspace updates a workspace.
func (s *SiteService) UpdateWorkspace(workspace Workspace) error {
	return s.sitesDbService.UpdateWorkspace(workspace)
}
