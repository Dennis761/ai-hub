import mongoose from 'mongoose';

const apiKeySchema = new mongoose.Schema(
  {
    modelName: {
      type: String,
      required: true,
      trim: true,
    },
    provider: {
      type: String,
      required: true,
      trim: true,
    },
    keyName: {
      type: String,
      required: true,
      trim: true,
    },
    keyValue: {
      type: String,
      required: true,
    },
    status: {
      type: String,
      enum: ['active', 'inactive'],
      default: 'active',
    },
    balance: {
      type: Number,
      default: null,
      min: 0,
    },
    usageEnv: {
      type: String,
      enum: ['dev', 'prod', 'test'],
      default: 'prod',
    },
    ownerId: {
      type: mongoose.Schema.Types.ObjectId,
      ref: 'Admin',
      required: true,
    },
  },
  { timestamps: true }
);

apiKeySchema.index({ keyName: 1, usageEnv: 1 }, { unique: true });

export default mongoose.models.ApiKey || mongoose.model('ApiKey', apiKeySchema);
