import jwt from 'jsonwebtoken';
import loadEnv from '../config/loadEnv.js';

const userSecretCode = loadEnv.JWT_SECRET;

export default (req, res, next) => {
  const token = (req.headers.authorization || '').replace(/Bearer\s?/, '');

  if (!token) {
    return res.status(401).json({ message: `You don't have access` });
  }

  try {
    const decoded = jwt.verify(token, userSecretCode);
    req.user = { _id: decoded.userId }; 
    next();
  } catch (error) {
    return res.status(401).json({ message: 'Your access failed', error: error.message });
  }
};
