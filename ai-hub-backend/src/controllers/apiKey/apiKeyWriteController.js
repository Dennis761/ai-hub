import { getModelInfoOrThrow } from '../../services/LLM/getModelInfoOrThrow.js';

class ApiKeyWriteController {
  constructor({ apiKeyWriteService, apiKeyInputService }) {
    this.apiKeyWriteService = apiKeyWriteService;
    this.apiKeyInputService = apiKeyInputService;
  }

  async create(req, res) {
    try {
      const cleaned = this.apiKeyInputService.normalize(req.body);
      this.apiKeyInputService.validate(cleaned);

      const modelInfo = getModelInfoOrThrow(cleaned.provider, cleaned.modelName);

      const key = await this.apiKeyWriteService.create({
        ...cleaned,
        costPerToken: modelInfo.input_cost_per_token,
        maxTokens: modelInfo.max_input_tokens,
        ownerId: req.user?._id,
      });

      res.status(201).json(key);
    } catch (err) {
      res.status(400).json({
        error: 'Failed to create key',
        details: err.message,
      });
    }
  }

  async update(req, res) {
    try {
      const { id } = req.params;
      const cleaned = this.apiKeyInputService.normalize(req.body);
      this.apiKeyInputService.validate(cleaned);

      const modelInfo = getModelInfoOrThrow(cleaned.provider, cleaned.modelName);
      cleaned.costPerToken = modelInfo.input_cost_per_token;
      cleaned.maxTokens = modelInfo.max_input_tokens;

      const updatedKey = await this.apiKeyWriteService.update(id, cleaned, req.user?._id);
      res.json(updatedKey);
    } catch (err) {
      res.status(400).json({ error: 'Failed to update API key', details: err.message });
    }
  }

  async delete(req, res) {
    try {
      await this.apiKeyWriteService.delete(req.params.id, req.user?._id);
      res.json({ message: 'Key deleted' });
    } catch (err) {
      res.status(500).json({ error: 'Failed to delete key', details: err.message });
    }
  }
}

export default ApiKeyWriteController;
