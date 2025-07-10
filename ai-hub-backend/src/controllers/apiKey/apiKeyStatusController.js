class ApiKeyStatusController {
    constructor(apiKeyStatusService) {
      this.apiKeyStatusService = apiKeyStatusService;
    }
  
    async activateKey(req, res) {
      try {
        const key = await this.apiKeyStatusService.activateKey(req.params.id, req.user?._id);
        res.json(key);
      } catch (err) {
        const status = err.message === 'Forbidden' ? 403 : 500;
        res.status(status).json({ error: 'Failed to activate key', details: err.message });
      }
    }
  
    async deactivateKey(req, res) {
      try {
        const key = await this.apiKeyStatusService.deactivateKey(req.params.id, req.user?._id);
        res.json(key);
      } catch (err) {
        const status = err.message === 'Forbidden' ? 403 : 500;
        res.status(status).json({ error: 'Failed to deactivate key', details: err.message });
      }
    }
  }
  
  export default ApiKeyStatusController;
  