class TaskReadManager {
  constructor(taskReadRepository) {
    this.taskReadRepository = taskReadRepository;
  }

  async getAllTasks(filter) {
    return this.taskReadRepository.findAll(filter);
  }

  async getById(id) {
    return this.taskReadRepository.findById(id);
  }

  async checkOwnership(taskId, adminId) {
    const task = await this.taskReadRepository.findById(taskId);
    return task && String(task.createdBy) === String(adminId);
  }
}

export default TaskReadManager;
