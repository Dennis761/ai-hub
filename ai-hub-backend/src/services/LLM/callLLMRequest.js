import extractLLMResponse from './extractLLMResponse.js';
import LLM from "@themaximalist/llm.js";

function cleanModelResponse(rawText) {
  if (!rawText || typeof rawText !== 'string') return '';

  let text = rawText.trim();

  text = text.replace(/^(AI|Response|Answer)[:\-–]\s*/i, '');

  text = text.replace(/^"(.*)"$/, '$1');

  text = text
    .replace(/\*\*(.*?)\*\*/g, '$1') // **bold**
    .replace(/\*(.*?)\*/g, '$1')     // *italic*
    .replace(/^\s*[-*]\s+/gm, '')    // lists with * or -
    .replace(/^#+\s*/gm, '');        // headers #, ##, ###

  text = text.replace(/\n+/g, ' ').replace(/\s+/g, ' ').trim();

  return text;
}

export async function callLLM(prompt, { service, model, apiKey, extended }) {
  return await LLM(prompt, { service, model, apiKey, extended });
}

export async function handleResponse(response, modelName) {
  const rawExtracted = await extractLLMResponse(response, modelName);
  return cleanModelResponse(rawExtracted?.answer);
}
