class ApiKeyStatusManager {
    constructor(apiKeyWriteManager) {
      this.apiKeyWriteManager = apiKeyWriteManager;
    }
  
    activateKey(id, adminId) {
      return this.apiKeyWriteManager.update(id, { status: 'active' }, adminId);
    }
  
    deactivateKey(id, adminId) {
      return this.apiKeyWriteManager.update(id, { status: 'inactive' }, adminId);
    }
  }
  
  export default ApiKeyStatusManager;
  