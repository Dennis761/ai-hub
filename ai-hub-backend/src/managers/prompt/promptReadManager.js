class PromptReadManager {
  constructor(promptReadRepository) {
    this.promptReadRepository = promptReadRepository;
  }

  async getById(promptId) {
    return this.promptReadRepository.findById(promptId);
  }

  async getByTask(taskId) {
    return this.promptReadRepository.findByTaskId(taskId);
  }
}

export default PromptReadManager;
