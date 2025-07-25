class PromptRunController {
  constructor( promptRunService ) {
    this.promptRunService = promptRunService;
  }

  async run(req, res) {
    try {
      const result = await this.promptRunService.runPrompt(req.params.id, req.user._id);
      res.json(result);
    } catch (error) {
      if (error.missing) {
        res.status(400).json({ error: error.message, missing: error.missing });
      } else {
        res.status(400).json({ error: 'Failed to run prompt', details: error.message });
      }
    }
  }
}

export default PromptRunController;
