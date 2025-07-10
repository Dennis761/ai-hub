class PromptWriteController {
    constructor({ promptWriteService, promptReadService, promptHistoryService, promptInputService, taskReadService }) {
      this.promptWriteService = promptWriteService;
      this.promptReadService = promptReadService;
      this.promptHistoryService = promptHistoryService;
      this.promptInputService = promptInputService;
      this.taskReadService = taskReadService;
    }
  
    async create(req, res) {
      try {
        const input = this.promptInputService.normalize(req.body);
        this.promptInputService.validate(input);
  
        const task = await this.taskReadService.getTaskById(input.taskId);
        if (!task?.apiMethod) {
          return res.status(400).json({ error: 'Task or its apiMethod is missing' });
        }
  
        this.promptInputService.validatePromptAgainstApi(input.promptText, task.apiMethod);
        const { prompt } = await this.promptWriteService.create(input);
        res.status(201).json(prompt);
      } catch (error) {
        res.status(400).json({ error: 'Failed to create prompt', details: error.message });
      }
    }
  
    async update(req, res) {
      try {
        const { id } = req.params;
        const input = this.promptInputService.normalize(req.body);

        const existing = await this.promptReadService.getById(id);

        if (!existing) return res.status(404).json({ error: 'Prompt not found' });
  
        if (input.promptText && input.promptText !== existing.promptText) {
          await this.promptHistoryService.addToHistory(id, input.promptText, null);
        }

        const updated = await this.promptWriteService.update(id, input);
        res.json(updated);
      } catch (error) {
        res.status(400).json({ error: 'Failed to update prompt', details: error.message });
      }
    }
  
    async delete(req, res) {
      try {
        await this.promptWriteService.delete(req.params.id);
        res.json({ message: 'Prompt deleted' });
      } catch (error) {
        res.status(500).json({ error: 'Failed to delete prompt', details: error.message });
      }
    }
  }
  
  export default PromptWriteController;