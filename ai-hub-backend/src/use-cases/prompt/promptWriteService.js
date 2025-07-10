class PromptWriteService {
  constructor({promptWriteManager, apiKeyReadService}) {
    this.promptWriteManager = promptWriteManager;
    this.apiKeyReadService = apiKeyReadService;
  }

  async create(data) {
    const prompt = await this.promptWriteManager.create(data);
    const { modelName, modelProvider } = await this.apiKeyReadService.getModelById(prompt.modelId);
    return { prompt, modelName, modelProvider };
  }

  async update(promptId, updates) {
    return this.promptWriteManager.update(promptId, updates);
  }

  async delete(promptId) {
    return this.promptWriteManager.delete(promptId);
  }
}

export default PromptWriteService;
