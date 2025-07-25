class AdminController {
  constructor({ adminService, adminInputService }) {
    this.adminService = adminService;
    this.adminInputService = adminInputService;
  }

  async register(req, res) {
    try {
      const input = this.adminInputService.normalize(req.body);
      this.adminInputService.validateRegister(input);

      const response = await this.adminService.register(input);
      res.status(response.status).json({ message: response.message });
    } catch (error) {
      res.status(error.status || 400).json({ message: error.message });
    }
  }

  async verifyEmail(req, res) {
    try {
      const input = this.adminInputService.normalize(req.body);
      this.adminInputService.validateVerification(input);

      const response = await this.adminService.verifyEmail(input.email, input.code);
      res.json({ message: response.message });
    } catch (error) {
      res.status(error.status || 400).json({ message: error.message });
    }
  }

  async login(req, res) {
    try {
      const input = this.adminInputService.normalize(req.body);
      this.adminInputService.validateLogin(input);

      const result = await this.adminService.login(input.email, input.password);
      res.json(result);
    } catch (error) {
      res.status(error.status || 400).json({ message: error.message });
    }
  }

  async requestPasswordReset(req, res) {
    try {
      const { email } = this.adminInputService.normalize(req.body);
      if (!email) throw new Error('Email is required');

      const response = await this.adminService.requestPasswordReset(email);
      res.json({ message: response.message });
    } catch (error) {
      res.status(error.status || 400).json({ message: error.message });
    }
  }

  async verifyResetCode(req, res) {
    try {
      const input = this.adminInputService.normalize(req.body);
      this.adminInputService.validateVerification(input);

      const response = await this.adminService.verifyResetCode(input.email, input.code);
      res.json({ message: response.message });
    } catch (error) {
      res.status(error.status || 400).json({ message: error.message });
    }
  }

  async setNewPassword(req, res) {
    try {
      const { email, password } = this.adminInputService.normalize(req.body);
      if (!email || !password) {
        throw new Error('Email and new password are required');
      }

      const response = await this.adminService.setNewPassword(email, password);
      res.json({ message: response.message });
    } catch (error) {
      res.status(error.status || 400).json({ message: error.message });
    }
  }
}

export default AdminController;
