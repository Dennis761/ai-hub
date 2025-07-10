class ApiKeyStatusService {
    constructor(apiKeyStatusManager) {
      this.apiKeyStatusManager = apiKeyStatusManager;
    }
  
    activateKey(id, adminId) {
      return this.apiKeyStatusManager.activateKey(id, adminId);
    }
  
    deactivateKey(id, adminId) {
      return this.apiKeyStatusManager.deactivateKey(id, adminId);
    }
  }
  
  export default ApiKeyStatusService;
  