class ProjectAccessController {
    constructor(projectAccessService) {
      this.projectAccessService = projectAccessService;
    }
  
    async joinProject(req, res) {
      try {
        const { name, apiKey } = req.body;
        const adminId = req.user?._id;
   
        if (!name || !apiKey) {
          return res.status(400).json({ error: 'Missing "name" or "apiKey".' });
        }
  
        const result = await this.projectAccessService.joinProjectByName(name.trim(), apiKey.trim(), adminId);
  
        res.status(200).json(result);
      } catch (error) {
        res.status(400).json({ error: 'Failed to join project', details: error.message });
      }
    }
  }
  
  export default ProjectAccessController;
  