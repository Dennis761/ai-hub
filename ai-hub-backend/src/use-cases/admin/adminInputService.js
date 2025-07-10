class AdminInputService {
  normalize({ email, password, name, code }) {
    return {
      email: email?.trim(),
      password: password?.trim(),
      name: name?.trim(),
      code: code?.trim()
    };
  }

  validateRegister({ email, password, name }) {
    if (!email || !password || !name)
      throw new Error('All fields are required');
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
      throw new Error('Invalid email format');
    if (password.length < 6)
      throw new Error('Password must be at least 6 characters');
  }

  validateLogin({ email, password }) {
    if (!email || !password)
      throw new Error('Email and password are required');
  }

  validateVerification({ email, code }) {
    if (!email || !code)
      throw new Error('Email and code are required');
  }
}

export default AdminInputService;
