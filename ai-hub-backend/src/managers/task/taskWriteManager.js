class TaskWriteManager {
  constructor({taskWriteRepository, taskReadRepository}) {
    this.taskWriteRepository = taskWriteRepository,
    this.taskReadRepository = taskReadRepository;
  }

  async create(data) {
    return this.taskWriteRepository.create(data);
  }

  async update(id, updates) {
    return this.taskWriteRepository.update(id, updates);
  }

  async delete(id) {
    const readTask = await this.taskReadRepository.findById(id)
    if (!readTask) throw new Error('Task not found');

    const writeTask = await this.taskWriteRepository.delete(id);
    if (!writeTask) throw new Error('Task not found');

    return { message: 'Task and related prompts deleted', projectId: readTask.projectId };
  }

}

export default TaskWriteManager;