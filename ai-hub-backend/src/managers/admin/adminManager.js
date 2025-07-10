import bcrypt from 'bcrypt';
import { generateVerificationCode, hashCode } from '../../utils/cryptoUtils.js';
import { generateToken } from '../../utils/tokenUtils.js';
import sendVerificationEmail from '../../services/nodemailer/nodemailer.js';

class AdminManager {
  constructor({ adminRepository }) {
    this.adminRepository = adminRepository;
  }

  async register({ email, password, name }) {
    const existingAdmin = await this.adminRepository.findByEmail(email);
    const { code, hash } = generateVerificationCode();
    const expiresAt = Date.now() + 10 * 60 * 1000;
    const hashedPassword = await bcrypt.hash(password, 10);

    if (existingAdmin) {
      if (existingAdmin.isVerified)
        throw new Error("This email is already registered");

      if (existingAdmin.verificationCodeExpires &&
          new Date() < new Date(existingAdmin.verificationCodeExpires)) {
        throw new Error("Please wait before requesting a new verification code");
      }

      Object.assign(existingAdmin, {
        name,
        password: hashedPassword,
        verificationCode: hash,
        verificationCodeExpires: expiresAt,
      });

      await this.adminRepository.save(existingAdmin);
      await sendVerificationEmail(email, code);
      return { message: "Verification code resent to email" };
    }

    await this.adminRepository.create({
      email,
      password: hashedPassword,
      name,
      verificationCode: hash,
      verificationCodeExpires: expiresAt,
      isVerified: false,
    });

    await sendVerificationEmail(email, code);
    return { message: "Verification code sent to email" };
  }

  async verifyEmail(email, code) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("User not found");
    if (admin.isVerified) throw new Error("Email already verified");

    if (!admin.verificationCode || Date.now() > admin.verificationCodeExpires)
      throw new Error("Code is invalid or expired");

    const hashedInput = hashCode(code);
    if (hashedInput !== admin.verificationCode)
      throw new Error("Invalid code");

    Object.assign(admin, {
      isVerified: true,
      verificationCode: null,
      verificationCodeExpires: null,
    });

    await this.adminRepository.save(admin);
    return { message: "Email successfully verified" };
  }

  async login(email, password) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("Invalid credentials");

    const isMatch = await bcrypt.compare(password, admin.password);
    if (!isMatch) throw new Error("Invalid credentials");

    if (!admin.isVerified) throw new Error("Email is not verified");

    const token = generateToken({ userId: admin._id });
    return {
      token,
      admin: {
        id: admin._id,
        email: admin.email,
        name: admin.name,
      },
    };
  }

  async requestPasswordReset(email) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("Admin not found");

    const { code, hash } = generateVerificationCode();
    const expiresAt = Date.now() + 10 * 60 * 1000;

    Object.assign(admin, {
      verificationCode: hash,
      verificationCodeExpires: expiresAt,
      isResetCodeConfirmed: false,
    });

    await this.adminRepository.save(admin);
    await sendVerificationEmail(email, code);
    return { message: "Verification code sent to email" };
  }

  async verifyResetCode(email, code) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin) throw new Error("Admin not found");

    const hashedInput = hashCode(code);
    if (
      !admin.verificationCode ||
      hashedInput !== admin.verificationCode ||
      Date.now() > admin.verificationCodeExpires
    ) {
      throw new Error("Code is invalid or expired");
    }

    admin.isResetCodeConfirmed = true;
    await this.adminRepository.save(admin);
    return { message: "Code confirmed. You can set a new password now" };
  }

  async setNewPassword(email, password) {
    const admin = await this.adminRepository.findByEmail(email);
    if (!admin || !admin.isResetCodeConfirmed) {
      throw new Error("Please confirm the reset code first");
    }

    admin.password = await bcrypt.hash(password, 10);
    Object.assign(admin, {
      verificationCode: null,
      verificationCodeExpires: null,
      isResetCodeConfirmed: false,
    });

    await this.adminRepository.save(admin);
    return { message: "Password changed successfully" };
  }
}

export default AdminManager;