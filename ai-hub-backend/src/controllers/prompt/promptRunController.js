import { handleResponse } from '../../services/LLM/callLLMRequest.js';
import { checkBalance } from '../../services/LLM/checkLLMBalance.js';

class PromptRunController {
  constructor({ promptReadService, taskReadService, apiKeyWriteService, apiKeyReadService, promptInputService, promptHistoryService }) {
    this.promptReadService = promptReadService;
    this.taskReadService = taskReadService;
    this.apiKeyWriteService = apiKeyWriteService;
    this.apiKeyReadService = apiKeyReadService;
    this.promptInputService = promptInputService;
    this.promptHistoryService = promptHistoryService;
  }

  async run(req, res) {
    try {
      const { id } = req.params;
      const prompt = await this.promptReadService.getById(id);
      if (!prompt) return res.status(404).json({ error: 'Prompt not found' });

      const task = await this.taskReadService.getTaskById(String(prompt.taskId));
      if (!task?.apiMethod) return res.status(400).json({ error: 'Task or its apiMethod is missing' });

      const apiParams = this.promptInputService.parseQueryParamsFromApiMethod(task.apiMethod);
      const placeholders = this.promptInputService.extractPlaceholders(prompt.promptText);
      const missing = placeholders.filter(p => !(p in apiParams));
      if (missing.length) return res.status(400).json({ error: 'Missing parameters', missing });

      const formattedPrompt = this.promptInputService.formatPrompt(prompt.promptText, apiParams);
      const { modelName } = await this.apiKeyReadService.getModelById(prompt.modelId);
      const apiKey = await this.apiKeyReadService.getActiveKey(prompt.modelId);

      const balanceResult = await checkBalance(apiKey, formattedPrompt, modelName);
      if (!balanceResult.success) return res.status(400).json({ error: balanceResult.error });

      const responseText = await handleResponse(balanceResult.response, modelName);

      await this.apiKeyWriteService.update(apiKey._id, { balance: balanceResult.remainingBalance }, req.user._id);
      const updated = await this.promptHistoryService.addToHistory(id, prompt.promptText, responseText);

      res.json({
        ...updated.toObject(),
        cost: (balanceResult.cost || 0).toFixed(6),
        balanceRemaining: (balanceResult.remainingBalance || 0).toFixed(6)
      });
    } catch (error) {
      res.status(400).json({ error: 'Failed to run prompt', details: error.message });
    }
  }
}

export default PromptRunController;
