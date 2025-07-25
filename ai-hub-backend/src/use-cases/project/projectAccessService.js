export default class ProjectAccessService {
    constructor(projectAccessManager) {
      this.projectAccessManager = projectAccessManager;
    }
  
    joinProjectByName(name, apiKey, adminId) {
      return this.projectAccessManager.joinByName(name, apiKey, adminId);
    }
  }
  