package apikeydomain

func (a *APIKey) Rename(newName APIKeyName) {
	if a.name.Value() == newName.Value() {
		return
	}
	a.name = newName
	a.touch()
}

func (a *APIKey) RebindProviderModel(provider ProviderName, modelName ModelName) {
	if a.provider.Value() == provider.Value() && a.model.Value() == modelName.Value() {
		return
	}
	a.provider = provider
	a.model = modelName
	a.touch()
}

func (a *APIKey) RotateEncryptedValue(newEncrypted EncryptedKeyValue) {
	a.keyValueEnc = newEncrypted
	a.touch()
}

func (a *APIKey) MoveToEnvironment(env UsageEnv) error {
	if _, ok := allowedApiKeyEnviroments[string(env)]; !ok {
		return InvalidUsageEnv()
	}

	if a.usageEnv == env {
		return nil
	}

	a.usageEnv = env
	a.touch()
	return nil
}

func (a *APIKey) SetStatus(newStatus APIKeyStatus) error {
	if _, ok := allowedApiKeyStatuses[string(newStatus)]; !ok {
		return InvalidStatus()
	}

	if a.status == newStatus {
		return nil
	}

	a.status = newStatus
	a.touch()
	return nil
}

func (a *APIKey) Activate() error {
	return a.SetStatus(APIKeyStatusActive)
}

func (a *APIKey) Deactivate() error {
	return a.SetStatus(APIKeyStatusInactive)
}

func (a *APIKey) SetBalance(newBalance APIKeyBalance) {
	if a.balance.Value() == newBalance.Value() {
		return
	}
	a.balance = newBalance
	a.touch()
}
