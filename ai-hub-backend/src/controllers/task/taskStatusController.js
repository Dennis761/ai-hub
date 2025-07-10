class TaskStatusController {
    constructor({ taskStatusService, taskInputService }) {
      this.taskStatusService = taskStatusService;
      this.taskInputService = taskInputService;
    }
  
    async setStatus(req, res) {
      try {
        const { status } = req.body;
        this.taskInputService.validateStatus(status);
  
        const updated = await this.taskStatusService.setStatus(req.params.id, status.toLowerCase());
        res.json(updated);
      } catch (error) {
        res.status(400).json({ error: 'Failed to update task status', details: error.message });
      }
    }
  }
  
  export default TaskStatusController;
  