import express from 'express';
import checkAuth from '../middlewares/authMiddleware.js';
import apiKeyContainer from '../di/apiKey/rootApiKeyContainer.js';

const router = express.Router();

const {
  apiKeyReadController,
  apiKeyWriteController,
  apiKeyStatusController
} = apiKeyContainer;

// Read
router.get('/', checkAuth, apiKeyReadController.getAllKeys.bind(apiKeyReadController));

// Write
router.post('/', checkAuth, apiKeyWriteController.create.bind(apiKeyWriteController));
router.put('/:id', checkAuth, apiKeyWriteController.update.bind(apiKeyWriteController));
router.delete('/:id', checkAuth, apiKeyWriteController.delete.bind(apiKeyWriteController));

// Status
router.patch('/:id/activate', checkAuth, apiKeyStatusController.activateKey.bind(apiKeyStatusController));
router.patch('/:id/deactivate', checkAuth, apiKeyStatusController.deactivateKey.bind(apiKeyStatusController));

export default router;
