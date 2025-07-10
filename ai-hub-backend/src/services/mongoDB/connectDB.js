import mongoose from 'mongoose';
import loadEnv from '../../config/loadEnv.js';

const db = loadEnv.MONGODB_API_KEY

const connectDB = async () => {
    try {
        mongoose
        .connect(db, { useNewUrlParser: true, useUnifiedTopology: true })
        .then(() => {
          console.log('Connected to the database');
        })
        .catch((err) => {
          console.error('Error connecting to the database:', err);
        });
    } catch (err) {
        console.error('Error connecting to the database:', err);
        process.exit(1);
    }
};

export default connectDB;
