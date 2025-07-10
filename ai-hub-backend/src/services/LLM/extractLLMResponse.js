import { getAsync, setAsync } from '../redis/redisClient.js';
import loadEnv from '../../config/loadEnv.js';

async function extractLLMResponse(data, modelName) {
  const cacheKey = `llm:path:${modelName}`;

  const start = Date.now();

  const cachedPathJson = await getAsync(cacheKey);

  if (cachedPathJson) {
    const cachedPath = JSON.parse(cachedPathJson);
    const value = getValueByPath(data, cachedPath);
    if (typeof value === 'string' && value.startsWith(loadEnv.RESPONSE_PREFIX)) {
      const end = Date.now();
      console.log(`[CACHE HIT] ${modelName} — ${(end - start)}ms`);
      return { path: cachedPath, answer: value.replace(loadEnv.RESPONSE_PREFIX, '').trim() };
    }
  }

  const searchStart = Date.now();

  function find(data, path = []) {
    if (typeof data === 'string' && data.startsWith(loadEnv.RESPONSE_PREFIX)) {
      return { path, answer: data.replace(loadEnv.RESPONSE_PREFIX, '').trim() };
    }
    if (Array.isArray(data)) {
      for (let i = 0; i < data.length; i++) {
        const result = find(data[i], [...path, i]);
        if (result) return result;
      }
    } else if (typeof data === 'object' && data !== null) {
      for (const key in data) {
        const result = find(data[key], [...path, key]);
        if (result) return result;
      }
    }
    return null;
  }

  const result = find(data);
  const searchEnd = Date.now();

  console.log(`[CACHE MISS] ${modelName} — recursive search ${(searchEnd - searchStart)}ms`);

  if (result?.path) {
    await setAsync(cacheKey, JSON.stringify(result.path), 86400);
  }

  return result;
}

function getValueByPath(data, path) {
  return path.reduce((acc, key) => acc?.[key], data);
}

export default extractLLMResponse;
