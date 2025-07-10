class TaskReadController {
    constructor( taskReadService ) {
      this.taskReadService = taskReadService;
    }
  
    async getTasksByProject(req, res) {
      try {
        const tasks = await this.taskReadService.getTasksByProject(req.params.projectId);
        res.json(tasks);
      } catch (error) {
        res.status(500).json({ error: 'Failed to get project tasks', details: error.message });
      }
    }
  
    async getTaskById(req, res) {
      try {
        const task = await this.taskReadService.getTaskById(req.params.id);
        if (!task) {
          return res.status(404).json({ error: 'Task not found' });
        }
        res.json(task);
      } catch (error) {
        res.status(500).json({ error: 'Failed to get task', details: error.message });
      }
    }
  }
  
  export default TaskReadController;
  