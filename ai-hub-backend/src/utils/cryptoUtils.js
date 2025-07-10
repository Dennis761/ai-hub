import crypto from 'crypto';
import loadEnv from '../config/loadEnv.js';

const algorithm = loadEnv.CRYPTO_ALGORITHM;
const secretKey = loadEnv.KEY_ENCRYPT_SECRET;
const ivLength = parseInt(loadEnv.IV_LENGTH);

export function encrypt(text) {
  if (!secretKey || Buffer.byteLength(secretKey) !== 32) {
    throw new Error('Invalid secret key length: must be 32 bytes');
  }
  if (!ivLength || isNaN(ivLength)) {
    throw new Error('Invalid IV length');
  }

  const iv = crypto.randomBytes(ivLength);
  const cipher = crypto.createCipheriv(algorithm, Buffer.from(secretKey), iv);
  let encrypted = cipher.update(text);
  encrypted = Buffer.concat([encrypted, cipher.final()]);
  return iv.toString('hex') + ':' + encrypted.toString('hex');
}

export function decrypt(text) {
  if (!secretKey || Buffer.byteLength(secretKey) !== 32) {
    throw new Error('Invalid secret key length: must be 32 bytes');
  }

  const [ivHex, encryptedHex] = text.split(':');
  const iv = Buffer.from(ivHex, 'hex');
  const encryptedText = Buffer.from(encryptedHex, 'hex');
  const decipher = crypto.createDecipheriv(algorithm, Buffer.from(secretKey), iv);
  let decrypted = decipher.update(encryptedText);
  decrypted = Buffer.concat([decrypted, decipher.final()]);
  return decrypted.toString();
}

export function generateVerificationCode() {
  const code = crypto.randomInt(100000, 999999).toString();
  const hash = crypto.createHash('sha256').update(code).digest('hex');
  return { code, hash };
}

export function hashCode(code) {
  return crypto.createHash('sha256').update(code).digest('hex');
}

