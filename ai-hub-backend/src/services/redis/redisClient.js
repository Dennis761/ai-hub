import { createClient } from 'redis';
import loadEnv from '../../config/loadEnv.js';

const redisClient = createClient({
    url: loadEnv.REDIS_URL,
    socket: {
        connectTimeout: parseInt(loadEnv.TIMEOUT) 
      }
});

redisClient.on('error', (err) => {
    console.error('Redis error:', err);
});

redisClient.on('connect', () => {
    console.log('Connected to Redis');
});

await redisClient.connect();
  
const getAsync = async (key) => {
    return await redisClient.get(key);
};

const setAsync = async (key, value, expiration) => {
    return await redisClient.set(key, value, {
        EX: expiration,
    });
};

export { redisClient, getAsync, setAsync };