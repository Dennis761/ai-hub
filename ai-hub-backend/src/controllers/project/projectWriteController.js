class ProjectWriteController {
    constructor({ projectWriteService, projectInputService }) {
      this.projectWriteService = projectWriteService;
      this.projectInputService = projectInputService;
    }
  
    async create(req, res) {
      try {
        const adminId = req.user?._id;
        if (!adminId) return res.status(401).json({ error: 'Unauthorized' });
  
        const input = this.projectInputService.normalize(req.body);
        this.projectInputService.validate(input);
  
        const project = await this.projectWriteService.create({
          ...input,
          ownerId: adminId,
          adminAccess: [adminId],
        });
  
        res.status(201).json(project);
      } catch (error) {
        if (error.code === 11000) {
          return res.status(409).json({
            error: 'Conflict: Project already exists.',
            fields: error.keyValue,
          });
        }
        res.status(400).json({ error: error.message });
      }
    }
  
    async update(req, res) {
      try {
        const updated = await this.projectWriteService.update(req.params.id, req.body);
        res.json(updated);
      } catch (error) {
        res.status(error.status || 400).json({
          error: 'Failed to update project',
          details: error.message,
        });
      }
    }
  
    async delete(req, res) {
      try {
        await this.projectWriteService.delete(req.params.id);
        res.json({ message: 'Project deleted' });
      } catch (error) {
        res.status(500).json({ error: 'Failed to delete project', details: error.message });
      }
    }
  }
  
  export default ProjectWriteController;
  