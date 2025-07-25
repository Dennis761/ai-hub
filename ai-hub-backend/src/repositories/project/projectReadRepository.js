import ProjectModel from '../../models/projectModel.js';

export default class ProjectReadRepository {
  async findById(id) {
    return await ProjectModel.findById(id);
  }

  async findByName(name) {
    return await ProjectModel.findOne({ name });
  }
 
  async findOwned(ownerId) {
    return await ProjectModel.find({ ownerId }).sort({ createdAt: -1 });
  }

  async findParticipating(adminId) {
    return await ProjectModel.find({ adminAccess: adminId }).sort({ createdAt: -1 });
  }
  
}
