class PromptReadController {
    constructor( promptReadService ) {
      this.promptReadService = promptReadService;
    }
   
    async getPromptsByTask(req, res) {
      try {
        const prompts = await this.promptReadService.getByTask(req.params.taskId);
        res.json(prompts);
      } catch (error) {
        res.status(500).json({ error: 'Failed to fetch prompts', details: error.message });
      }
    }
  }
  
  export default PromptReadController;  