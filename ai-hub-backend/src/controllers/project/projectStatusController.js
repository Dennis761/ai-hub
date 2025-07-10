class ProjectStatusController {
    constructor(projectStatusService) {
      this.projectStatusService = projectStatusService;
    }
  
    async setStatus(req, res) {
      try {
        const status = req.body.status?.toLowerCase();
        const updated = await this.projectStatusService.setStatus(req.params.id, status);
        res.json(updated);
      } catch (error) {
        res.status(400).json({ error: 'Failed to update status', details: error.message });
      }
    }
  }
  
  export default ProjectStatusController;
  