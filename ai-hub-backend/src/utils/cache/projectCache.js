import { redisClient } from '../../services/redis/redisClient.js';
import loadEnv from '../../config/loadEnv.js';

const incrementProjectEditCount = async (projectId) => {
  const key = `project:edit:count:${String(projectId)}`;
  await redisClient.incr(key);
  await redisClient.expire(key, loadEnv.REDIS_PROJECT_EDIT_TTL);
};

const cacheProject = async (projectId, projectData) => {
  const value = JSON.stringify(projectData);
  await redisClient.set(`project:cache:${String(projectId)}`, value, { EX: loadEnv.REDIS_PROJECT_CACHE_TTL });
};

const getCachedProject = async (projectId) => {
  const value = await redisClient.get(`project:cache:${String(projectId)}`);
  return value ? JSON.parse(value) : null;
};

const isInTop100Projects = async (projectId) => {
  const rank = await redisClient.zRevRank('project:edit:count', String(projectId));
  return rank !== null && rank < 100;
};

export {
  incrementProjectEditCount,
  cacheProject,
  getCachedProject,
  isInTop100Projects
};



