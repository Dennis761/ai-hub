class PromptWriteManager {
  constructor({promptReadRepository, promptWriteRepository}) {
    this.promptReadRepository = promptReadRepository;
    this.promptWriteRepository = promptWriteRepository;
  }

  async create(data) {
    return this.promptWriteRepository.create(data);
  }

  async update(promptId, updates) {
    const prompt = await this.promptReadRepository.findById(promptId);
    if (!prompt) throw new Error('Prompt not found');

    if (typeof updates.promptText === 'string') {
      prompt.promptText = updates.promptText;
    }

    if (typeof updates.name === 'string') {
      prompt.name = updates.name;
    }

    return this.promptWriteRepository.save(prompt);
  }

  async delete(promptId) {
    return this.promptWriteRepository.delete(promptId);
  }
}

export default PromptWriteManager;