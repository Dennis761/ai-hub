class ApiKeyReadService {
    constructor(apiKeyReadManager) {
      this.apiKeyReadManager = apiKeyReadManager;
    }
  
    getAllKeys(filter) {
      return this.apiKeyReadManager.getAllKeys(filter);
    }
   
    getModelById(modelId) {
      return this.apiKeyReadManager.getDecryptedKeyById(modelId);
    }
  }
  
  export default ApiKeyReadService;
  