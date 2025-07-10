import LLM from '@themaximalist/llm.js'

const { ModelUsage } = LLM;

export function getModelInfoOrThrow(provider, modelName) {
    const modelInfo = ModelUsage.get(provider, modelName);
    if (!modelInfo) {
      throw new Error(`Model "${modelName}" is not supported by provider "${provider}".`);
    }
    return modelInfo;
  }