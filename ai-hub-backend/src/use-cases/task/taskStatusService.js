class TaskStatusService {
    constructor(taskStatusManager) {
      this.taskStatusManager = taskStatusManager;
    }
  
    setStatus(id, status) {
      return this.taskStatusManager.setStatus(id, status);
    }
  
  }
  
  export default TaskStatusService;
  