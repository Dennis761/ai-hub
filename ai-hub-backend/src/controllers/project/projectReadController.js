import mongoose from 'mongoose';

class ProjectReadController {
  constructor(projectReadService) {
    this.projectReadService = projectReadService;
  }

  async getMyProjects(req, res) {
    try {
      const adminId = req.user?._id;
      if (!adminId) return res.status(401).json({ error: 'Unauthorized' });

      if (!mongoose.Types.ObjectId.isValid(adminId)) {
        return res.status(400).json({ error: 'Invalid admin ID' });
      }

      const projects = await this.projectReadService.getMyProjects(adminId);
      res.json({ projects });
    } catch (error) {
      res.status(500).json({ error: 'Failed to fetch projects', details: error.message });
    }
  }

  async getProjectById(req, res) {
    try {
      const adminId = req.user?._id;
      const { id: projectId } = req.params;

      if (!adminId) return res.status(401).json({ error: 'Unauthorized' });

      if (!mongoose.Types.ObjectId.isValid(projectId)) {
        return res.status(400).json({ error: 'Invalid project ID' });
      }

      const project = await this.projectReadService.getById(projectId, adminId);

      if (!project) {
        return res.status(404).json({ error: 'Project not found or access denied' });
      }

      res.json(project);
    } catch (error) {
      res.status(500).json({ error: 'Failed to fetch project', details: error.message });
    }
  }
}

export default ProjectReadController;
