class TaskReadService {
    constructor(taskReadManager) {
      this.taskReadManager = taskReadManager;
    }
  
    getTaskById(id) {
      return this.taskReadManager.getById(id);
    }
  
    getTasksByProject(projectId) {
      return this.taskReadManager.getAllTasks({ projectId });
    }
  
    getAllTasks(filter) {
      return this.taskReadManager.getAllTasks(filter);
    }
  
    getTasksByCreator(adminId) {
      return this.taskReadManager.getAllTasks({ createdBy: adminId });
    }
  
    checkOwnership(taskId, adminId) {
      return this.taskManager.checkOwnership(taskId, adminId);
    }
  }
  
  export default TaskReadService;
  