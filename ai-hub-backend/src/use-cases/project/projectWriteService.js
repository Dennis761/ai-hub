export default class ProjectWriteService {
    constructor(projectWriteManager) {
      this.projectWriteManager = projectWriteManager;
    }
  
    create(data) {
      return this.projectWriteManager.create(data);
    }
  
    update(id, data) {
      return this.projectWriteManager.update(id, data);
    }
  
    delete(id) {
      return this.projectWriteManager.delete(id);
    }
  }
  