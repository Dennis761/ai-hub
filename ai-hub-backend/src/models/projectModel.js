import mongoose from 'mongoose';

const projectSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
    unique: true
  },
  status: {
    type: String,
    enum: ['active', 'inactive', 'archived'],
    default: 'active'
  },
  apiKey: {
    type: String,
    required: true,
  },
  ownerId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Admin',
    required: true
  },
  adminAccess: [{
    type: mongoose.Schema.Types.ObjectId,
    ref: 'Admin'
  }]
}, { timestamps: true });

export default mongoose.model('Project', projectSchema);
