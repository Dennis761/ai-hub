class AdminService {
  constructor({ adminManager, adminRepository }) {
    this.adminManager = adminManager;
    this.adminRepository = adminRepository;
  } 

  async register(data) {
    const existing = await this.adminRepository.findByEmail(data.email);

    if (existing?.isVerified) {
      throw new Error("This email is already registered");
    }

    if (existing && existing.verificationCodeExpires &&
        new Date() < new Date(existing.verificationCodeExpires)) {
      throw new Error("Please wait before requesting a new verification code");
    }

    return this.adminManager.register(data, existing);
  }

  async verifyEmail(email, code) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("User not found");
    if (admin.isVerified) throw new Error("Email already verified");

    return this.adminManager.verifyEmail(admin, code);
  }

  async login(email, password) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("Invalid credentials");

    return this.adminManager.login(admin, password);
  }

  async requestPasswordReset(email) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("Admin not found");

    return this.adminManager.requestPasswordReset(admin);
  }

  async verifyResetCode(email, code) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("Admin not found");

    return this.adminManager.verifyResetCode(admin, code);
  }

  async setNewPassword(email, password) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin || !admin.isResetCodeConfirmed) {
      throw new Error("Please confirm the reset code first");
    }

    return this.adminManager.setNewPassword(admin, password);
  }
}

export default AdminService;
