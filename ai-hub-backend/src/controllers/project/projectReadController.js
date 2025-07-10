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

      const [owned, participating] = await Promise.all([
        this.projectReadService.getProjectsByOwner(adminId),
        this.projectReadService.getProjectsByParticipant(adminId)
      ]);

      res.json({ owned, participating });
    } catch (error) {
      res.status(500).json({ error: 'Failed to fetch projects', details: error.message });
    }
  }
}

export default ProjectReadController;
