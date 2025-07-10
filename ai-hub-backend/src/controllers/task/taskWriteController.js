class TaskWriteController {
    constructor({ taskWriteService, taskInputService }) {
      this.taskWriteService = taskWriteService;
      this.taskInputService = taskInputService;
    }
  
    async create(req, res) {
      try {
        const adminId = req.user?._id;
        const input = this.taskInputService.normalize(req.body);
        this.taskInputService.validate(input);
  
        const task = await this.taskWriteService.create({
          ...input,
          createdBy: adminId
        });
  
        res.status(201).json(task);
      } catch (error) {
        res.status(400).json({ error: 'Failed to create task', details: error.message });
      }
    }
  
    async update(req, res) {
      try {
        const input = this.taskInputService.normalize(req.body);
        const updated = await this.taskWriteService.update(req.params.id, input);
        res.json(updated);
      } catch (error) {
        res.status(400).json({ error: 'Failed to update task', details: error.message });
      }
    }
  
    async delete(req, res) {
      try {
        const result = await this.taskWriteService.delete(req.params.id);
        res.json(result);
      } catch (error) {
        res.status(500).json({ error: 'Failed to delete task', details: error.message });
      }
    }
  }
  
  export default TaskWriteController;
  