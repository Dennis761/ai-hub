import { decrypt } from '../../utils/cryptoUtils.js';

export default class ProjectReadService {
  constructor(projectReadManager) {
    this.projectReadManager = projectReadManager;
  }

  async getMyProjects(adminId) {
    const [owned, participating] = await Promise.all([
      this.projectReadManager.getByOwner(adminId),
      this.projectReadManager.getByParticipant(adminId)
    ]);

    const all = [...owned, ...participating];

    const unique = Array.from(
      new Map(all.map(project => [project._id.toString(), project])).values()
    );

    return unique.map(project => ({
      _id: project._id,
      name: project.name,
      apiKey: decrypt(project.apiKey)
    }));
  }

  async getById(projectId, adminId) {
    const project = await this.projectReadManager.getById(projectId);
    if (!project) return null;

    const hasAccess =
      String(project.ownerId) === String(adminId) ||
      (Array.isArray(project.adminAccess) &&
        project.adminAccess.some(id => String(id) === String(adminId)));

    return hasAccess ? project : null;
  }
}
