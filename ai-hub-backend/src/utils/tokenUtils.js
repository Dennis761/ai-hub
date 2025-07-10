import jwt from 'jsonwebtoken';
import loadEnv from '../config/loadEnv.js';

export function generateToken(payload, expiresIn = '3h') {
  return jwt.sign(payload, loadEnv.JWT_SECRET, { expiresIn });
}
