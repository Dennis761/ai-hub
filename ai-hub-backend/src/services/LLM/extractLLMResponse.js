import loadEnv from '../../config/loadEnv.js';
import { getCachedLLMPath, cacheLLMPath } from '../../utils/cache/llmPathCache.js';

function getValueByPath(data, path) {
  return path.reduce((acc, key) => acc?.[key], data);
}

async function extractLLMResponse(data, modelName) {
  const start = Date.now();
  const cachedPath = await getCachedLLMPath(modelName);

  if (cachedPath) {
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
    await cacheLLMPath(modelName, result.path);
  }

  return result;
}

export default extractLLMResponse;
