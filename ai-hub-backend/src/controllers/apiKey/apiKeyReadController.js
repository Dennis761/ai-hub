class ApiKeyReadController {
    constructor(apiKeyReadService) {
      this.apiKeyReadService = apiKeyReadService;
    }
  
    async getAllKeys(req, res) {
      try {
        const keys = await this.apiKeyReadService.getAllKeys({ ownerId: req.user._id });
        res.json(keys);
      } catch (err) {
        res.status(500).json({ error: 'Failed to fetch keys', details: err.message });
      }
    } 
  }
  
  export default ApiKeyReadController;
  