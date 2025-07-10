class PromptReadService {
  constructor(promptReadManager) {
    this.promptReadManager = promptReadManager;
  }

  async getById(promptId) {
    return this.promptReadManager.getById(promptId);
  }

  async getByTask(taskId) {
    return this.promptReadManager.getByTask(taskId);
  }
}

export default PromptReadService;
