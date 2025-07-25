import { handleResponse } from '../../services/LLM/callLLMRequest.js';
import { checkBalance } from '../../services/LLM/checkLLMBalance.js';

class PromptRunService {
  constructor({
    promptReadService,
    taskReadService,
    apiKeyWriteService,
    apiKeyReadService,
    promptInputService,
    promptHistoryService
  }) {
    this.promptReadService = promptReadService;
    this.taskReadService = taskReadService;
    this.apiKeyWriteService = apiKeyWriteService;
    this.apiKeyReadService = apiKeyReadService;
    this.promptInputService = promptInputService;
    this.promptHistoryService = promptHistoryService;
  }
 
  async runPrompt(promptId, userId) {
    // Retrieve the prompt by ID
    const prompt = await this.promptReadService.getById(promptId);
    if (!prompt) throw new Error('Prompt not found');
  
    // Retrieve the associated task and validate apiMethod presence
    const task = await this.taskReadService.getTaskById(String(prompt.taskId));
    if (!task?.apiMethod) throw new Error('Task or its apiMethod is missing');
  
    // Parse query parameters from task's apiMethod
    const apiParams = this.promptInputService.parseQueryParamsFromApiMethod(task.apiMethod);
  
    // Extract placeholders from prompt text
    const placeholders = this.promptInputService.extractPlaceholders(prompt.promptText);
  
    // Validate that all placeholders are present in the API params
    const missing = placeholders.filter(p => !(p in apiParams));
    if (missing.length) {
      const err = new Error('Missing parameters');
      err.missing = missing;
      throw err;
    }
  
    // Format the prompt text using parsed parameters
    const formattedPrompt = this.promptInputService.formatPrompt(prompt.promptText, apiParams);
  
    // Retrieve model name and a valid API key
    const apiKey = await this.apiKeyReadService.getModelById(prompt.modelId);
    
    // Check balance and estimate cost for this prompt execution
    const balanceResult = await checkBalance(apiKey, formattedPrompt, apiKey.modelName);
    if (!balanceResult.success) throw new Error(balanceResult.error);
  
    // Send prompt to LLM and get the response
    const responseText = await handleResponse(balanceResult.response, apiKey.modelName);
    
    // Update API key balance
    await this.apiKeyWriteService.update(apiKey._id, { balance: balanceResult.remainingBalance }, userId);
    
    // Save the prompt and response in the prompt history
    const updated = await this.promptHistoryService.addToHistory(promptId, prompt.promptText, responseText);

    // Return result with cost and remaining balance info
    return {
      ...updated.toObject(),
      cost: (balanceResult.cost || 0).toFixed(6),
      balanceRemaining: (balanceResult.remainingBalance || 0).toFixed(6)
    };
  }  
}

export default PromptRunService;
