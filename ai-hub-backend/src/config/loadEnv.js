import dotenv from 'dotenv'
dotenv.config()

const loadEnv = {
  PORT: process.env.PORT,
  MONGODB_API_KEY: process.env.MONGODB_API_KEY,
  JWT_SECRET: process.env.JWT_SECRET,
  EMAIL_USER: process.env.EMAIL_USER,
  EMAIL_PASS: process.env.EMAIL_PASS,
  REDIS_URL: process.env.REDIS_URL,
  REDIS_PROJECT_EDIT_TTL: process.env.REDIS_PROJECT_EDIT_TTL,
  REDIS_PROJECT_CACHE_TTL: process.env.REDIS_PROJECT_CACHE_TTL,
  TIMEOUT: process.env.TIMEOUT,
  RESPONSE_PREFIX: process.env.RESPONSE_PREFIX,
  CRYPTO_ALGORITHM: process.env.CRYPTO_ALGORITHM,
  KEY_ENCRYPT_SECRET: process.env.KEY_ENCRYPT_SECRET,
  IV_LENGTH: process.env.IV_LENGTH
}

export default loadEnv
