class PromptHistoryManager {
  constructor({promptReadRepository, promptWriteRepository}) {
    this.promptReadRepository = promptReadRepository;
    this.promptWriteRepository = promptWriteRepository;
  }

  async addToHistory(promptId, newPromptText, newResponseText) {
    const prompt = await this.promptReadRepository.findById(promptId);

    if (!prompt) throw new Error('Prompt not found');

    if (!newResponseText || newResponseText.trim() === '') {
      prompt.promptText = newPromptText;
      prompt.responseText = '';
      return this.promptWriteRepository.save(prompt);
    }

    const usedVersions = [prompt.version, ...prompt.history.map(h => h.version || 0)];
    const newVersion = Math.max(...usedVersions) + 1;

    prompt.history.push({
      prompt: prompt.promptText,
      response: newResponseText,
      version: prompt.version,
      createdAt: new Date(),
    });

    prompt.promptText = newPromptText;
    prompt.responseText = newResponseText;
    prompt.version = newVersion;

    return this.promptWriteRepository.save(prompt);
  }

  async rollback(promptId, targetVersion) {
    const prompt = await this.promptReadRepository.findById(promptId);
    if (!prompt) throw new Error('Prompt not found');

    const historyItem = prompt.history.find(h => h.version === targetVersion);
    if (!historyItem) throw new Error(`Version ${targetVersion} not found in history`);

    const usedVersions = [prompt.version, ...prompt.history.map(h => h.version || 0)];
    const currentVersion = Math.max(...usedVersions);

    prompt.history.push({
      prompt: historyItem.prompt,
      response: historyItem.response,
      version: currentVersion,
      createdAt: new Date(),
    });

    prompt.promptText = historyItem.prompt;
    prompt.responseText = historyItem.response;
    prompt.version = currentVersion + 1;

    await this.promptWriteRepository.save(prompt);

    return { version: currentVersion };
  }
}

export default PromptHistoryManager;
