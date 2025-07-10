import mongoose from 'mongoose';

const historySchema = new mongoose.Schema({
  prompt: String,
  response: String,
  version: Number,
  createdAt: { type: Date, default: Date.now }
}, { _id: false });

const promptSchema = new mongoose.Schema({
  taskId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Task',
    required: true
  },
  name: { type: String, required: true },
  modelId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Model',
    required: true
  },
  promptText: { type: String, required: true },
  responseText: { type: String },
  history: [historySchema],
  executionOrder: {
    type: Number,
    required: true,
    default: 0
  },
  version: {
    type: Number,
    default: 1
  }
}, { timestamps: true });

export default mongoose.model('Prompt', promptSchema);