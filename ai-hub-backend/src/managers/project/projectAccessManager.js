import { decrypt } from '../../utils/cryptoUtils.js';

export default class ProjectAccessManager {
  constructor(projectReadRepository, projectWriteRepository) {
    this.projectReadRepository = projectReadRepository;
    this.projectWriteRepository = projectWriteRepository;
  }

  async getByOwner(ownerId) {
    const projects = await this.projectReadRepository.findOwned(ownerId);
    return projects.map(p => ({ ...p.toObject(), apiKey: decrypt(p.apiKey) }));
  }

  async getByParticipant(adminId) {
    const projects = await this.projectReadRepository.findParticipating(adminId);
    return projects.map(p => ({ ...p.toObject(), apiKey: decrypt(p.apiKey) }));
  }

  async joinByName(projectName, providedApiKey, adminId) {
    const project = await this.projectReadRepository.findByName(projectName);
    if (!project) throw new Error('Project with this name does not exist.');
    if (project.status !== 'active') throw new Error('Project is not active.');

    const decryptedKey = decrypt(project.apiKey);
    if (decryptedKey !== providedApiKey) throw new Error('Incorrect API key.');

    const adminIdStr = adminId.toString();
    const isOwner = project.ownerId.toString() === adminIdStr;
    const isAlreadyInList = project.adminAccess.some(id => id.toString() === adminIdStr);

    if (isOwner || isAlreadyInList) {
      return {
        alreadyJoined: true,
        project,
        message: 'You are already a member of this project.',
      };
    }

    project.adminAccess.push(adminId);
    await project.save();

    return {
      alreadyJoined: false,
      project,
      message: 'Successfully joined the project.',
    };
  }
}
