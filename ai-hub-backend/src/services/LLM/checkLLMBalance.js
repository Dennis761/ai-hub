import { callLLM } from './callLLMRequest.js';

export async function checkBalance(apiKeyDoc, prompt, modelName) {
  const response = await callLLM(prompt, {
    service: apiKeyDoc.provider,
    model: modelName,
    apiKey: apiKeyDoc.keyValue,
    extended: true
  });

  const cost = response?.usage?.total_cost || 0;
  const remaining = apiKeyDoc.balance - cost;

  if (remaining < 0) {
    return {
      success: false,
      error: {
        error: 'Insufficient balance for prompt execution',
        required: cost.toFixed(6),
        available: apiKeyDoc.balance.toFixed(6)
      }
    };
  }

  return {
    success: true,
    cost,
    remainingBalance: remaining,
    response
  };
}
