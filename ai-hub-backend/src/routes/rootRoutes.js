import { Router } from 'express';

import apiKeyRoutes from './apiKeyRoutes.js';
import projectRoutes from './projectRoutes.js';
import taskRoutes from './taskRoutes.js';
import promptRoutes from './promptRoutes.js';
import adminRoutes from './adminRoutes.js';

const router = Router();

router.use('/api-keys', apiKeyRoutes);
router.use('/projects', projectRoutes);
router.use('/tasks', taskRoutes);
router.use('/prompts', promptRoutes);
router.use('/admin', adminRoutes); 

export default router;
