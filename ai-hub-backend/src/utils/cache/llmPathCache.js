import { getAsync, setAsync } from '../../services/redis/redisClient.js';

export const getCachedLLMPath = async (modelName) => {
  const cacheKey = `llm:path:${modelName}`;
  const json = await getAsync(cacheKey);
  return json ? JSON.parse(json) : null;
};

export const cacheLLMPath = async (modelName, path, ttl = 86400) => {
  const cacheKey = `llm:path:${modelName}`;
  await setAsync(cacheKey, JSON.stringify(path), ttl);
};
