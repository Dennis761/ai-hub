class ApiKeyWriteService {
    constructor(apiKeyWriteManager) {
      this.apiKeyWriteManager = apiKeyWriteManager;
    }
  
    create(data) {
      return this.apiKeyWriteManager.create(data);
    }
  
    update(id, updates, adminId) {
      return this.apiKeyWriteManager.update(id, updates, adminId);
    }
  
    delete(id, adminId) {
      return this.apiKeyWriteManager.delete(id, adminId);
    }
  }
  
  export default ApiKeyWriteService;
  