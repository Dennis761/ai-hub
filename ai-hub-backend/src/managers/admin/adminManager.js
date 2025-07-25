import bcrypt from 'bcrypt';
import { generateVerificationCode, hashCode } from '../../utils/cryptoUtils.js';
import { generateToken } from '../../utils/tokenUtils.js';
import sendVerificationEmail from '../../services/nodemailer/nodemailer.js';

class AdminManager {
  constructor({ adminRepository }) {
    this.adminRepository = adminRepository;
  }

  async register(data, existingAdmin) {
    const { code, hash } = generateVerificationCode();
    const hashedPassword = await bcrypt.hash(data.password, 10);
    const expiresAt = Date.now() + 10 * 60 * 1000;

    if (existingAdmin) {
      Object.assign(existingAdmin, {
        name: data.name,
        password: hashedPassword,
        verificationCode: hash,
        verificationCodeExpires: expiresAt,
      });

      await this.adminRepository.save(existingAdmin);
      await sendVerificationEmail(data.email, code);
      return { message: "Verification code resent to email", status: 200 };
    }

    await this.adminRepository.create({
      email: data.email,
      password: hashedPassword,
      name: data.name,
      verificationCode: hash,
      verificationCodeExpires: expiresAt,
      isVerified: false,
    });

    await sendVerificationEmail(data.email, code);
    return { message: "Verification code sent to email", status: 201 };
  }

  async verifyEmail(admin, code) {
    if (!admin.verificationCode || Date.now() > admin.verificationCodeExpires) {
      throw new Error("Code is invalid or expired");
    }

    const hashedInput = hashCode(code);
    if (hashedInput !== admin.verificationCode) {
      throw new Error("Invalid code");
    }

    Object.assign(admin, {
      isVerified: true,
      verificationCode: null,
      verificationCodeExpires: null,
    });

    await this.adminRepository.save(admin);
    return { message: "Email successfully verified" };
  }

  async login(admin, password) {
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

  async requestPasswordReset(admin) {
    const { code, hash } = generateVerificationCode();
    const expiresAt = Date.now() + 10 * 60 * 1000;

    Object.assign(admin, {
      verificationCode: hash,
      verificationCodeExpires: expiresAt,
      isResetCodeConfirmed: false,
    });

    await this.adminRepository.save(admin);
    await sendVerificationEmail(admin.email, code);
    return { message: "Verification code sent to email" };
  }

  async verifyResetCode(admin, code) {
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

  async setNewPassword(admin, password) {
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
