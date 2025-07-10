export default class ProjectStatusService {
    constructor(projectStatusManager) {
      this.projectStatusManager = projectStatusManager;
    }
  
    setStatus(id, status) {
      return this.projectStatusManager.setStatus(id, status);
    }
  }
  