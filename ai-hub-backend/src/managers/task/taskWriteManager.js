class TaskWriteManager {
  constructor(taskWriteRepository) {
    this.taskWriteRepository = taskWriteRepository;
  }

  async create(data) {
    return this.taskWriteRepository.create(data);
  }

  async update(id, updates) {
    return this.taskWriteRepository.update(id, updates);
  }

  async delete(id) {
    const task = await this.taskWriteRepository.delete(id);
    if (!task) throw new Error('Task not found');
    return { message: 'Task and related prompts deleted' };
  }

}

export default TaskWriteManager;