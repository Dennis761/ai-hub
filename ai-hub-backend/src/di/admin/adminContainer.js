import AdminController from '../../controllers/admin/adminController.js';
import AdminService from '../../use-cases/admin/adminService.js';
import AdminInputService from '../../use-cases/admin/adminInputService.js';
import AdminManager from '../../managers/admin/adminManager.js';
import AdminRepository from '../../repositories/admin/adminRepository.js';
// Repositories
const adminRepository = new AdminRepository();

// Managers
const adminManager = new AdminManager({ adminRepository }); 

// Services
const adminService = new AdminService({ 
  adminManager, 
  adminRepository 
});
const adminInputService = new AdminInputService();

// Controller
const adminController = new AdminController({
  adminService,
  adminInputService,
});

export default adminController;
