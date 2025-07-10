export default class ProjectAccessService {
    constructor(projectAccessManager) {
      this.projectAccessManager = projectAccessManager;
    }
  
    getProjectsByOwner(ownerId) {
      return this.projectAccessManager.getByOwner(ownerId);
    }
  
    getProjectsByParticipant(adminId) {
      return this.projectAccessManager.getByParticipant(adminId);
    }
  
    joinProjectByName(name, apiKey, adminId) {
      return this.projectAccessManager.joinByName(name, apiKey, adminId);
    }
  }
  