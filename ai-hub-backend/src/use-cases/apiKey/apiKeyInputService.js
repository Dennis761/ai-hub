class ApiKeyInputService {
  normalize(data) {
    return {
      ...data,
      modelName: data.modelName?.trim(),
      provider: data.provider?.trim(),
      keyName: data.keyName?.trim(),
      usageEnv: data.usageEnv?.trim(),
      keyValue: data.keyValue?.trim()
    };
  }

  validate(data) {
    const required = ['modelName', 'provider', 'keyName', 'keyValue'];
    const missing = required.filter(field => !data[field]);
    if (missing.length) {
      throw new Error(`Missing fields: ${missing.join(', ')}`);
    }

    if (data.status && !['active', 'inactive'].includes(data.status)) {
      throw new Error('Invalid status. Allowed: active, inactive');
    }

    if (data.usageEnv && !['prod', 'dev', 'test'].includes(data.usageEnv)) {
      throw new Error('Invalid usageEnv. Allowed: prod, dev, test');
    }
  }
}

export default ApiKeyInputService;
