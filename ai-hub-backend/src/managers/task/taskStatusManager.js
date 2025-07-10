class TaskStatusManager {
    constructor(taskWriteRepository) {
      this.taskWriteRepository = taskWriteRepository;
    }
    
  async setStatus(id, status) {
    return this.taskWriteRepository.update(id, { status });
   }
  
  }
  
  export default TaskStatusManager;
  