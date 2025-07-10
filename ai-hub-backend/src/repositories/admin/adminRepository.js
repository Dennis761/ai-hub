import AdminModel from '../../models/adminModel.js';

class AdminRepository {
  async findByEmail(email) {
    return AdminModel.findOne({ email });
  }

  async create(data) {
    return AdminModel.create(data);
  }

  async save(admin) {
    return admin.save();
  }
}

export default AdminRepository;
