class TaskWriteService {
    constructor(taskWriteManager) {
      this.taskWriteManager = taskWriteManager;
    }
  
    create(data) {
      return this.taskWriteManager.create(data);
    }
  
    update(id, updates) {
      return this.taskWriteManager.update(id, updates);
    }
  
    delete(id) {
      return this.taskWriteManager.delete(id);
    }
  }
  
  export default TaskWriteService;  