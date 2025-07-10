class AdminService {
  constructor(adminManager) {
    this.adminManager = adminManager;
  }

  register(data) {
    return this.adminManager.register(data);
  }

  verifyEmail(email, code) {
    return this.adminManager.verifyEmail(email, code);
  }

  login(email, password) {
    return this.adminManager.login(email, password);
  }

  requestPasswordReset(email) {
    return this.adminManager.requestPasswordReset(email);
  }

  verifyResetCode(email, code) {
    return this.adminManager.verifyResetCode(email, code);
  }

  setNewPassword(email, password) {
    return this.adminManager.setNewPassword(email, password);
  }
}

export default AdminService;
