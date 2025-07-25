export default class ProjectReadManager {
  constructor(projectReadRepository) {
    this.projectReadRepository = projectReadRepository;
  }

  async getByOwner(ownerId) {
    return this.projectReadRepository.findOwned(ownerId);
  }

  async getByParticipant(adminId) {
    return this.projectReadRepository.findParticipating(adminId);
  }

  async getById(projectId) {
    return this.projectReadRepository.findById(projectId);
  }
}
