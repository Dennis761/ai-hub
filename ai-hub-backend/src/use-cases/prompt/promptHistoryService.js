class PromptHistoryService {
  constructor(promptHistoryManager) {
    this.promptHistoryManager = promptHistoryManager;
  }

  async addToHistory(promptId, promptText, responseText) {
    return this.promptHistoryManager.addToHistory(promptId, promptText, responseText);
  }

  async rollbackToHistoryVersion(promptId, versionIndex) {
    return this.promptHistoryManager.rollback(promptId, versionIndex);
  }
}

export default PromptHistoryService;
