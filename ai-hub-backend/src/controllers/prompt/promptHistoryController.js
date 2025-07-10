class PromptHistoryController {
    constructor( promptHistoryService ) {
      this.promptHistoryService = promptHistoryService;
    }
  
    async rollback(req, res) {
      try {
        const { id } = req.params;
        const { versionIndex } = req.body;
        const updated = await this.promptHistoryService.rollbackToHistoryVersion(id, versionIndex);
        res.json(updated);
      } catch (error) {
        res.status(400).json({ error: 'Failed to rollback prompt', details: error.message });
      }
    }
  }
  
  export default PromptHistoryController;
  