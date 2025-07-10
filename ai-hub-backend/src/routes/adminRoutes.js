import express from 'express';
import adminContainer from '../di/admin/rootAdminContainer.js';

const router = express.Router();
const controller = adminContainer.adminController;

router.post('/register', controller.register.bind(controller));
router.post('/verify', controller.verifyEmail.bind(controller));
router.post('/login', controller.login.bind(controller));
router.post('/forgot-password', controller.requestPasswordReset.bind(controller));
router.post('/verify-reset-code', controller.verifyResetCode.bind(controller));
router.post('/set-new-password', controller.setNewPassword.bind(controller));

export default router;
