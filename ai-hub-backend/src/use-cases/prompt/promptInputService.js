import loadEnv from "../../config/loadEnv.js";

class PromptInputService {
  normalize(data) {
    return {
      name: data.name?.trim(),
      promptText: data.promptText?.trim(),
      taskId: data.taskId?.trim(),
      modelId: data.modelId?.trim(),
      executionOrder: data.executionOrder
    };
  }

  validate(data) {
    if (!data.name || !data.promptText || !data.taskId || !data.modelId) {
      throw new Error('Missing required fields for prompt');
    }
  }

  extractPlaceholders(promptText) {
    return [...promptText.matchAll(/{{(.*?)}}/g)].map(m => m[1]);
  }

  parseQueryParamsFromApiMethod(apiMethod) {
    const queryPart = apiMethod.split('?')[1] || '';
    const params = new URLSearchParams(queryPart);
    const result = {};
    for (const [key, value] of params.entries()) {
      result[key] = value;
    }
    return result;
  }

  validatePromptAgainstApi(promptText, apiMethod) {
    const placeholders = this.extractPlaceholders(promptText);
    const apiParams = this.parseQueryParamsFromApiMethod(apiMethod);

    const missing = placeholders.filter(key => !(key in apiParams));
    if (missing.length > 0) {
      throw new Error(`Missing parameters in API method: ${missing.join(', ')}`);
    }
  }

  formatPrompt(promptText, apiParams) {
    // Builds a prompt for the AI with an instruction to start its response with the value of loadEnv.RESPONSE_PREFIX (e.g., "#$#"),
    // so the answer can be easily detected and extracted.
    const instruction = `Відповідай, починаючи з символів "${loadEnv.RESPONSE_PREFIX}", а далі пиши саму відповідь.`;
    // Replaces all placeholders like {{key}} in the promptText with their corresponding values from apiParams.
    const filled = Object.entries(apiParams).reduce(
      (acc, [key, value]) => acc.replaceAll(`{{${key}}}`, value),
      promptText
    );
    return `${instruction}\n\n${filled}`;
  }
}

export default PromptInputService;
